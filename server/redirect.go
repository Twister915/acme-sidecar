package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) startRedirectServer() {
	http.ListenAndServe(fmt.Sprintf(":%d", srv.RedirectPort), http.HandlerFunc(redirectHandler))
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	u := *req.URL
	u.Host = req.Host
	u.Scheme = "https"
	w.Header().Add("Location", u.String())
	w.WriteHeader(http.StatusPermanentRedirect)
}