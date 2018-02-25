package main

import (
	"github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/logger"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

var (
	accountLogger   = logger.Default
	interruptLogger = logger.Default
	nodeLogger      = logger.Default
	p2pLogger       = logger.Default
	rpcLogger       = logger.Default
	storageLogger   = logger.Default
)

func init() {
	accountLogger.SetLevel(logrus.InfoLevel)
	interruptLogger.SetLevel(logrus.InfoLevel)
	nodeLogger.SetLevel(logrus.WarnLevel)
	p2pLogger.SetLevel(logrus.WarnLevel)
	rpcLogger.SetLevel(logrus.InfoLevel)
	storageLogger.SetLevel(logrus.WarnLevel)

	account.UseLogger(accountLogger)
	interrupt.UseLogger(interruptLogger)
	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
