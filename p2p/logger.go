package p2p

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/logger"
)

var log = logger.Disabled

func UseLogger(logger *logrus.Entry) {
	log = logger
}
