package mining

import "github.com/tclchiam/block_n_go/blockchain/block"

type Miner interface {
	MineBlock(header *block.Header) *block.Block
}
