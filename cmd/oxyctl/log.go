package main

import (
	"github.com/sirupsen/logrus"

	account "github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/logger"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

func init() {
	accountLogger := logger.NewLogger(logrus.InfoLevel)
	interruptLogger := logger.NewLogger(logrus.InfoLevel)
	nodeLogger := logger.NewLogger(logrus.InfoLevel)
	p2pLogger := logger.NewLogger(logrus.InfoLevel)
	rpcLogger := logger.NewLogger(logrus.InfoLevel)
	storageLogger := logger.NewLogger(logrus.InfoLevel)

	account.UseLogger(accountLogger)
	interrupt.UseLogger(interruptLogger)
	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
