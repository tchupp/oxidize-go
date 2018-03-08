package http

import (
	"net"
	"net/http"
	"time"
)

const readTimeoutSec = 10

type Server struct {
	server   *http.Server
	mux      *http.ServeMux
	listener net.Listener
}

func NewServer(listener net.Listener) *Server {
	mux := http.NewServeMux()
	server := &http.Server{Handler: mux, ReadTimeout: time.Second * readTimeoutSec}

	return &Server{
		server:   server,
		mux:      mux,
		listener: listener,
	}
}

func (s *Server) Register(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func (s *Server) RegisterMany(path string, handlers []http.Handler) {
	for _, handler := range handlers {
		s.mux.Handle(path, handler)
	}
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Serve() {
	if s.listener == nil {
		return
	}

	log.WithField("addr", s.listener.Addr()).Info("starting rpc server")
	go func() {
		log.Error(s.server.Serve(s.listener))
	}()
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	log.WithField("addr", s.listener.Addr()).Info("shutting down rpc server")
	return s.server.Close()
}
