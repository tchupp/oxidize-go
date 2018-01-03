package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/wallet"
)

type Blockchain struct {
	repository Repository
}

func Open(repository Repository, ownerAddress string) (*Blockchain, error) {
	head, err := repository.Head()
	if err == HeadBlockNotFoundError {
		head = NewGenesisBlock(ownerAddress)
		err = repository.SaveBlock(head)
	}
	if err != nil {
		return nil, err
	}

	return &Blockchain{repository: repository}, nil
}

func (bc *Blockchain) Send(sender, receiver, miner *wallet.Wallet, expense uint) (error) {
	expenseTransaction, err := bc.buildExpenseTransaction(sender, receiver, expense)
	if err != nil {
		return err
	}
	rewardTransaction := tx.NewCoinbaseTx(miner.GetAddress())

	return bc.mineBlock([]*tx.Transaction{expenseTransaction, rewardTransaction})
}

func (bc *Blockchain) mineBlock(transactions []*tx.Transaction) (error) {
	currentHead, err := bc.repository.Head()
	if err != nil {
		return err
	}

	newHead := NewBlock(transactions, currentHead.Hash, currentHead.Index+1)
	if err = bc.repository.SaveBlock(newHead); err != nil {
		return err
	}

	return nil
}
