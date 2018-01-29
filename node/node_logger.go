package node

import (
	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/block_n_go/p2p"
)

type loggingNodeDecorator struct {
	Node
}

func WrapWithLogger(inner Node) Node {
	return &loggingNodeDecorator{Node: inner}
}

func (n *loggingNodeDecorator) AddPeer(address string) (*p2p.Peer, error) {
	peer, err := n.Node.AddPeer(address)
	if err != nil {
		log.WithError(err).Warnf("unable to add peer '%s'", address)
		return nil, err
	}

	log.Debugf("added peer: %s", address)
	return peer, nil
}
