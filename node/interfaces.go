package node

import (
	"github.com/tclchiam/block_n_go/p2p"
)

type Node interface {
	AddPeer(address string) error
	ActivePeers() p2p.Peers
}
