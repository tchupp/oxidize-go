package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/engine"
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/identity"
)

type Blockchain interface {
	ForEachBlock(consume func(*entity.Block)) error

	ReadBalance(identity *identity.Identity) (uint32, error)

	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
	GetHeaderByIndex(index uint64) (*entity.BlockHeader, error)
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)

	SaveHeaders(headers entity.BlockHeaders) error
	SaveHeader(header *entity.BlockHeader) error
	SaveBlock(block *entity.Block) error

	Send(spender, receiver, coinbase *identity.Identity, expense uint32) error
}

type blockchain struct {
	repository  entity.ChainRepository
	miner            mining.Miner
	utxoEngine       utxo.Engine
}

func Open(repository entity.ChainRepository, miner mining.Miner) (Blockchain, error) {
	bc := &blockchain{
		repository: repository,
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

	genesisBlock := entity.DefaultGenesisBlock()
	if err = repository.SaveBlock(genesisBlock); err != nil {
		return nil, err
	}

	return bc, nil
}

func genesisBlockExists(repository entity.ChainRepository) (bool, error) {
	head, err := repository.BestBlock()
	if err != nil {
		return false, err
	}
	if head == nil {
		return false, nil
	}
	return true, nil
}

func (bc *blockchain) ForEachBlock(consume func(*entity.Block)) error {
	return iter.ForEachBlock(bc.repository, consume)
}

func (bc *blockchain) ReadBalance(identity *identity.Identity) (uint32, error) {
	return engine.ReadBalance(identity, bc.utxoEngine)
}

func (bc *blockchain) GetBestHeader() (*entity.BlockHeader, error) {
	return bc.repository.BestHeader()
}

func (bc *blockchain) GetHeader(hash *entity.Hash) (*entity.BlockHeader, error) {
	bestHeader, err := bc.repository.HeaderByHash(hash)
	if err != nil {
		return nil, err
	}
	return bestHeader, nil
}

func (bc *blockchain) GetHeaderByIndex(index uint64) (*entity.BlockHeader, error) {
	bestHeader, err := bc.repository.HeaderByIndex(index)
	if err != nil {
		return nil, err
	}
	return bestHeader, nil
}

func (bc *blockchain) GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error) {
	startingHeader, err := bc.GetHeader(hash)
	if err != nil {
		return nil, err
	}
	if startingHeader == nil {
		return entity.NewBlockHeaders(), nil
	}

	headers := entity.BlockHeaders{startingHeader}
	nextHeader := startingHeader
	for {
		nextHeader, err = bc.GetHeaderByIndex(nextHeader.Index + 1)
		if err != nil {
			return nil, err
		}
		if nextHeader == nil {
			return headers, nil
		}

		headers = headers.Add(nextHeader)
	}

	panic("unexpected")
}

func (bc *blockchain) SaveHeaders(headers entity.BlockHeaders) error {
	return engine.SaveHeaders(headers, bc)
}

func (bc *blockchain) SaveHeader(header *entity.BlockHeader) error {
	// TODO verify header
	return bc.repository.SaveHeader(header)
}

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	err := bc.SaveHeader(block.Header())
	if err != nil {
		return err
	}

	err = bc.repository.SaveBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (bc *blockchain) Send(spender, receiver, coinbase *identity.Identity, expense uint32) error {
	expenseTransaction, err := engine.BuildExpenseTransaction(spender, receiver, expense, bc.utxoEngine)
	if err != nil {
		return err
	}

	newBlock, err := engine.MineBlock(entity.Transactions{expenseTransaction}, bc.miner, bc.repository)
	if err != nil {
		return err
	}
	return bc.repository.SaveBlock(newBlock)
}
