package rpc

import (
	log "github.com/sirupsen/logrus"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
)

var (
	logrusEntry = log.NewEntry(log.StandardLogger())
)

func init() {
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
}
