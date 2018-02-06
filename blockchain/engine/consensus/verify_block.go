package consensus

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func VerifyBlock(block *entity.Block) error {
	if err := VerifyHeader(block.Header()); err != nil {
		return err
	}
	// TODO verify transactions
	if transactionsHash := mining.CalculateTransactionsHash(block.Transactions()); !transactionsHash.IsEqual(block.TransactionsHash()) {
		return errInvalidTxHash
	}

	return nil
}
