package store

import (
	"github.com/Twister915/acme-sidecar/store/kube"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
)

var registeredProviders = make(map[string]func() autocert.Cache)

func register(id string, factory func() autocert.Cache) {
	registeredProviders[id] = factory
}

func init() {
	register("kubernetes", func() autocert.Cache {
		return kube.NewProvider()
	})
}

func GetProvider(name string) autocert.Cache {
	factory, has := registeredProviders[name]
	if !has {
		panic(fmt.Sprintf("could not find provider %s", name))
	}
	return factory()
}