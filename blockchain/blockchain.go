package blockchain

import (
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/encoding"
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
		blockHeader := entity.NewGenesisBlockHeader(ownerAddress, encoding.NewTransactionGobEncoder())
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
	rewardTransaction := entity.NewCoinbaseTx(miner.GetAddress(), encoding.NewTransactionGobEncoder())

	return bc.mineBlock([]*entity.Transaction{expenseTransaction, rewardTransaction})
}

func (bc *Blockchain) mineBlock(transactions []*entity.Transaction) (error) {
	currentHead, err := bc.reader.Head()
	if err != nil {
		return err
	}

	newBlockHeader := entity.NewBlockHeader(currentHead.Index+1, currentHead.Hash, transactions)
	newBlock := bc.miner.MineBlock(newBlockHeader)
	if err = bc.reader.SaveBlock(newBlock); err != nil {
		return err
	}

	return nil
}
