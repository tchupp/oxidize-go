package account

import (
	"sync"

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
	bc   blockchain.Blockchain
	sub  blockchain.Subscription
	lock sync.RWMutex
	repo *accountRepo
}

func NewEngine(bc blockchain.Blockchain) Engine {
	channel := make(chan blockchain.Event, 10)
	engine := &engine{
		bc:   bc,
		sub:  bc.Subscribe(channel),
		repo: NewAccountRepository(),
	}

	go engine.indexChain(channel)

	return engine
}

func (e *engine) Balance(address *identity.Address) (*Account, error) {
	return balance(address, utxo.NewCrawlerEngine(e.bc))
}

func (e *engine) Transactions(address *identity.Address) (Transactions, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

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
	e.sub.Unsubscribe()
	return e.bc.Close()
}

func (e *engine) indexChain(c <-chan blockchain.Event) {
	var (
		currentIndex = uint64(0)
	)

	for {
		log.Debugf("resuming indexing with index '%d'...", currentIndex)

		currentIndex = e.processNewBlocks(currentIndex)

		log.Debugf("caught up with index '%d', waiting for event...", currentIndex)
		waitForBlockEvent(c)
	}
}

func (e *engine) processNewBlocks(currentIndex uint64) uint64 {
	e.lock.Lock()
	defer e.lock.Unlock()

	bestBlock, err := e.bc.BestBlock()
	for ; err == nil && currentIndex <= bestBlock.Index(); currentIndex += 1 {
		block, err := e.bc.BlockByIndex(currentIndex)
		if err != nil {
			log.Panicf("error reading block for indexing: %s")
		}

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

				e.repo.SaveTx(accountTx.spender, accountTx)
				e.repo.SaveTx(accountTx.receiver, accountTx)
			}
		}
	}

	return currentIndex
}

func waitForBlockEvent(c <-chan blockchain.Event) {
	for event, ok := <-c; ; {
		if !ok {
			log.Debug("indexing finished!")
			return
		}
		if event == blockchain.BlockSaved {
			return
		}
	}
}
