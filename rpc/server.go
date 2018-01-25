package rpc

import (
	log "github.com/sirupsen/logrus"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewServer(listener net.Listener) *Server {
	return &Server{server: grpc.NewServer(), listener: listener}
}

func (s *Server) RegisterSyncServer(service SyncServiceServer) {
	RegisterSyncServiceServer(s.server, service)
}

func (s *Server) RegisterDiscoveryServer(service DiscoveryServiceServer) {
	RegisterDiscoveryServiceServer(s.server, service)
}

func (s *Server) Serve() {
	log.WithField("addr", s.listener.Addr()).Debug("starting server")
	go s.server.Serve(s.listener)
}

func (s *Server) Shutdown() {
	log.WithField("addr", s.listener.Addr()).Debug("shutting down server")
	s.server.GracefulStop()
}
