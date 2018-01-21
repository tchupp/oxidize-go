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

func (s *Server) Start() {
	log.Debugf("starting server on: %s", s.listener.Addr())
	go s.server.Serve(s.listener)
}

func (s *Server) Shutdown() error {
	log.Debugf("shutting down server on: %s", s.listener.Addr())
	return s.listener.Close()
}
