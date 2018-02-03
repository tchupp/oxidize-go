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
			grpc_logrus.StreamServerInterceptor(log),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(log),
		)),
	)

	return &Server{server: grpcServer, listener: listener}
}

func (s *Server) RegisterSyncServer(service SyncServiceServer) {
	RegisterSyncServiceServer(s.server, service)
}

func (s *Server) RegisterDiscoveryServer(service DiscoveryServiceServer) {
	RegisterDiscoveryServiceServer(s.server, service)
}

func (s *Server) Serve() {
	log.WithField("addr", s.listener.Addr()).Info("starting server")
	go s.server.Serve(s.listener)
}

func (s *Server) Shutdown() {
	log.WithField("addr", s.listener.Addr()).Info("shutting down server")
	s.server.GracefulStop()
}
