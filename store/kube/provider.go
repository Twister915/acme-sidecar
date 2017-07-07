package kube

import (
	"encoding/base64"
	"fmt"
	"context"
	"golang.org/x/crypto/acme/autocert"
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

func (p *Provider) Get(ctx context.Context, key string) (secretContents []byte, err error) {
	secret, err := p.client.request().prepare("GET", "v1", "secrets", key)()
	if err != nil {
		return
	}

	if secret.Status() == 404 {
		err = autocert.ErrCacheMiss
		return
	}

	data, err := secret.AsMap()
	if err != nil {
		return
	}

	secretContents, err = base64.StdEncoding.DecodeString(data["data"].(map[string]interface{})["data"].(string))
	return
}

func (p *Provider) Put(ctx context.Context, key string, data []byte) (err error) {
	err = p.sendSecret(ctx, key, data, "PUT")
	if err != nil {
		err = p.sendSecret(ctx, key, data, "POST")
	}
	return
}

func (p *Provider) sendSecret(ctx context.Context, key string, data []byte, method string) (err error) {
	resp, err := p.client.request().prepare(method, "v1", "secrets", key)(map[string]interface{}{
		"apiVersion": "v1",
		"data": map[string]string{
			"data": base64.StdEncoding.EncodeToString(data),
		},
		"kind": "Secret",
		"metadata": map[string]interface{}{
			"name": key,
			"namespace": string(p.client.namespace),
		},
		"type": "Opaque",
	})

	if err != nil {
		return
	}

	if resp.statusCode / 100 != 2 {
		err = fmt.Errorf("err: %s", string(resp.data))
	}

	return
}

func (p *Provider) Delete(ctx context.Context, key string) (err error) {
	_, err = p.client.request().prepare("DELETE", "v1", "secrets", key)()
	return
}