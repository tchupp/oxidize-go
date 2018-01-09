package engine

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/engine/txsigning"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/encoding"
)

func MineBlock(transactions entity.Transactions, miner mining.Miner, repository storage.BlockRepository) (*entity.Block, error) {
	headBlock, err := repository.Head()
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		for index, input := range transaction.Inputs {
			verified := txsigning.VerifySignature(input, transaction.Outputs, encoding.NewTransactionGobEncoder())
			if !verified {
				return nil, fmt.Errorf(TransactionInputHasBadSignatureMessage, transaction.ID, index)
			}
		}
	}

	newBlockHeader := entity.NewBlockHeaderNow(headBlock.Index()+1, headBlock.Hash(), transactions)
	solution := miner.MineBlock(newBlockHeader)
	return entity.NewBlock(newBlockHeader, solution, transactions), nil
}
