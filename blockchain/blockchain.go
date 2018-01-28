package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/engine"
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/identity"
)

type Blockchain interface {
	ForEachBlock(consume func(*entity.Block)) error
	ReadBalance(identity *identity.Identity) (uint32, error)
	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
	Send(spender, receiver, coinbase *identity.Identity, expense uint32) error
}

type blockchain struct {
	blockRepository storage.BlockRepository
	miner           mining.Miner
	utxoEngine      utxo.Engine
}

func Open(repository storage.BlockRepository, miner mining.Miner) (Blockchain, error) {
	bc := &blockchain{
		blockRepository: repository,
		miner:           miner,
		utxoEngine:      utxo.NewCrawlerEngine(repository),
	}

	exists, err := genesisBlockExists(repository)
	if err != nil {
		return nil, err
	}

	if exists {
		return bc, nil
	}

	err = repository.SaveBlock(entity.DefaultGenesisBlock())
	if err != nil {
		return nil, err
	}

	return bc, nil
}

func genesisBlockExists(repository storage.BlockRepository) (bool, error) {
	head, err := repository.Head()
	if err != nil {
		return false, err
	}
	if head == nil {
		return false, nil
	}
	return true, nil
}

func (bc *blockchain) ForEachBlock(consume func(*entity.Block)) error {
	return iter.ForEachBlock(bc.blockRepository, consume)
}

func (bc *blockchain) ReadBalance(identity *identity.Identity) (uint32, error) {
	return engine.ReadBalance(identity, bc.utxoEngine)
}

func (bc *blockchain) GetBestHeader() (*entity.BlockHeader, error) {
	head, err := bc.blockRepository.Head()
	if err != nil {
		return nil, err
	}
	return head.Header(), nil
}

func (bc *blockchain) GetHeader(hash *entity.Hash) (*entity.BlockHeader, error) {
	head, err := bc.blockRepository.Block(hash)
	if err != nil {
		return nil, err
	}
	return head.Header(), nil
}

func (bc *blockchain) Send(spender, receiver, coinbase *identity.Identity, expense uint32) error {
	expenseTransaction, err := engine.BuildExpenseTransaction(spender, receiver, expense, bc.utxoEngine)
	if err != nil {
		return err
	}

	newBlock, err := engine.MineBlock(entity.Transactions{expenseTransaction}, bc.miner, bc.blockRepository)
	if err != nil {
		return err
	}
	return bc.blockRepository.SaveBlock(newBlock)
}
