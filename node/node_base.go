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

	bc *blockchain.Blockchain
}

func NewNode(bc *blockchain.Blockchain, server *rpc.Server) Node {
	node := &baseNode{
		bc:          bc,
		PeerManager: p2p.NewPeerManager(),
		Server:      server,
	}

	server.RegisterSyncServer(blockrpc.NewSyncServer(bc))
	server.RegisterDiscoveryServer(p2p.NewDiscoveryServer(bc))

	return node
}
