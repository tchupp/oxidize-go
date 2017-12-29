package tx

import (
	"fmt"
	"strings"
)

type TransactionOutputSet struct {
	transactionsToOutputs map[*Transaction][]*Output
}

func NewTransactionSet() *TransactionOutputSet {
	return &TransactionOutputSet{make(map[*Transaction][]*Output, 0)}
}

func (set *TransactionOutputSet) Len() int {
	return len(set.transactionsToOutputs)
}

func (set *TransactionOutputSet) Add(transaction *Transaction, output *Output) *TransactionOutputSet {
	outputs := set.transactionsToOutputs[transaction]

	newTransactionsToOutputs := copyTransactionOutputs(set)
	newTransactionsToOutputs[transaction] = append(outputs, output)

	return &TransactionOutputSet{
		transactionsToOutputs: newTransactionsToOutputs,
	}
}

func (set *TransactionOutputSet) Plus(other *TransactionOutputSet) *TransactionOutputSet {
	addToTxSet := func(res interface{}, transaction *Transaction, output *Output) interface{} {
		return res.(*TransactionOutputSet).Add(transaction, output)
	}

	return other.Reduce(set, addToTxSet).(*TransactionOutputSet)
}

func (set *TransactionOutputSet) Reduce(res interface{}, apply func(res interface{}, transaction *Transaction, output *Output) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for transaction, outputs := range set.transactionsToOutputs {
			for _, output := range outputs {
				res = apply(res, transaction, output)
			}
		}
		c <- res
	}()
	return <-c
}

func (set *TransactionOutputSet) Filter(predicate func(transaction *Transaction, output *Output) bool) *TransactionOutputSet {
	c := make(chan *TransactionOutputSet)

	go func() {
		newTransactionsToOutputIds := make(map[*Transaction][]*Output, 0)

		for transaction, outputs := range set.transactionsToOutputs {
			for _, output := range outputs {
				if predicate(transaction, output) {
					outputs := newTransactionsToOutputIds[transaction]
					newTransactionsToOutputIds[transaction] = append(outputs, output)
				}
			}
		}
		c <- &TransactionOutputSet{newTransactionsToOutputIds}
	}()
	return <-c
}

func (set *TransactionOutputSet) String() string {
	var lines []string

	for transactionId, outputs := range set.transactionsToOutputs {
		lines = append(lines, fmt.Sprintf("Transaction %x:", transactionId))
		for _, output := range outputs {
			lines = append(lines, output.string()...)
		}
	}

	return strings.Join(lines, "")
}

func copyTransactionOutputs(set *TransactionOutputSet) map[*Transaction][]*Output {
	newTransactionsToOutputIds := make(map[*Transaction][]*Output, 0)
	for k, v := range set.transactionsToOutputs {
		newTransactionsToOutputIds[k] = v
	}
	return newTransactionsToOutputIds
}
