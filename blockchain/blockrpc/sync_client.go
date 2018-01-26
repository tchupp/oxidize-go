package blockrpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/rpc"
)

type SyncClient interface {
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)
}

type syncClient struct {
	client rpc.SyncServiceClient
}

func NewSyncClient(conn *grpc.ClientConn) SyncClient {
	client := rpc.NewSyncServiceClient(conn)

	return &syncClient{client: client}
}

func (c *syncClient) GetBestHeader() (*entity.BlockHeader, error) {
	request := &rpc.GetBestHeaderRequest{}

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
	request := &rpc.GetHeadersRequest{
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
