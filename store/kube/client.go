package kube

import (
	"io/ioutil"
	"path/filepath"
)

type client struct {
	token      []byte
	ca         []byte
	clientCert []byte
	namespace  []byte
}

func newClient() (client *client, err error) {
	const (
		authDir = "/var/run/secrets/kubernetes.io/serviceaccount"

		token     = "token"
		ca        = "ca.crt"
		namespace = "namespace"
	)

	tokenData, err := ioutil.ReadFile(filepath.Join(authDir, token))
	if err != nil {
		return
	}

	caData, err := ioutil.ReadFile(filepath.Join(authDir, ca))
	if err != nil {
		caData = nil
	}

	nsData, err := ioutil.ReadFile(filepath.Join(authDir, namespace))
	if err != nil {
		return
	}

	client = &client{token: tokenData, ca: caData, namespace: nsData}
	return
}

func (c *client) requestIn(namespace string) *requestIn {
	return &requestIn{client: c, namespace: namespace}
}

func (c *client) request() *requestIn {
	if len(c.namespace) == 0 {
		panic("no default namespace known...")
	}
	return c.requestIn(string(c.namespace))
}
