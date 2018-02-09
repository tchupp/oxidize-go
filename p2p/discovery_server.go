package p2p

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/rpc"
)

type discoveryBackend interface {
	GetBestHeader() (*entity.BlockHeader, error)
}

type discoveryServer struct {
	backend discoveryBackend
}

func NewDiscoveryServer(backend discoveryBackend) DiscoveryServiceServer {
	return &discoveryServer{backend: backend}
}

func RegisterDiscoveryServer(s *rpc.Server, srv DiscoveryServiceServer) {
	s.Register(&_DiscoveryService_serviceDesc, srv)
}

func (s *discoveryServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	return &PingResponse{}, nil
}

func (s *discoveryServer) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	header, err := s.backend.GetBestHeader()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding best header: %s", err)
	}

	return &VersionResponse{
		LatestIndex: proto.Uint64(header.Index),
		LatestHash:  header.Hash.Slice(),
	}, nil
}
