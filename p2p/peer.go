package p2p

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type Peer struct {
	Address string
	Head    *entity.Hash
}

func (p *Peer) String() string {
	return fmt.Sprintf("{Address: %s, Head: %s}", p.Address, p.Head)
}

type Peers []*Peer

func (peers Peers) Add(peer *Peer) Peers {
	return append(peers, peer)
}
