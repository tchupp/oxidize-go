package main

import (
	"github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/logger"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/server/httpserver"
	"github.com/tclchiam/oxidize-go/server/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

func init() {
	accountLogger := logger.NewLogger(logrus.InfoLevel)
	httpLogger := logger.NewLogger(logrus.InfoLevel)
	interruptLogger := logger.NewLogger(logrus.InfoLevel)
	nodeLogger := logger.NewLogger(logrus.WarnLevel)
	p2pLogger := logger.NewLogger(logrus.WarnLevel)
	rpcLogger := logger.NewLogger(logrus.WarnLevel)
	storageLogger := logger.NewLogger(logrus.WarnLevel)

	account.UseLogger(accountLogger)
	httpserver.UseLogger(httpLogger)
	interrupt.UseLogger(interruptLogger)
	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
