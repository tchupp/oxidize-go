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
		log.Warnf("unable to add peer '%s': %s", address, err)
		return err
	}
	log.Debugf("added peer: %s", address)
	return nil
}

func (n *loggingNodeDecorator) ActivePeers() p2p.Peers {
	return n.inner.ActivePeers()
}
