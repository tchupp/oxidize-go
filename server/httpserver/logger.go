package httpserver

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/logger"
)

var log = logger.Disabled

func UseLogger(logger *logrus.Logger) {
	log = logger
}
