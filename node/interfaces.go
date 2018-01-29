package node

import (
	"github.com/tclchiam/block_n_go/p2p"
)

type Node interface {
	AddPeer(address string) (*p2p.Peer, error)
	ActivePeers() p2p.Peers
	Serve()
	Shutdown()
}
