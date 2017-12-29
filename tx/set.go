package tx

import (
	"fmt"
	"strings"
)

type TransactionOutputSet struct {
	transactionsToOutputs map[string][]*Output
}

func NewTransactionSet() *TransactionOutputSet {
	return &TransactionOutputSet{make(map[string][]*Output, 0)}
}

func (set *TransactionOutputSet) Contains(transactionId string, output *Output) bool {
	if spentOutputIds, ok := set.transactionsToOutputs[transactionId]; ok {
		for _, spentOutput := range spentOutputIds {
			if spentOutput.IsEqual(output) {
				return true
			}
		}
	}
	return false
}

func (set *TransactionOutputSet) Add(transactionId string, output *Output) *TransactionOutputSet {
	outputs := set.transactionsToOutputs[transactionId]

	newTransactionsToOutputs := copyTransactionOutputs(set)
	newTransactionsToOutputs[transactionId] = append(outputs, output)

	return &TransactionOutputSet{
		transactionsToOutputs: newTransactionsToOutputs,
	}
}

func (set *TransactionOutputSet) Plus(other *TransactionOutputSet) *TransactionOutputSet {
	addToTxSet := func(res interface{}, transactionId string, output *Output) interface{} {
		return res.(*TransactionOutputSet).Add(transactionId, output)
	}

	return other.Reduce(set, addToTxSet).(*TransactionOutputSet)
}

func copyTransactionOutputs(set *TransactionOutputSet) map[string][]*Output {
	newTransactionsToOutputIds := make(map[string][]*Output, 0)
	for k, v := range set.transactionsToOutputs {
		newTransactionsToOutputIds[k] = v
	}
	return newTransactionsToOutputIds
}

func (set *TransactionOutputSet) Reduce(res interface{}, apply func(res interface{}, transactionId string, output *Output) interface{}) interface{} {
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

func (set *TransactionOutputSet) Filter(predicate func(transactionId string, output *Output) bool) *TransactionOutputSet {
	c := make(chan *TransactionOutputSet)

	go func() {
		newTransactionsToOutputIds := make(map[string][]*Output, 0)

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
