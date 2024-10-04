package web

import (
	"net/http"
)

func NewServerConnection(addr string, router *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
