package node

import (
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/p2p"
)

type Node interface {
	AddPeer(address string) (*p2p.Peer, error)
	Addr() string
	ActivePeers() p2p.Peers
	Serve()
	Close() error
}

func StartNode(bc blockchain.Blockchain, config Config) (Node, error) {
	node, err := NewNode(bc, config)
	if err != nil {
		return nil, err
	}

	return WrapWithLogger(node), nil
}
