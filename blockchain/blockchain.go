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
	GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)

	SaveHeaders(headers entity.BlockHeaders) error
	SaveBlock(block *entity.Block) error

	Send(spender, receiver, coinbase *identity.Identity, expense uint32) error
}

type blockchain struct {
	blockRepository  entity.BlockRepository
	headerRepository entity.HeaderRepository
	miner            mining.Miner
	utxoEngine       utxo.Engine
}

func Open(blockRepository entity.BlockRepository, headerRepository entity.HeaderRepository, miner mining.Miner) (Blockchain, error) {
	bc := &blockchain{
		blockRepository:  blockRepository,
		headerRepository: headerRepository,
		miner:            miner,
		utxoEngine:       utxo.NewCrawlerEngine(blockRepository),
	}

	exists, err := genesisBlockExists(blockRepository)
	if err != nil {
		return nil, err
	}

	if exists {
		return bc, nil
	}

	genesisBlock := entity.DefaultGenesisBlock()
	if err = blockRepository.SaveBlock(genesisBlock); err != nil {
		return nil, err
	}
	if err = headerRepository.SaveHeader(genesisBlock.Header()); err != nil {
		return nil, err
	}

	return bc, nil
}

func genesisBlockExists(repository entity.BlockRepository) (bool, error) {
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
	return bc.headerRepository.Head()
}

func (bc *blockchain) GetHeader(hash *entity.Hash) (*entity.BlockHeader, error) {
	head, err := bc.blockRepository.Block(hash)
	if err != nil {
		return nil, err
	}
	return head.Header(), nil
}

func (bc *blockchain) GetHeaders(hash *entity.Hash, index uint64) (entity.BlockHeaders, error) {
	// TODO finish unit tests
	startingHeader, err := bc.GetBestHeader()
	if err != nil {
		return nil, err
	}

	headers := entity.BlockHeaders{startingHeader}
	nextHeader := startingHeader
	for {
		if nextHeader.PreviousHash.IsEqual(&entity.EmptyHash) {
			return headers, nil
		}

		nextHeader, err = bc.GetHeader(nextHeader.PreviousHash)
		if err != nil {
			return nil, err
		}

		headers = headers.Add(nextHeader)
	}

	panic("unexpected")
}

func (bc *blockchain) SaveHeaders(headers entity.BlockHeaders) error {
	// TODO verify headers
	for _, header := range headers {
		err := bc.headerRepository.SaveHeader(header)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	err := bc.blockRepository.SaveBlock(block)
	if err != nil {
		return err
	}

	err = bc.headerRepository.SaveHeader(block.Header())
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

	newBlock, err := engine.MineBlock(entity.Transactions{expenseTransaction}, bc.miner, bc.blockRepository)
	if err != nil {
		return err
	}
	return bc.blockRepository.SaveBlock(newBlock)
}
