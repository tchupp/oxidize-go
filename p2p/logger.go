package p2p

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/oxylogger"
)

var log = oxylogger.Disabled
var grpcLogger = logrus.NewEntry(log)

func UseLogger(logger *logrus.Logger) {
	log = logger
	grpcLogger = logrus.NewEntry(log)
}
