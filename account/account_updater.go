package account

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type accountUpdater struct {
	repo *accountRepo
}

func NewAccountUpdater(repo *accountRepo) BlockProcessor {
	return &accountUpdater{repo: repo}
}

func (p *accountUpdater) Process(block *entity.Block) {
	var accountTxs Transactions
	for _, tx := range block.Transactions() {
		spenderAddress := findSpenderAddress(tx)

		for _, out := range tx.Outputs {
			accountTx := &Transaction{
				amount:   out.Value,
				spender:  spenderAddress,
				receiver: identity.FromPublicKeyHash(out.PublicKeyHash),
			}
			accountTxs = accountTxs.Add(accountTx)
		}
	}

	p.repo.SaveTxs(accountTxs)
}

// We trust that every transaction with one or more inputs has one sender, rewards have no sender
func findSpenderAddress(tx *entity.Transaction) *identity.Address {
	var spenderAddress *identity.Address
	if len(tx.Inputs) > 0 {
		spenderAddress = identity.FromPublicKey(tx.Inputs[0].PublicKey)
	}
	return spenderAddress
}
