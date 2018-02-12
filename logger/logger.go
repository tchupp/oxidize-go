package logger

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	Disabled *logrus.Logger
)

func init() {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	Disabled = logger
}
