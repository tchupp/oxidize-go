package logger

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	Disabled *logrus.Logger
	Default  *logrus.Logger
)

func init() {
	Disabled = logrus.New()
	Disabled.Out = ioutil.Discard

	Default = logrus.New()
	Default.SetLevel(logrus.InfoLevel)
}
