package storage

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type BlockRepository interface {
	Head() (head *entity.Block, err error)

	Block(hash *entity.Hash) (*entity.Block, error)

	SaveBlock(*entity.Block) error

	Close() error
}
