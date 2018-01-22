package p2p

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/rpc"
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
	doRequest := func() error {
		request := &rpc.PingRequest{}

		ctx := context.Background()
		_, err := c.client.Ping(ctx, request)
		return err
	}

	requestLogger := log.WithFields(log.Fields{"request": "Ping"})
	err := doRequest()
	if err != nil {
		requestLogger.Warnf("error: %s", err)
		return err
	}

	requestLogger.Debug("success")
	return nil
}

func (c *discoveryClient) Version() (*entity.Hash, error) {
	doRequest := func() (*entity.Hash, error) {
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

	requestLogger := log.WithFields(log.Fields{"request": "Version"})
	hash, err := doRequest()
	if err != nil {
		requestLogger.Warnf("error: %s", err)
		return nil, err
	}

	requestLogger.Debug("success")
	return hash, nil
}
