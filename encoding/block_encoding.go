package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type (
	blockData struct {
		Index        int
		PreviousHash chainhash.Hash
		Timestamp    int64
		Transactions []*txData
		Hash         chainhash.Hash
		Nonce        int
	}
)

func toBlockData(block *entity.Block) *blockData {
	var transactions []*txData
	for _, transaction := range block.Transactions {
		transactions = append(transactions, toTxData(transaction))
	}

	return &blockData{
		Index:        block.Index,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
		Transactions: transactions,
		Hash:         block.Hash,
		Nonce:        block.Nonce,
	}
}

func fromBlockData(block *blockData) (*entity.Block, error) {
	var transactions []*entity.Transaction
	for _, txData := range block.Transactions {
		transaction, err := fromTxData(txData)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return &entity.Block{
		Index:        block.Index,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
		Transactions: transactions,
		Hash:         block.Hash,
		Nonce:        block.Nonce,
	}, nil
}
