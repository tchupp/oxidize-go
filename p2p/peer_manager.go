package p2p

import (
	"sync"
	"time"

	"github.com/tclchiam/block_n_go/rpc"
)

type PeerManager interface {
	rpc.ConnectionManager
	AddPeer(address string) error
	ActivePeers() Peers
}

type peerManager struct {
	rpc.ConnectionManager
	peers             Peers
	lock              sync.RWMutex
}

func NewPeerManager() PeerManager {
	return newPeerManager(rpc.NewConnectionManager())
}

func newPeerManager(connectionManager rpc.ConnectionManager) *peerManager {
	return &peerManager{
		ConnectionManager: connectionManager,
		peers:             Peers{},
	}
}

func (m *peerManager) AddPeer(address string) error {
	conn, err := m.ConnectionManager.OpenConnection(address)
	if err != nil {
		return err
	}

	hash, err := NewDiscoveryClient(conn).Version()
	if err != nil {
		m.ConnectionManager.CloseConnection(address)
		return err
	}

	peer := &Peer{
		Address: address,
		Head:    hash,
	}

	m.addPeer(peer)
	return nil
}
func (m *peerManager) addPeer(peer *Peer) {
	m.lock.Lock()
	if !m.peers.Contains(peer) {
		m.peers = m.peers.Add(peer)
		go m.peerMonitor(peer)
	}
	m.lock.Unlock()
}

func (m *peerManager) removePeer(target *Peer) {
	if target == nil {
		return
	}

	m.lock.Lock()
	m.ConnectionManager.CloseConnection(target.Address)
	m.peers = m.peers.Remove(target)
	m.lock.Unlock()
}

func (m *peerManager) ActivePeers() Peers {
	return append(Peers(nil), m.peers...)
}

// Expected to be run in a goroutine
func (m *peerManager) peerMonitor(peer *Peer) {
	for {
		conn := m.ConnectionManager.GetConnection(peer.Address)
		if conn == nil {
			m.removePeer(peer)
			return
		}

		client := NewDiscoveryClient(conn)
		if err := client.Ping(); err != nil {
			m.removePeer(peer)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
