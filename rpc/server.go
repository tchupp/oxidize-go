package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewServer(listener net.Listener) *Server {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(grpcLogger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(grpcLogger),
		)),
	)

	return &Server{server: grpcServer, listener: listener}
}

func (s *Server) Register(sd *grpc.ServiceDesc, ss interface{}) {
	s.server.RegisterService(sd, ss)
}

func (s *Server) Serve() {
	log.WithField("addr", s.listener.Addr()).Info("starting server")
	go s.server.Serve(s.listener)
}

func (s *Server) Shutdown() {
	log.WithField("addr", s.listener.Addr()).Info("shutting down server")
	s.server.GracefulStop()
}
