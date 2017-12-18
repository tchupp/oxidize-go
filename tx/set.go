package tx

import (
	"fmt"
	"strings"
)

type TransactionSet struct {
	transactionsToOutputs map[string][]*Output
}

func NewTransactionSet() *TransactionSet {
	return &TransactionSet{make(map[string][]*Output, 0)}
}

func (set *TransactionSet) Contains(transactionId string, output *Output) bool {
	if spentOutputIds, ok := set.transactionsToOutputs[transactionId]; ok {
		for _, spentOutput := range spentOutputIds {
			if *spentOutput == *output {
				return true
			}
		}
	}
	return false
}

func (set *TransactionSet) Add(transactionId string, output *Output) *TransactionSet {
	outputs := set.transactionsToOutputs[transactionId]

	newTransactionsToOutputs := copyTransactionOutputs(set)
	newTransactionsToOutputs[transactionId] = append(outputs, output)

	return &TransactionSet{
		transactionsToOutputs: newTransactionsToOutputs,
	}
}

func (set *TransactionSet) Plus(other *TransactionSet) *TransactionSet {
	addToTxSet := func(res interface{}, transactionId string, output *Output) interface{} {
		return res.(*TransactionSet).Add(transactionId, output)
	}

	return other.Reduce(set, addToTxSet).(*TransactionSet)
}

func copyTransactionOutputs(set *TransactionSet) map[string][]*Output {
	newTransactionsToOutputIds := make(map[string][]*Output, 0)
	for k, v := range set.transactionsToOutputs {
		newTransactionsToOutputIds[k] = v
	}
	return newTransactionsToOutputIds
}

func (set *TransactionSet) Reduce(res interface{}, apply func(res interface{}, transactionId string, output *Output) interface{}) interface{} {
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

func (set *TransactionSet) Filter(predicate func(transactionId string, output *Output) bool) *TransactionSet {
	c := make(chan *TransactionSet)

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
		c <- &TransactionSet{newTransactionsToOutputIds}
	}()
	return <-c
}

func (set *TransactionSet) String() string {
	var lines []string

	for transactionId, outputs := range set.transactionsToOutputs {
		lines = append(lines, fmt.Sprintf("Transaction %x:", transactionId))
		for _, output := range outputs {
			lines = append(lines, output.string()...)
		}
	}

	return strings.Join(lines, "")
}
