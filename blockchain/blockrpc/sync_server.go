package blockrpc

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/rpc"
)

type syncBackend interface {
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)
	GetBlock(hash *entity.Hash) (*entity.Block, error)
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

	headers, err := s.backend.GetHeaders(hash, req.GetLatestIndex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding previous headers")
	}

	return &rpc.GetHeadersResponse{
		HeaderCount: proto.Uint32(uint32(len(headers))),
		Headers:     encoding.ToWireBlockHeaders(headers),
	}, nil
}

func (s *syncServer) GetBlock(ctx context.Context, req *rpc.GetBlockRequest) (*rpc.GetBlockResponse, error) {
	hash, err := entity.NewHash(req.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested starting header hash was invalid: '%s'", req.GetHash())
	}

	block, err := s.backend.GetBlock(hash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding previous headers")
	}

	return &rpc.GetBlockResponse{
		Block: encoding.ToWireBlock(block),
	}, nil
}
