package rpc

import (
	"google.golang.org/grpc/grpclog"

	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/logger"
)

var log = logger.Disabled

func UseLogger(logger *logrus.Entry) {
	log = logger
	grpclog.SetLogger(logger)
}
