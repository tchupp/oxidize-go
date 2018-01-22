package rpc

import (
	"golang.org/x/net/context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

func LoggerFromContext(ctx context.Context) *logrus.Entry {
	if p, ok := peer.FromContext(ctx); ok {
		return logrus.WithField("remote", p.Addr)
	}
	return logrus.WithFields(logrus.Fields{})
}
