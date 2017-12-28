package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/wallet"
)

type Blockchain struct {
	head       *Block
	repository Repository
}

func Open(repository Repository, ownerAddress string) (*Blockchain, error) {
	head, err := repository.Head()
	if err == LatestHashNotFoundError {
		head = NewGenesisBlock(ownerAddress)
		err = repository.SaveBlock(head)
	}
	if err != nil {
		return nil, err
	}

	return &Blockchain{head: head, repository: repository}, nil
}

func (bc *Blockchain) Send(sender, receiver *wallet.Wallet, expense uint) (*Blockchain, error) {
	expenseTransaction, err := bc.buildExpenseTransaction(sender, receiver, expense)
	if err != nil {
		return bc, err
	}
	rewardTransaction := tx.NewCoinbaseTx(sender.GetAddress())

	return bc.mineBlock([]*tx.Transaction{expenseTransaction, rewardTransaction})
}

func (bc *Blockchain) mineBlock(transactions []*tx.Transaction) (*Blockchain, error) {
	currentHead, err := bc.repository.Head()
	if err != nil {
		return nil, err
	}

	newHead := NewBlock(transactions, currentHead.Hash, currentHead.Index+1)
	if err = bc.repository.SaveBlock(newHead); err != nil {
		return nil, err
	}

	return &Blockchain{head: newHead, repository: bc.repository}, nil
}
