package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/wallet"
)

type Blockchain struct {
	repository    Repository
	miningService MiningService
}

func Open(repository Repository, miningService MiningService, ownerAddress string) (*Blockchain, error) {
	exists, err := genesisBlockExists(repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		blockHeader := NewGenesisBlockHeader(ownerAddress)
		block := miningService.MineBlock(blockHeader)
		err = repository.SaveBlock(block)
	}
	if err != nil {
		return nil, err
	}

	return &Blockchain{repository: repository, miningService: miningService}, nil
}

func genesisBlockExists(repository Repository) (bool, error) {
	_, err := repository.Head()
	if err == HeadBlockNotFoundError {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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

	newBlockHeader := NewBlockHeader(currentHead.Index+1, currentHead.Hash, transactions)
	newBlock := bc.miningService.MineBlock(newBlockHeader)
	if err = bc.repository.SaveBlock(newBlock); err != nil {
		return err
	}

	return nil
}
