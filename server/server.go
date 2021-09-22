package server

import (
	"context"
	"net/http"
	"time"
)

var AppConf *ServerConf

type Server struct {
	httpServer *http.Server
}

type ServerConf struct {
	IP       string
	Port     string
	Protocol string
}

func (s *Server) Run(ip, port, protocol string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           "0.0.0.0:" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 28,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	AppConf = &ServerConf{
		IP:       ip,
		Port:     port,
		Protocol: protocol,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
