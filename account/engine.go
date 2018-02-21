package account

import (
	"bytes"

	"github.com/tclchiam/oxidize-go/account/utxo"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type Engine interface {
	Balance(address *identity.Address) (*Account, error)
	Transactions(address *identity.Address) (Transactions, error)

	Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error
}

type engine struct {
	bc blockchain.Blockchain
}

func NewEngine(bc blockchain.Blockchain) Engine {
	return &engine{
		bc: bc,
	}
}

func (e *engine) Balance(address *identity.Address) (*Account, error) {
	return balance(address, utxo.NewCrawlerEngine(e.bc))
}

func (e *engine) Transactions(address *identity.Address) (Transactions, error) {
	entityTxs := entity.Transactions{}
	e.bc.ForEachBlock(func(block *entity.Block) {
		filteredTxs := block.Transactions().
			FilterReward(false).
			Filter(func(tx *entity.Transaction) bool {
				for _, in := range tx.Inputs {
					if identity.FromPublicKey(in.PublicKey).IsEqual(address) {
						return true
					}
				}

				for _, out := range tx.Outputs {
					if bytes.Compare(out.PublicKeyHash, address.PublicKeyHash()) == 0 {
						return true
					}
				}

				return false
			})

		entityTxs = entityTxs.Add(filteredTxs...)
	})

	accountTxs := Transactions{}
	for _, tx := range entityTxs {
		for _, out := range tx.Outputs {
			accountTxs = accountTxs.Add(
				&Transaction{
					amount: out.Value,
					from:   identity.FromPublicKey(tx.Inputs[0].PublicKey),
					to:     identity.FromPublicKeyHash(out.PublicKeyHash),
				},
			)
		}
	}
	return accountTxs, nil
}

func (e *engine) Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error {
	expenseTransaction, err := buildExpenseTransaction(spender, receiver, expense, utxo.NewCrawlerEngine(e.bc))
	if err != nil {
		return err
	}

	newBlock, err := e.bc.MineBlock(entity.Transactions{expenseTransaction})
	if err != nil {
		return err
	}
	return e.bc.SaveBlock(newBlock)
}
