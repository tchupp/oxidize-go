package blockrpc

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/server/rpc"
	"github.com/tclchiam/oxidize-go/wire"
)

type syncBackend interface {
	BestHeader() (*entity.BlockHeader, error)
	HeaderByHash(hash *entity.Hash) (*entity.BlockHeader, error)
	Headers(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)
	BlockByHash(hash *entity.Hash) (*entity.Block, error)
}

type syncServer struct {
	backend syncBackend
}

func NewSyncServer(backend syncBackend) SyncServiceServer {
	return &syncServer{backend: backend}
}

func RegisterSyncServer(s *rpc.Server, srv SyncServiceServer) {
	s.Register(&_SyncService_serviceDesc, srv)
}

func (s *syncServer) GetBestHeader(ctx context.Context, req *GetBestHeaderRequest) (*GetBestHeaderResponse, error) {
	header, err := s.backend.BestHeader()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding best header: %s", err)
	}

	return &GetBestHeaderResponse{
		Header: wire.ToWireBlockHeader(header),
	}, nil
}

func (s *syncServer) GetHeaders(ctx context.Context, req *GetHeadersRequest) (*GetHeadersResponse, error) {
	hash, err := entity.NewHash(req.GetLatestHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested starting header hash was invalid: '%s'", req.GetLatestHash())
	}

	headers, err := s.backend.Headers(hash, req.GetLatestIndex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding previous headers")
	}

	return &GetHeadersResponse{
		HeaderCount: proto.Uint32(uint32(len(headers))),
		Headers:     wire.ToWireBlockHeaders(headers),
	}, nil
}

func (s *syncServer) GetBlock(ctx context.Context, req *GetBlockRequest) (*GetBlockResponse, error) {
	hash, err := entity.NewHash(req.GetHash())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "requested starting header hash was invalid: '%s'", req.GetHash())
	}

	block, err := s.backend.BlockByHash(hash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding previous headers")
	}

	return &GetBlockResponse{
		Block: wire.ToWireBlock(block),
	}, nil
}
