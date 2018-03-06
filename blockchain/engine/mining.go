package engine

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine/consensus"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func MineBlock(transactions entity.Transactions, miner mining.Miner, reader entity.ChainReader) (*entity.Block, error) {
	if err := consensus.VerifyTransaction(transactions); err != nil {
		return nil, err
	}

	parent, err := reader.BestBlock()
	if err != nil {
		return nil, err
	}
	return miner.MineBlock(parent.Header(), transactions), nil
}
