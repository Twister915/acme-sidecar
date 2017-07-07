package main

import (
	"github.com/Twister915/acme-sidecar/store"
	"github.com/Twister915/acme-sidecar/server"
	"os"
	"strconv"
)

func main() {
	srv := server.Server{
		TargetPort: getPort(),
		Domain: getDomain(),
		Store: store.GetProvider("kubernetes"),
	}

	srv.Start()
}

func getPort() int {
	p := os.Getenv("PORT")
	if len(p) > 0 {
		port, err := strconv.Atoi(p)
		if err == nil {
			return port
		}
	}

	return 443
}

func getDomain() string {
	d := os.Getenv("DOMAIN")
	if len(d) == 0 {
		panic("must specify DOMAIN env var")
	}
	return d
}
