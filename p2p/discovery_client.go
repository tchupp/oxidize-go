package p2p

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/rpc"
)

type DiscoveryClient interface {
	Ping() error
	Version() (*entity.Hash, error)
}

type discoveryClient struct {
	client rpc.DiscoveryServiceClient
}

func NewDiscoveryClient(conn *grpc.ClientConn) DiscoveryClient {
	client := rpc.NewDiscoveryServiceClient(conn)

	return &discoveryClient{client: client}
}

func (c *discoveryClient) Ping() error {
	request := &rpc.PingRequest{}

	ctx := context.Background()
	_, err := c.client.Ping(ctx, request)
	return err
}

func (c *discoveryClient) Version() (*entity.Hash, error) {
	request := &rpc.VersionRequest{}

	ctx := context.Background()
	response, err := c.client.Version(ctx, request)
	if err != nil {
		return nil, err
	}

	hash, err := entity.NewHash(response.GetLatestHash())
	if err != nil {
		return nil, err
	}

	return hash, nil
}
