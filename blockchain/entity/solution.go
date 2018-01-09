package entity

import "github.com/tclchiam/block_n_go/blockchain/chainhash"

type BlockSolution struct {
	Nonce uint64
	Hash  *chainhash.Hash
}
