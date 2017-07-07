package kube

import (
	"io/ioutil"
	"path/filepath"
	"github.com/apex/log"
)

type client struct {
	token      []byte
	ca         []byte
	clientCert []byte
	namespace  []byte
}

func newClient() (c *client, err error) {
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

	log.Infof("token => %s", tokenData)

	caData, err := ioutil.ReadFile(filepath.Join(authDir, ca))
	if err != nil {
		caData = nil
		log.WithError(err).Warn("could not read ca from dir")
	}

	nsData, err := ioutil.ReadFile(filepath.Join(authDir, namespace))
	if err != nil {
		return
	}

	log.Infof("namespace => %s", nsData)

	c = &client{token: tokenData, ca: caData, namespace: nsData}
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
