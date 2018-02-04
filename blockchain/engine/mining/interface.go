package mining

import "github.com/tclchiam/block_n_go/blockchain/entity"

type Miner interface {
	MineBlock(parent *entity.BlockHeader, transactions entity.Transactions) *entity.Block
}
