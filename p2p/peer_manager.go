package p2p

import (
	"sync"

	"google.golang.org/grpc"
)

type PeerManager interface {
	AddPeer(address string) (*Peer, error)
	GetPeerConnection(peer *Peer) *grpc.ClientConn
	ActivePeers() Peers
}

type peerManager struct {
	ConnectionManager
	peers Peers
	lock  sync.RWMutex
}

func NewPeerManager() PeerManager {
	return newPeerManager(NewConnectionManager())
}

func newPeerManager(connectionManager ConnectionManager) *peerManager {
	return &peerManager{
		ConnectionManager: connectionManager,
		peers:             Peers{},
	}
}

func (m *peerManager) AddPeer(address string) (*Peer, error) {
	conn, err := m.ConnectionManager.OpenConnection(address)
	if err != nil {
		return nil, err
	}

	hash, err := NewDiscoveryClient(conn).Version()
	if err != nil {
		m.ConnectionManager.CloseConnection(address)
		return nil, err
	}

	peer := &Peer{
		Address:  address,
		BestHash: hash,
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.peers.Contains(peer) {
		m.peers = m.peers.Add(peer)
	}

	return peer, nil
}

func (m *peerManager) ActivePeers() Peers {
	peers := Peers{}
	for _, peer := range m.peers {
		if m.ConnectionManager.HasConnection(peer.Address) {
			peers = append(peers, peer)
		}
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.peers = peers
	return peers
}

func (m *peerManager) GetPeerConnection(peer *Peer) *grpc.ClientConn {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if !m.peers.Contains(peer) {
		return nil
	}

	return m.ConnectionManager.GetConnection(peer.Address)
}
