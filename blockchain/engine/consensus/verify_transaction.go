package consensus

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

const (
	transactionInputHasBadSignatureMessage = "transaction '%s' has input '%d' with bad signature"
)

func VerifyTransaction(transactions entity.Transactions) error {
	var result *multierror.Error
	for _, transaction := range transactions {
		for index, input := range transaction.Inputs {
			if verified := txsigning.VerifySignature(input, transaction.Outputs); !verified {
				result = multierror.Append(result, fmt.Errorf(transactionInputHasBadSignatureMessage, transaction.ID, index))
			}
		}
	}

	return result.ErrorOrNil()
}
