package p2p

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type Peer struct {
	Address  string
	BestHash *entity.Hash
}

func (p *Peer) String() string {
	return fmt.Sprintf("{Address: %s, BestHash: %s}", p.Address, p.BestHash)
}

type Peers []*Peer

func (peers Peers) Add(peer *Peer) Peers {
	return append(peers, peer)
}

func (peers Peers) Remove(target *Peer) Peers {
	for i, peer := range peers {
		if peer.Address == target.Address {
			peers[i] = peers[0]
			return peers[1:]
		}
	}
	return append(Peers(nil), peers...)
}

func (peers Peers) ContainsAddress(address string) bool {
	for _, peer := range peers {
		if peer.Address == address {
			return true
		}
	}
	return false
}

func (peers Peers) Contains(target *Peer) bool {
	if peers.ContainsAddress(target.Address) {
		return true
	}
	return false
}
