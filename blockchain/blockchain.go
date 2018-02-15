package blockchain

import (
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain/engine"
	"github.com/tclchiam/oxidize-go/blockchain/engine/iter"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/engine/utxo"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type Blockchain interface {
	entity.ChainRepository
	ForEachBlock(consume func(*entity.Block)) error

	Balance(address *identity.Address) (*account.Account, error)

	Headers(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)

	SaveHeaders(headers entity.BlockHeaders) error

	Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error
}

type blockchain struct {
	entity.ChainRepository
	miner      mining.Miner
	utxoEngine utxo.Engine
}

func Open(repository entity.ChainRepository, miner mining.Miner) (Blockchain, error) {
	err := engine.ResetGenesis(repository)
	if err != nil {
		return nil, err
	}

	return &blockchain{
		ChainRepository: repository,
		miner:           miner,
		utxoEngine:      utxo.NewCrawlerEngine(repository),
	}, nil
}

func (bc *blockchain) ForEachBlock(consume func(*entity.Block)) error {
	return iter.ForEachBlock(bc.ChainRepository, consume)
}

func (bc *blockchain) Balance(address *identity.Address) (*account.Account, error) {
	return engine.Balance(address, bc.utxoEngine)
}

func (bc *blockchain) BestHeader() (*entity.BlockHeader, error) {
	return bc.ChainRepository.BestHeader()
}

func (bc *blockchain) HeaderByHash(hash *entity.Hash) (*entity.BlockHeader, error) {
	return bc.ChainRepository.HeaderByHash(hash)
}

func (bc *blockchain) HeaderByIndex(index uint64) (*entity.BlockHeader, error) {
	return bc.ChainRepository.HeaderByIndex(index)
}

func (bc *blockchain) Headers(hash *entity.Hash, index uint64) (entity.BlockHeaders, error) {
	startingHeader, err := bc.HeaderByHash(hash)
	if err != nil {
		return nil, err
	}
	if startingHeader == nil {
		return entity.NewBlockHeaders(), nil
	}

	headers := entity.BlockHeaders{startingHeader}
	nextHeader := startingHeader
	for {
		nextHeader, err = bc.HeaderByIndex(nextHeader.Index + 1)
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
	return bc.ChainRepository.SaveHeader(header)
}

func (bc *blockchain) BestBlock() (*entity.Block, error) {
	return bc.ChainRepository.BestBlock()
}

func (bc *blockchain) BlockByHash(hash *entity.Hash) (*entity.Block, error) {
	return bc.ChainRepository.BlockByHash(hash)
}

func (bc *blockchain) BlockByIndex(index uint64) (*entity.Block, error) {
	return bc.ChainRepository.BlockByIndex(index)
}

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	return bc.ChainRepository.SaveBlock(block)
}

func (bc *blockchain) Send(spender *identity.Identity, receiver *identity.Address, expense uint64) error {
	expenseTransaction, err := engine.BuildExpenseTransaction(spender, receiver, expense, bc.utxoEngine)
	if err != nil {
		return err
	}

	newBlock, err := engine.MineBlock(entity.Transactions{expenseTransaction}, bc.miner, bc.ChainRepository)
	if err != nil {
		return err
	}
	return bc.ChainRepository.SaveBlock(newBlock)
}

func (bc *blockchain) Close() error {
	return bc.ChainRepository.Close()
}
