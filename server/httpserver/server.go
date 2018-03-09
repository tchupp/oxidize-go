package httpserver

import (
	"context"

	"github.com/labstack/echo"
)

type Server struct {
	addr string
	*echo.Echo
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		Echo: echo.New(),
	}
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Serve() {
	log.WithField("addr", s.addr).Info("starting http server")
	go func() {
		log.Error(s.Echo.Start(s.addr))
	}()
}

func (s *Server) Close() error {
	log.WithField("addr", s.addr).Info("shutting down http server")
	return s.Echo.Shutdown(context.Background())
}
