package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

const systemField = "system"

var (
	nodeLogger    = logrus.WithField(systemField, "node")
	rpcLogger     = logrus.WithField(systemField, "rpc")
	storageLogger = logrus.WithField(systemField, "storage")
)

func init() {
	node.UseLogger(nodeLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
