package blockchain

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine"
	"github.com/tclchiam/oxidize-go/blockchain/engine/iter"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/engine/utxo"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type Blockchain interface {
	ForEachBlock(consume func(*entity.Block)) error

	Balance(identity *identity.Address) (uint32, error)

	GetBestHeader() (*entity.BlockHeader, error)
	GetHeader(hash *entity.Hash) (*entity.BlockHeader, error)
	GetHeaderByIndex(index uint64) (*entity.BlockHeader, error)
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)

	SaveHeaders(headers entity.BlockHeaders) error
	SaveHeader(header *entity.BlockHeader) error

	GetBestBlock() (*entity.Block, error)
	GetBlock(hash *entity.Hash) (*entity.Block, error)
	GetBlockByIndex(index uint64) (*entity.Block, error)
	SaveBlock(block *entity.Block) error

	Send(spender *identity.Identity, receiver *identity.Address, expense uint32) error
}

type blockchain struct {
	repository entity.ChainRepository
	miner      mining.Miner
	utxoEngine utxo.Engine
}

func Open(repository entity.ChainRepository, miner mining.Miner) (Blockchain, error) {
	err := engine.ResetGenesis(repository)
	if err != nil {
		return nil, err
	}

	return &blockchain{
		repository: repository,
		miner:      miner,
		utxoEngine: utxo.NewCrawlerEngine(repository),
	}, nil
}

func (bc *blockchain) ForEachBlock(consume func(*entity.Block)) error {
	return iter.ForEachBlock(bc.repository, consume)
}

func (bc *blockchain) Balance(address *identity.Address) (uint32, error) {
	return engine.ReadBalance(address, bc.utxoEngine)
}

func (bc *blockchain) GetBestHeader() (*entity.BlockHeader, error) {
	return bc.repository.BestHeader()
}

func (bc *blockchain) GetHeader(hash *entity.Hash) (*entity.BlockHeader, error) {
	return bc.repository.HeaderByHash(hash)
}

func (bc *blockchain) GetHeaderByIndex(index uint64) (*entity.BlockHeader, error) {
	return bc.repository.HeaderByIndex(index)
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

func (bc *blockchain) GetBestBlock() (*entity.Block, error) {
	return bc.repository.BestBlock()
}

func (bc *blockchain) GetBlock(hash *entity.Hash) (*entity.Block, error) {
	return bc.repository.BlockByHash(hash)
}

func (bc *blockchain) GetBlockByIndex(index uint64) (*entity.Block, error) {
	return bc.repository.BlockByIndex(index)
}

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	return bc.repository.SaveBlock(block)
}

func (bc *blockchain) Send(spender *identity.Identity, receiver *identity.Address, expense uint32) error {
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
