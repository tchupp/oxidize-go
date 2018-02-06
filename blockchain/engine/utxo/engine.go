package utxo

import "github.com/tclchiam/oxidize-go/identity"

type Engine interface {
	FindUnspentOutputs(spender *identity.Identity) (*TransactionOutputSet, error)
}
