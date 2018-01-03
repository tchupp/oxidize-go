package blockchain

import "github.com/tclchiam/block_n_go/chainhash"

type Repository interface {
	Head() (head *Block, err error)

	Block(hash chainhash.Hash) (*Block, error)

	SaveBlock(*Block) error
}
