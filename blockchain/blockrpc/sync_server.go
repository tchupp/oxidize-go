package blockrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

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

func NewSyncServer(backend syncBackend) rpc.SyncServiceServer {
	return &syncServer{backend: backend}
}

func (s *syncServer) GetBestHeader(ctx context.Context, req *rpc.GetBestHeaderRequest) (*rpc.GetBestHeaderResponse, error) {
	header, err := s.backend.GetBestHeader()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding best header: %s", err)
	}

	return &rpc.GetBestHeaderResponse{
		Header: encoding.ToWireBlockHeader(header),
	}, nil
}

func (s *syncServer) GetHeaders(ctx context.Context, req *rpc.GetHeadersRequest) (*rpc.GetHeadersResponse, error) {
	hash, err := entity.NewHash(req.GetLatestHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested starting header hash was invalid: '%s'", req.GetLatestHash())
	}

	header, err := s.backend.GetHeader(hash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "requested starting header was not found with hash: '%s'", hash)
	}

	return &rpc.GetHeadersResponse{
		HeaderCount: proto.Uint32(1),
		Headers: []*encoding.BlockHeader{
			encoding.ToWireBlockHeader(header),
		},
	}, nil
}
