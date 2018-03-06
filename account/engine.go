package account

import (
	"time"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/closer"
	"github.com/tclchiam/oxidize-go/identity"
)

type Engine interface {
	Account(address *identity.Address) (*Account, error)

	SpendableOutputs(*identity.Address) (*utxo.OutputSet, error)
	Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error
	ProposeTransaction(*entity.Transaction) error

	Close() error
}

type engine struct {
	bc      blockchain.Blockchain
	repo    *accountRepo
	indexer *chainIndexer
}

func NewEngine(bc blockchain.Blockchain) Engine {
	repo := NewAccountRepository()
	indexer := NewChainIndexer(bc, NewAccountUpdater(repo))

	return &engine{
		bc:      bc,
		repo:    repo,
		indexer: indexer,
	}
}

func (e *engine) Account(address *identity.Address) (*Account, error) {
	<-e.waitForIndexer()

	account, err := e.repo.Account(address)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (e *engine) SpendableOutputs(address *identity.Address) (*utxo.OutputSet, error) {
	return e.bc.SpendableOutputs(address)
}

func (e *engine) Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error {
	spendableOutputs, err := e.bc.SpendableOutputs(spender.Address())
	if err != nil {
		return err
	}

	expenseTransaction, err := buildExpenseTransaction(spender, receiver, expense, spendableOutputs)
	if err != nil {
		return err
	}

	newBlock, err := e.bc.MineBlock(entity.Transactions{expenseTransaction})
	if err != nil {
		return err
	}
	return e.bc.SaveBlock(newBlock)
}

func (e *engine) ProposeTransaction(tx *entity.Transaction) error {
	newBlock, err := e.bc.MineBlock(entity.Transactions{tx})
	if err != nil {
		return err
	}
	return e.bc.SaveBlock(newBlock)
}

func (e *engine) Close() error {
	return closer.CloseMany(e.bc, e.indexer)
}

func (e *engine) waitForIndexer() <-chan struct{} {
	c := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Millisecond) // give the indexer a moment to actually update

		status := e.indexer.Status()
		for status != Idle && status != Done {
			time.Sleep(1 * time.Millisecond) // give the indexer a moment to actually update
			status = e.indexer.Status()
		}
		close(c)
	}()

	return c
}
