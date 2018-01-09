package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type (
	blockData struct {
		Index        int             `json:"index"`
		PreviousHash *chainhash.Hash `json:"previous_hash"`
		Timestamp    int64           `json:"timestamp"`
		Transactions []*txData       `json:"transactions"`
		Hash         *chainhash.Hash `json:"hash"`
		Nonce        int             `json:"nonce"`
	}
)

func toBlockData(block *entity.Block) *blockData {
	var transactions []*txData
	for _, transaction := range block.Transactions() {
		transactions = append(transactions, toTxData(transaction))
	}

	return &blockData{
		Index:        block.Index(),
		PreviousHash: block.PreviousHash(),
		Timestamp:    block.Timestamp(),
		Transactions: transactions,
		Hash:         block.Hash(),
		Nonce:        block.Nonce(),
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

	return entity.NewBlock(
		entity.NewBlockHeader(block.Index, block.PreviousHash, transactions, block.Timestamp),
		&entity.BlockSolution{Hash: block.Hash, Nonce: block.Nonce},
		transactions,
	), nil
}
