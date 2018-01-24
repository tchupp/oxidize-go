package node

import (
	log "github.com/sirupsen/logrus"
)

type loggingNodeDecorator struct {
	Node
}

func WrapWithLogger(inner Node) Node {
	return &loggingNodeDecorator{Node: inner}
}

func (n *loggingNodeDecorator) AddPeer(address string) error {
	if err := n.Node.AddPeer(address); err != nil {
		log.WithError(err).Warnf("unable to add peer '%s'", address)
		return err
	}
	log.Debugf("added peer: %s", address)
	return nil
}

func (n *loggingNodeDecorator) Shutdown() error {
	if err := n.Node.Shutdown(); err != nil {
		log.WithError(err).Warn("error shutting down server")
	}
	return nil
}
