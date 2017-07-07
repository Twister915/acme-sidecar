package server

import (
	"fmt"
	"net/http"
	"golang.org/x/crypto/acme/autocert"
	"os"
	"crypto/tls"
)

type Server struct {
	TargetPort int
	Domain string
	Store autocert.Cache
}

func (srv *Server) Start() {
	manager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: srv.Store,
		Email: os.Getenv("LETSENCRYPT_EMAIL"),
		HostPolicy: autocert.HostWhitelist(srv.Domain),
	}

	s := &http.Server{
		Addr:      fmt.Sprintf(":%d", srv.TargetPort),
		TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
	}
	must(s.ListenAndServeTLS("", ""))
}