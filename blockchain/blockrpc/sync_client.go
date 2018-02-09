package blockrpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
)

type SyncClient interface {
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)

	GetBlock(hash *entity.Hash, index uint64) (*entity.Block, error)
}

type syncClient struct {
	client SyncServiceClient
}

func NewSyncClient(conn *grpc.ClientConn) SyncClient {
	client := NewSyncServiceClient(conn)

	return &syncClient{client: client}
}

func (c *syncClient) GetBestHeader() (*entity.BlockHeader, error) {
	request := &GetBestHeaderRequest{}

	ctx := context.Background()
	response, err := c.client.GetBestHeader(ctx, request)
	if err != nil {
		return nil, err
	}

	header, err := encoding.FromWireBlockHeader(response.GetHeader())
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *syncClient) GetHeaders(latestHash *entity.Hash, latestIndex uint64) (entity.BlockHeaders, error) {
	request := &GetHeadersRequest{
		LatestHash:  latestHash.Slice(),
		LatestIndex: proto.Uint64(latestIndex),
	}

	ctx := context.Background()
	response, err := c.client.GetHeaders(ctx, request)
	if err != nil {
		return nil, err
	}

	headers, err := encoding.FromWireBlockHeaders(response.GetHeaders())
	if err != nil {
		return nil, err
	}
	return headers, nil
}

func (c *syncClient) GetBlock(hash *entity.Hash, index uint64) (*entity.Block, error) {
	request := &GetBlockRequest{
		Hash:  hash.Slice(),
		Index: proto.Uint64(index),
	}

	ctx := context.Background()
	response, err := c.client.GetBlock(ctx, request)
	if err != nil {
		return nil, err
	}

	block, err := encoding.FromWireBlock(response.GetBlock())
	if err != nil {
		return nil, err
	}
	return block, nil
}
