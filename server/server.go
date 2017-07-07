package server

import (
	"golang.org/x/crypto/acme/autocert"
	"sync"
)

type Server struct {
	ListenPort int
	TargetPort int
	RedirectPort int
	Domain string
	Store autocert.Cache
}

func (srv *Server) Start() {
	var wg sync.WaitGroup
	do := func(f func()) {
		defer wg.Done()
		f()
	}

	wg.Add(1)
	go do(srv.startTlsServer)
	if srv.RedirectPort > 0 {
		wg.Add(1)
		go do(srv.startRedirectServer)
	}

	wg.Wait()
}