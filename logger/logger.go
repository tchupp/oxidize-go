package logger

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	Disabled *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	Disabled = logger.WithFields(logrus.Fields{})
}
