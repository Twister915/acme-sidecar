package server

import (
	"fmt"
	"net/http"
	"golang.org/x/crypto/acme/autocert"
	"os"
	"crypto/tls"
	"github.com/apex/log"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	ListenPort int
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

	ctx := log.WithFields(log.Fields{
		"email": manager.Email,
	})
	ctx.Infof("configured let's encrypt")

	target, err := url.Parse(fmt.Sprintf("http://localhost:%d", srv.TargetPort))
	if err != nil {
		ctx.WithError(err).Error("failed to put url together")
		return
	}

	revProxy := httputil.NewSingleHostReverseProxy(target)

	s := &http.Server{
		Addr:      fmt.Sprintf(":%d", srv.ListenPort),
		TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
		Handler: revProxy,
	}

	must(s.ListenAndServeTLS("", ""))
}