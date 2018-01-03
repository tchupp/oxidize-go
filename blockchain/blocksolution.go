package blockchain

import "github.com/tclchiam/block_n_go/chainhash"

type BlockSolution struct {
	Nonce int
	Hash  chainhash.Hash
}
