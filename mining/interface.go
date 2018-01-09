package mining

import "github.com/tclchiam/block_n_go/blockchain/entity"

type Miner interface {
	MineBlock(header *entity.BlockHeader) *entity.BlockSolution
}
