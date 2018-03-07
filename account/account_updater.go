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

func (updater *accountUpdater) Process(block *entity.Block) {
	var updates []Update
	for _, tx := range block.Transactions() {
		for _, input := range tx.Inputs {
			output := input.OutputReference.Output

			updates = append(updates, &spendUpdate{
				address: identity.FromPublicKeyHash(output.PublicKeyHash),
				amount:  output.Value,
			})
		}

		for _, output := range tx.Outputs {
			updates = append(updates, &receiveUpdate{
				address: identity.FromPublicKeyHash(output.PublicKeyHash),
				amount:  output.Value,
			})
		}
	}

	updater.repo.ProcessUpdates(updates)
}
