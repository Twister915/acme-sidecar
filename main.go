package main

import (
	"github.com/Twister915/acme-sidecar/store"
	"github.com/Twister915/acme-sidecar/server"
	"os"
	"strconv"
	"github.com/apex/log"
)

func main() {
	log.Info("starting application")
	srv := server.Server{
		ListenPort: getPort("LISTEN"),
		TargetPort: getPort("TARGET"),
		RedirectPort: getPort("REDIRECT"),
		Domain: getDomain(),
		Store: store.GetProvider("kubernetes"),
	}
	ctx := log.WithFields(log.Fields{
		"listen": srv.ListenPort,
		"target": srv.TargetPort,
		"domain": srv.Domain,
	})
	ctx.Info("server configured")

	srv.Start()
}

func getPort(name string) int {
	p := os.Getenv(name + "_PORT")
	if len(p) > 0 {
		port, err := strconv.Atoi(p)
		if err == nil {
			return port
		}
	}

	switch name {
	case "LISTEN":
		log.Warn("using default port 443 to listen")
		return 443
	case "REDIRECT":
		switch os.Getenv("REDIRECT_HTTP") {
		case "1", "true":
			return 80
		case "0", "false":
			return 0
		default:
			log.Warn("redirect server is not enabled")
			return 0
		}
	default:
		panic("must specify the TARGET_PORT env var")
	}
}

func getDomain() string {
	d := os.Getenv("DOMAIN")
	if len(d) == 0 {
		panic("must specify DOMAIN env var")
	}
	return d
}
