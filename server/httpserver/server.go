package httpserver

import (
	"context"

	"net"

	"github.com/labstack/echo"
)

type Server struct {
	lis net.Listener
	*echo.Echo
}

func NewServer(lis net.Listener) *Server {
	e := echo.New()
	e.Listener = lis

	return &Server{
		lis:  lis,
		Echo: e,
	}
}

func (s *Server) Addr() string {
	return s.lis.Addr().String()
}

func (s *Server) Serve() {
	log.WithField("addr", s.Addr()).Info("starting http server")
	go func() {
		log.Error(s.Echo.Start(s.Addr()))
	}()
}

func (s *Server) Close() error {
	log.WithField("addr", s.Addr()).Info("shutting down http server")
	return s.Echo.Shutdown(context.Background())
}
