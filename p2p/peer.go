package p2p

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type Peer struct {
	Address string
	Head    *entity.Hash
}

type Peers []*Peer

func (peers Peers) Add(peer *Peer) Peers {
	return append(peers, peer)
}
