package server

import (
	"golang.org/x/crypto/acme/autocert"
	"os"
	"net/url"
	"fmt"
	"net/http/httputil"
	"net/http"
	"crypto/tls"
	"github.com/apex/log"
)

func (srv *Server) startTlsServer() {
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