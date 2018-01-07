package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/blockchain/block"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/storage"
)

type Blockchain struct {
	reader storage.BlockReader
	miner  mining.Miner
}

func Open(reader storage.BlockReader, miner mining.Miner, ownerAddress string) (*Blockchain, error) {
	exists, err := genesisBlockExists(reader)
	if err != nil {
		return nil, err
	}

	if !exists {
		blockHeader := block.NewGenesisBlockHeader(ownerAddress)
		b := miner.MineBlock(blockHeader)
		err = reader.SaveBlock(b)
	}
	if err != nil {
		return nil, err
	}

	return &Blockchain{reader: reader, miner: miner}, nil
}

func genesisBlockExists(reader storage.BlockReader) (bool, error) {
	_, err := reader.Head()
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
	currentHead, err := bc.reader.Head()
	if err != nil {
		return err
	}

	newBlockHeader := block.NewBlockHeader(currentHead.Index+1, currentHead.Hash, transactions)
	newBlock := bc.miner.MineBlock(newBlockHeader)
	if err = bc.reader.SaveBlock(newBlock); err != nil {
		return err
	}

	return nil
}
