package p2p

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type discoveryBackend interface {
	GetBestHeader() (*entity.BlockHeader, error)
}

type discoveryServer struct {
	backend discoveryBackend
}

func NewDiscoveryServer(backend discoveryBackend) rpc.DiscoveryServiceServer {
	return &discoveryServer{backend: backend}
}

func (s *discoveryServer) Ping(ctx context.Context, req *rpc.PingRequest) (*rpc.PingResponse, error) {
	rpc.LoggerFromContext(ctx).Debugf("handled ping")
	return &rpc.PingResponse{}, nil
}

func (s *discoveryServer) Version(ctx context.Context, req *rpc.VersionRequest) (*rpc.VersionResponse, error) {
	handleRequest := func() (*rpc.VersionResponse, error) {
		header, err := s.backend.GetBestHeader()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error finding best header")
		}

		return &rpc.VersionResponse{
			LatestIndex: proto.Uint64(header.Index),
			LatestHash:  header.Hash.Slice(),
		}, nil
	}

	logger := rpc.LoggerFromContext(ctx)

	response, err := handleRequest()
	if err != nil {
		logger.WithError(err).Warnf("error finding headers: %s", err)
		return nil, err
	}

	logger.Debugf("handled version")
	return response, nil
}
