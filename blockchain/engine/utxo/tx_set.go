package utxo

import "github.com/tclchiam/block_n_go/blockchain/entity"

type TransactionOutputSet struct {
	transactionsToOutputs map[*entity.Transaction][]*entity.Output
}

func NewTransactionSet() *TransactionOutputSet {
	return &TransactionOutputSet{make(map[*entity.Transaction][]*entity.Output, 0)}
}

func (set *TransactionOutputSet) Add(transaction *entity.Transaction, output *entity.Output) *TransactionOutputSet {
	outputs := set.transactionsToOutputs[transaction]

	newTransactionsToOutputs := copyTransactionOutputs(set)
	newTransactionsToOutputs[transaction] = append(outputs, output)

	return &TransactionOutputSet{
		transactionsToOutputs: newTransactionsToOutputs,
	}
}

func (set *TransactionOutputSet) Plus(other *TransactionOutputSet) *TransactionOutputSet {
	addToTxSet := func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{} {
		return res.(*TransactionOutputSet).Add(transaction, output)
	}

	return other.Reduce(set, addToTxSet).(*TransactionOutputSet)
}

func (set *TransactionOutputSet) Reduce(res interface{}, apply func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{}) interface{} {
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

func (set *TransactionOutputSet) Filter(predicate func(transaction *entity.Transaction, output *entity.Output) bool) *TransactionOutputSet {
	c := make(chan *TransactionOutputSet)

	go func() {
		newTransactionsToOutputIds := make(map[*entity.Transaction][]*entity.Output, 0)

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

func copyTransactionOutputs(set *TransactionOutputSet) map[*entity.Transaction][]*entity.Output {
	newTransactionsToOutputIds := make(map[*entity.Transaction][]*entity.Output, 0)
	for k, v := range set.transactionsToOutputs {
		newTransactionsToOutputIds[k] = v
	}
	return newTransactionsToOutputIds
}
