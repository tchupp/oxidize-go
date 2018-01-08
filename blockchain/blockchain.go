package blockchain

import (
	"fmt"

	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/blockchain/engine"
)

type Blockchain struct {
	blockRepository storage.BlockRepository
	miner  mining.Miner
}

func Open(repository storage.BlockRepository, miner mining.Miner, ownerAddress string) (*Blockchain, error) {
	exists, err := genesisBlockExists(repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		blockHeader := entity.NewGenesisBlockHeader(ownerAddress, encoding.NewTransactionGobEncoder())
		b := miner.MineBlock(blockHeader)
		err = repository.SaveBlock(b)
	}
	if err != nil {
		return nil, err
	}

	return &Blockchain{blockRepository: repository, miner: miner}, nil
}

func genesisBlockExists(repository storage.BlockRepository) (bool, error) {
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
	rewardTransaction := entity.NewCoinbaseTx(miner.GetAddress(), encoding.NewTransactionGobEncoder())

	return bc.mineBlock([]*entity.Transaction{expenseTransaction, rewardTransaction})
}

func (bc *Blockchain) mineBlock(transactions []*entity.Transaction) (error) {
	currentHead, err := bc.blockRepository.Head()
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		for index, input := range transaction.Inputs {
			verified := engine.VerifySignature(input, transaction.Outputs, encoding.NewTransactionGobEncoder())
			if !verified {
				return fmt.Errorf(TransactionInputHasBadSignatureMessage, transaction.ID, index)
			}
		}
	}

	newBlockHeader := entity.NewBlockHeader(currentHead.Index+1, currentHead.Hash, transactions)
	newBlock := bc.miner.MineBlock(newBlockHeader)
	if err = bc.blockRepository.SaveBlock(newBlock); err != nil {
		return err
	}

	return nil
}
