package kube

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"encoding/json"
	"github.com/apex/log"
)

type requestIn struct {
	client    *client
	namespace string
}

type preparedReq func(body ...interface{}) (response responseData, err error)

func (r *requestIn) prepare(method, version, resource string, urlComponents ...interface{}) preparedReq {
	const (
		kube_host_env = "KUBERNETES_SERVICE_HOST"
		kube_port_env = "KUBERNETES_SERVICE_PORT"
	)

	//prepare the request, with a new http client
	var httpClient http.Client
	httpClient.Timeout = time.Second * 5

	//a new transport
	var transport http.Transport
	transport.TLSClientConfig = &tls.Config{}

	tlsCfg := transport.TLSClientConfig
	if r.client.ca != nil {
		tlsCfg.RootCAs = x509.NewCertPool()
		tlsCfg.RootCAs.AppendCertsFromPEM(r.client.ca)
	}

	httpClient.Transport = transport

	//and a new request
	var req http.Request
	req.Method = method

	convertedComponents := make([]string, len(urlComponents))
	for i, cmp := range urlComponents {
		convertedComponents[i] = fmt.Sprintf("%v", cmp)
	}
	urlCmps := strings.Join(convertedComponents, "/")
	if len(urlCmps) > 0 {
		urlCmps = "/" + urlCmps
	}

	urlStr := fmt.Sprintf("https://%s:%s/api/%s/namespaces/%s/%s%s",
		os.Getenv(kube_host_env), os.Getenv(kube_port_env),
		version, r.namespace,
		resource, urlCmps)

	var urlErr error
	req.URL, urlErr = url.Parse(urlStr)
	if urlErr != nil {
		panic(urlErr)
	}

	return func(body ...interface{}) (resp responseData, err error) {
		req := req //copy of the current request
		req.Header = make(http.Header)
		if len(r.client.token) > 0 {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", string(r.client.token)))
		}

		req.Header.Add("Accept", "application/json; charset=utf-8")
		if len(body) != 0 {
			if method == "GET" {
				panic("cannot set body on GET")
			}

			b := body[0]
			var reader io.ReadCloser
			switch bod := b.(type) {
			case []byte:
				reader = ioutil.NopCloser(bytes.NewReader(bod))
			case string:
				reader = ioutil.NopCloser(bytes.NewReader([]byte(bod)))
			case io.Reader:
				reader = ioutil.NopCloser(bod)
			case io.ReadCloser:
				reader = bod
			default:
				data, err := json.Marshal(bod)
				if err != nil {
					panic(fmt.Sprintf("cannot encode body for request -> %s", err.Error()))
				}
				reader = ioutil.NopCloser(bytes.NewReader(data))
				req.Header.Add("Content-Type", "application/json")
			}
			req.Body = reader
		}

		ctx := log.WithFields(log.Fields{
			"method": method,
			"version": version,
			"resource": resource,
			"url": urlStr,
		})

		ctx.Info("sending kubernetes request")
		start := time.Now()
		respData, err := httpClient.Do(&req)
		end := time.Now()
		if err != nil {
			ctx.WithError(err).Warnf("kube request failed")
			return
		}
		defer respData.Body.Close()
		resp.statusCode = respData.StatusCode
		resp.data, err = ioutil.ReadAll(respData.Body)
		if err != nil {
			ctx.WithError(err).Warnf("could not read body")
			return
		}
		ctx.Infof("kube request took %s", end.Sub(start).String())
		return
	}
}
