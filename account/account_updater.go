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
	for _, tx := range block.Transactions() {
		var spenderAddress *identity.Address
		if len(tx.Inputs) > 0 {
			spenderAddress = identity.FromPublicKey(tx.Inputs[0].PublicKey)
		}

		for _, out := range tx.Outputs {
			accountTx := &Transaction{
				amount:   out.Value,
				spender:  spenderAddress,
				receiver: identity.FromPublicKeyHash(out.PublicKeyHash),
			}

			p.repo.SaveTx(accountTx.spender, accountTx)
			p.repo.SaveTx(accountTx.receiver, accountTx)
		}
	}
}
