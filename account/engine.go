package account

import (
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/account/utxo"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type Engine interface {
	Balance(address *identity.Address) (*Account, error)
	Transactions(address *identity.Address) (Transactions, error)

	Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error

	Close() error
}

type engine struct {
	bc        blockchain.Blockchain
	repo      *accountRepo
	processor *blockProcessor
}

func NewEngine(bc blockchain.Blockchain) Engine {
	repo := NewAccountRepository()
	processor := NewBlockProcessor(bc, repo)
	processor.Start()

	return &engine{
		bc:        bc,
		repo:      repo,
		processor: processor,
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
	closers := []io.Closer{e.bc, e.processor}

	var result *multierror.Error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result.ErrorOrNil()
}
