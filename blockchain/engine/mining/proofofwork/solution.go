package proofofwork

import "github.com/tclchiam/block_n_go/blockchain/entity"

type BlockSolution struct {
	Nonce uint64
	Hash  *entity.Hash
}
