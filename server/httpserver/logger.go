package httpserver

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/oxylogger"
)

var log = oxylogger.Disabled

func UseLogger(logger *logrus.Logger) {
	log = logger
}
