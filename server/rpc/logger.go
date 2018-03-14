package rpc

import (
	"google.golang.org/grpc/grpclog"

	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/oxylogger"
)

var log = oxylogger.Disabled
var grpcLogger = logrus.NewEntry(log)

type rpcLogger struct {
	*logrus.Logger
}

func (rpcLogger) V(l int) bool {
	return true
}

func UseLogger(logger *logrus.Logger) {
	log = logger
	grpcLogger = logrus.NewEntry(log)
	grpclog.SetLoggerV2(rpcLogger{Logger: logger})
}
