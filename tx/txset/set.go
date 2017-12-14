package txset

type TransactionSet struct {
	transactionsToOutputIds map[string][]int
}

func New() *TransactionSet {
	return &TransactionSet{make(map[string][]int, 0)}
}

func (spent *TransactionSet) Contains(transactionId string, outputId int) bool {
	if spentOutputIds, ok := spent.transactionsToOutputIds[transactionId]; ok {
		for _, spentOutputId := range spentOutputIds {
			if spentOutputId == outputId {
				return true
			}
		}
	}
	return false
}

func (spent *TransactionSet) Add(transactionId string, outputId int) *TransactionSet {
	outputIds := spent.transactionsToOutputIds[transactionId]

	newTransactionsToOutputIds := copyTransactionOutputs(spent)
	newTransactionsToOutputIds[transactionId] = append(outputIds, outputId)

	return &TransactionSet{transactionsToOutputIds: newTransactionsToOutputIds}
}

func copyTransactionOutputs(spent *TransactionSet) map[string][]int {
	newTransactionsToOutputIds := make(map[string][]int, 0)
	for k, v := range spent.transactionsToOutputIds {
		newTransactionsToOutputIds[k] = v
	}
	return newTransactionsToOutputIds
}
