package p2p

import (
	"sync"

	"github.com/tclchiam/block_n_go/rpc"
)

type PeerManager interface {
	AddPeer(address string) error
	ActivePeers() Peers
}

type peerManager struct {
	connectionManager rpc.ConnectionManager
	peers             Peers
	lock              sync.RWMutex
}

func NewPeerManager() PeerManager {
	return newPeerManager(rpc.NewConnectionManager())
}

func newPeerManager(connectionManager rpc.ConnectionManager) *peerManager {
	return &peerManager{
		connectionManager: connectionManager,
		peers:             Peers{},
	}
}

func (m *peerManager) AddPeer(address string) error {
	conn, err := m.connectionManager.OpenConnection(address)
	if err != nil {
		return err
	}

	client := NewDiscoveryClient(conn)
	hash, err := client.Version()
	if err != nil {
		conn.Close()
		return err
	}

	peer := &Peer{
		Address: address,
		Head:    hash,
	}

	m.lock.Lock()
	if !m.peers.ContainsAddress(peer.Address) {
		m.peers = m.peers.Add(peer)
	}
	m.lock.Unlock()
	return nil
}

func (m *peerManager) ActivePeers() Peers {
	return append(Peers(nil), m.peers...)
}
