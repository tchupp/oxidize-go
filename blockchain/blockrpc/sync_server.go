package blockrpc

import (
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/rpc"
)

type syncBackend interface {
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
}

type syncServer struct {
	backend syncBackend
}

func NewSyncServiceServer(backend syncBackend) rpc.SyncServiceServer {
	return &syncServer{backend: backend}
}

func (s *syncServer) GetBestHeader(context.Context, *rpc.GetBestHeaderRequest) (*rpc.GetBestHeaderResponse, error) {
	header, err := s.backend.GetBestHeader()
	if err != nil {
		return nil, err
	}

	return &rpc.GetBestHeaderResponse{
		Header: encoding.ToWireBlockHeader(header),
	}, nil
}

func (s *syncServer) GetHeaders(ctx context.Context, req *rpc.GetHeadersRequest) (*rpc.GetHeadersResponse, error) {
	hash, err := entity.NewHash(req.GetLatestHash())
	if err != nil {
		log.WithField("hash", req.GetLatestHash()).
			Warn("Requested starting header hash was invalid")
		return nil, status.Errorf(codes.InvalidArgument, "Requested starting header hash '%s' was invalid", req.GetLatestHash())
	}

	header, err := s.backend.GetHeader(hash)
	if err != nil {
		log.WithField("hash", hash).
			Warn("Requested starting header was not found")
		return nil, status.Errorf(codes.NotFound, "Requested starting header was not found", hash)
	}

	return &rpc.GetHeadersResponse{
		HeaderCount: proto.Uint32(1),
		Headers: []*encoding.BlockHeader{
			encoding.ToWireBlockHeader(header),
		},
	}, nil
}
