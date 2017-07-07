package store

import (
	"github.com/Twister915/acme-sidecar/certs"
	"github.com/Twister915/acme-sidecar/store/kube"
	"fmt"
)

var registeredProviders = make(map[string]func() Provider)

type Provider interface {
	Get(id string) (cert certs.Cert, err error)
	Put(id string, cert certs.Cert) (err error)
}

func register(id string, factory func() Provider) {
	registeredProviders[id] = factory
}

func init() {
	register("kubernetes", func() Provider {
		return kube.NewProvider()
	})
}

func GetProvider(name string) Provider {
	factory, has := registeredProviders[name]
	if !has {
		panic(fmt.Sprintf("could not find provider %s", name))
	}
	return factory()
}