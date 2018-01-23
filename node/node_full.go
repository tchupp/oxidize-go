package node

import (
	"sync"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/p2p"
	"github.com/tclchiam/block_n_go/rpc"
	"github.com/tclchiam/block_n_go/blockchain/blockrpc"
)

type baseNode struct {
	bc                *blockchain.Blockchain
	connectionManager *rpc.ConnectionManager
	peers             p2p.Peers
	server            *rpc.Server
	lock              sync.RWMutex
}

func NewNode(bc *blockchain.Blockchain, connectionManager *rpc.ConnectionManager, server *rpc.Server) Node {
	node := &baseNode{
		bc:                bc,
		connectionManager: connectionManager,
		peers:             p2p.Peers{},
		server:            server,
	}

	server.RegisterSyncServer(blockrpc.NewSyncServer(bc))
	server.RegisterDiscoveryServer(p2p.NewDiscoveryServer(bc))

	return node
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

	n.lock.Lock()
	if !n.peers.ContainsAddress(peer.Address) {
		n.peers = n.peers.Add(peer)
	}
	n.lock.Unlock()
	return nil
}

func (n *baseNode) ActivePeers() p2p.Peers {
	return append(p2p.Peers(nil), n.peers...)
}

func (n *baseNode) Serve() {
	n.server.Serve()
}

func (n *baseNode) Shutdown() error {
	return n.server.Shutdown()
}
