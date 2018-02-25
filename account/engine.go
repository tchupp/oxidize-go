package account

import (
	"github.com/tclchiam/oxidize-go/account/utxo"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/closer"
	"github.com/tclchiam/oxidize-go/identity"
)

type Engine interface {
	Balance(address *identity.Address) (*Account, error)
	Transactions(address *identity.Address) (Transactions, error)

	Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error

	Close() error
}

type engine struct {
	bc      blockchain.Blockchain
	repo    *accountRepo
	indexer *chainIndexer
}

func NewEngine(bc blockchain.Blockchain) Engine {
	repo := NewAccountRepository()
	indexer := NewChainIndexer(bc, repo)

	return &engine{
		bc:      bc,
		repo:    repo,
		indexer: indexer,
	}
}

func (e *engine) Balance(address *identity.Address) (*Account, error) {
	return balance(address, utxo.NewCrawlerEngine(e.bc))
}

func (e *engine) Transactions(address *identity.Address) (Transactions, error) {
	account, err := e.repo.Account(address)
	if err != nil {
		return nil, err
	}
	return account.Transactions, nil
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

func (e *engine) Close() error {
	return closer.CloseMany(e.bc, e.indexer)
}
