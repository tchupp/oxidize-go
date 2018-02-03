package node

import (
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/blockrpc"
	"github.com/tclchiam/block_n_go/p2p"
	"github.com/tclchiam/block_n_go/rpc"
)

type baseNode struct {
	p2p.PeerManager
	*rpc.Server

	bc blockchain.Blockchain
}

func NewNode(bc blockchain.Blockchain, server *rpc.Server) Node {
	return newNode(bc, server)
}

func newNode(bc blockchain.Blockchain, server *rpc.Server) *baseNode {
	server.RegisterSyncServer(blockrpc.NewSyncServer(bc))
	server.RegisterDiscoveryServer(p2p.NewDiscoveryServer(bc))

	return &baseNode{
		bc:          bc,
		PeerManager: p2p.NewPeerManager(),
		Server:      server,
	}
}

func (n *baseNode) AddPeer(address string) (*p2p.Peer, error) {
	peer, err := n.PeerManager.AddPeer(address)
	if err != nil {
		return nil, err
	}

	go startSyncHeadersFlow(peer, n.PeerManager, n.bc)

	return peer, nil
}
