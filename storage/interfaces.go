package storage

import (
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type BlockReader interface {
	Head() (head *entity.Block, err error)

	Block(hash chainhash.Hash) (*entity.Block, error)

	SaveBlock(*entity.Block) error

	Close() error
}
