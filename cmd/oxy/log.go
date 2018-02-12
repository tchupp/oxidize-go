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
	nodeLogger.SetLevel(logrus.WarnLevel)
	p2pLogger.SetLevel(logrus.WarnLevel)
	rpcLogger.SetLevel(logrus.WarnLevel)
	storageLogger.SetLevel(logrus.WarnLevel)

	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
