package mining

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type Miner interface {
	MineBlock(parent *entity.BlockHeader, transactions entity.Transactions) *entity.Block
}
