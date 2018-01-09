package engine

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/encoding"
)

func MineBlock(transactions []*entity.Transaction, miner mining.Miner, blockRepository storage.BlockRepository) (*entity.Block, error) {
	headBlock, err := blockRepository.Head()
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		for index, input := range transaction.Inputs {
			verified := VerifySignature(input, transaction.Outputs, encoding.NewTransactionGobEncoder())
			if !verified {
				return nil, fmt.Errorf(TransactionInputHasBadSignatureMessage, transaction.ID, index)
			}
		}
	}

	newBlockHeader := entity.NewBlockHeader(headBlock.Index+1, headBlock.Hash, transactions)
	return miner.MineBlock(newBlockHeader), nil
}
