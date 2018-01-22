package node

import (
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/p2p"
	"github.com/tclchiam/block_n_go/rpc"
)

type baseNode struct {
	connectionManager *rpc.ConnectionManager
	peers             p2p.Peers
}

func NewNode(bc *blockchain.Blockchain, connectionManager *rpc.ConnectionManager) Node {
	return &baseNode{
		connectionManager: connectionManager,
		peers:             p2p.Peers{},
	}
}

func (n *baseNode) AddPeer(address string) error {
	conn, err := n.connectionManager.OpenConnection(address)
	if err != nil {
		return err
	}

	client := p2p.NewDiscoveryClient(conn)
	hash, err := client.Version()
	if err != nil {
		conn.Close()
		return err
	}

	peer := &p2p.Peer{
		Address: address,
		Head:    hash,
	}

	n.peers = n.peers.Add(peer)
	return nil
}

func (n *baseNode) ActivePeers() p2p.Peers {
	return append(p2p.Peers(nil), n.peers...)
}
