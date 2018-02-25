package account

import (
	"fmt"
	"strings"

	"github.com/tclchiam/oxidize-go/identity"
)

type Transaction struct {
	amount   uint64
	spender  *identity.Address
	receiver *identity.Address
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("{amount: %d, spender: %s, receiver: %s}", tx.amount, tx.spender, tx.receiver)
}

type Transactions []*Transaction

func (txs Transactions) Add(tx *Transaction) Transactions { return append(txs, tx) }
func (txs Transactions) String() string {
	if txs == nil {
		return "accounts.Transactions(nil)"
	}

	var lines []string
	for _, tx := range txs {
		lines = append(lines, tx.String())
	}

	return "[" + strings.Join(lines, ", ") + "]"
}
