package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tclchiam/block_n_go/node"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/storage"
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
