package main

import (
	"github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/logger"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

var (
	nodeLogger    = logger.Default
	p2pLogger     = logger.Default
	rpcLogger     = logger.Default
	storageLogger = logger.Default
)

func init() {
	nodeLogger.SetLevel(logrus.InfoLevel)
	p2pLogger.SetLevel(logrus.InfoLevel)
	rpcLogger.SetLevel(logrus.InfoLevel)
	storageLogger.SetLevel(logrus.InfoLevel)

	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
