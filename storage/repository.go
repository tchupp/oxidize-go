package storage

import (
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/block"
)

type BlockReader interface {
	Head() (head *block.Block, err error)

	Block(hash chainhash.Hash) (*block.Block, error)

	SaveBlock(*block.Block) error
}
