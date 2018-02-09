package p2p

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type DiscoveryClient interface {
	Ping() error
	Version() (*entity.Hash, error)
}

type discoveryClient struct {
	client DiscoveryServiceClient
}

func NewDiscoveryClient(conn *grpc.ClientConn) DiscoveryClient {
	client := NewDiscoveryServiceClient(conn)

	return &discoveryClient{client: client}
}

func (c *discoveryClient) Ping() error {
	request := &PingRequest{}

	ctx := context.Background()
	_, err := c.client.Ping(ctx, request)
	return err
}

func (c *discoveryClient) Version() (*entity.Hash, error) {
	request := &VersionRequest{}

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
