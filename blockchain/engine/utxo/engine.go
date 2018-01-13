package utxo

import "github.com/tclchiam/block_n_go/identity"

type Engine interface {
	FindUnspentOutputs(address *identity.Address) (*TransactionOutputSet, error)
}
