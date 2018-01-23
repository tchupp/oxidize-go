package node

import (
	log "github.com/sirupsen/logrus"
	"github.com/tclchiam/block_n_go/p2p"
)

type loggingNodeDecorator struct {
	inner Node
}

func WrapWithLogger(inner Node) Node {
	return &loggingNodeDecorator{inner: inner}
}

func (n *loggingNodeDecorator) AddPeer(address string) error {
	if err := n.inner.AddPeer(address); err != nil {
		log.WithError(err).Warnf("unable to add peer '%s'", address)
		return err
	}
	log.Debugf("added peer: %s", address)
	return nil
}

func (n *loggingNodeDecorator) ActivePeers() p2p.Peers {
	return n.inner.ActivePeers()
}

func (n *loggingNodeDecorator) Serve() {
	n.inner.Serve()
}

func (n *loggingNodeDecorator) Shutdown() error {
	if err := n.inner.Shutdown(); err != nil {
		log.WithError(err).Warn("error shutting down server")
	}
	return nil
}
