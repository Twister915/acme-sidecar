package kube

import (
	"github.com/Twister915/acme-sidecar/certs"
	"github.com/Twister915/acme-sidecar/store/errors"
	"encoding/base64"
	"time"
	"fmt"
)

type Provider struct {
	client *client
}

func NewProvider() *Provider {
	client, err := newClient()
	if err != nil {
		panic(err)
	}

	return &Provider{client: client}
}

func (p *Provider) Get(id string) (cert certs.Cert, err error) {
	secret, err := p.client.request().prepare("GET", "v1", "secrets", id)()
	if err != nil {
		return
	}

	if secret.Status() == 404 {
		err = errors.ErrNotExist
		return
	}


	data, err := secret.AsMap()
	secretContents := data["data"].(map[string]interface{})
	cert.Expires, err = time.Parse(time.RFC822Z, string(readContentsSecret(secretContents, "expires")))
	if err != nil {
		panic(err)
	}
	cert.Domain = string(readContentsSecret(secretContents, "domain"))
	cert.Certificate = readContentsSecret(secretContents, "certificate")
	cert.Key = readContentsSecret(secretContents, "key")
	return
}

func (p *Provider) Put(id string, cert certs.Cert) (err error) {
	data := make(map[string]string)
	putContentsSecret(data, "expires", []byte(cert.Expires.Format(time.RFC822Z)))
	putContentsSecret(data, "domain", []byte(cert.Domain))
	putContentsSecret(data, "certificate", cert.Certificate)
	putContentsSecret(data, "key", cert.Key)

	secret := map[string]interface{}{
		"apiVersion": "v1",
		"data": data,
		"kind": "Secret",
		"metadata": map[string]interface{}{
			"name": id,
			"namespace": string(p.client.namespace),
		},
		"type": "Opaque",
	}

	resp, err := p.client.request().prepare("PUT", "v1", "secrets", id)(secret)
	if err != nil {
		return
	}

	if resp.statusCode / 100 != 2 {
		err = fmt.Errorf("err: %s", string(resp.data))
	}

	return
}

func readContentsSecret(data map[string]interface{}, key string) (out []byte) {
	var err error
	out, err = base64.StdEncoding.DecodeString(data[key].(string))
	if err != nil {
		panic(err)
	}
	return
}

func putContentsSecret(data map[string]string, key string, d []byte) {
	data[key] = base64.StdEncoding.EncodeToString(d)
}