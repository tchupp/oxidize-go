package blockrpc

import (
	log "github.com/sirupsen/logrus"
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
	doRequest := func() (*entity.BlockHeader, error) {
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

	requestLogger := log.WithFields(log.Fields{"request": "GetBestHeader"})
	headers, err := doRequest()
	if err != nil {
		requestLogger.Warnf("error requesting header: %s", err)
		return nil, err
	}

	requestLogger.Debug("successfully received header")
	return headers, nil
}

func (c *syncClient) GetHeaders(latestHash *entity.Hash, latestIndex uint64) (entity.BlockHeaders, error) {
	doRequest := func() (entity.BlockHeaders, error) {
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

	requestLogger := log.WithFields(log.Fields{"request": "GetHeaders", "latestHash": latestHash, "latestIndex": latestIndex})
	headers, err := doRequest()
	if err != nil {
		requestLogger.Warnf("error: %s", err)
		return nil, err
	}

	requestLogger.Debug("success")
	return headers, nil
}
