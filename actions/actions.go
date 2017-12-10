package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
)

type Action interface {
	Execute(db *bolt.DB) (*blockchain.Blockchain, error)
}
