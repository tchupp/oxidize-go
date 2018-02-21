package account

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/identity"
)

type Transaction struct {
	amount uint64
	from   *identity.Address
	to     *identity.Address
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("{amount: %d, from: %s, to: %s}", tx.amount, tx.from, tx.to)
}

type Transactions []*Transaction

func (txs Transactions) Add(tx *Transaction) Transactions { return append(txs, tx) }
