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

			address := identity.FromPublicKeyHash(output.PublicKeyHash)
			updates = append(updates, &spendUpdate{
				address: address,
				amount:  output.Value,
			})
			updates = append(updates, &txUpdate{
				address: address,
				tx:      tx,
			})
		}

		for _, output := range tx.Outputs {
			address := identity.FromPublicKeyHash(output.PublicKeyHash)

			updates = append(updates, &receiveUpdate{
				address: address,
				amount:  output.Value,
			})
			updates = append(updates, &txUpdate{
				address: address,
				tx:      tx,
			})
		}
	}

	updater.repo.ProcessUpdates(updates)
}
