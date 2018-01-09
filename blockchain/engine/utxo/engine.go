package utxo

import "github.com/tclchiam/block_n_go/blockchain/tx"

type Engine interface {
	FindUnspentOutputs(address string) (*tx.TransactionOutputSet, error)
}
