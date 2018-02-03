package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/engine"
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/mining"
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

func (bc *blockchain) ReadBalance(identity *identity.Identity) (uint32, error) {
	return engine.ReadBalance(identity, bc.utxoEngine)
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

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	err := bc.repository.SaveBlock(block)
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
