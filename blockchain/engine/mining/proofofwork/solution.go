package proofofwork

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type BlockSolution struct {
	Nonce uint64
	Hash  *entity.Hash
}
