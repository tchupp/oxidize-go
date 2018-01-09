package utxo

type Engine interface {
	FindUnspentOutputs(address string) (*TransactionOutputSet, error)
}
