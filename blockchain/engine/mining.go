package engine

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/engine/txsigning"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/mining"
)

func MineBlock(transactions entity.Transactions, miner mining.Miner, repository entity.BlockRepository) (*entity.Block, error) {
	parent, err := repository.BestBlock()
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		for index, input := range transaction.Inputs {
			if verified := txsigning.VerifySignature(input, transaction.Outputs, encoding.TransactionProtoEncoder()); !verified {
				return nil, fmt.Errorf(TransactionInputHasBadSignatureMessage, transaction.ID, index)
			}
		}
	}

	return miner.MineBlock(parent.Header(), transactions), nil
}
