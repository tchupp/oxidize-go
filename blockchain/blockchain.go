package blockchain

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/identity"
)

type Blockchain interface {
	entity.ChainRepository

	SpendableOutputs(*identity.Address) (*utxo.OutputSet, error)
	IsSpendable(*utxo.BlockIndex, *entity.Output) (bool, error)

	Headers(hash *entity.Hash, index uint64) (entity.BlockHeaders, error)
	SaveHeaders(headers entity.BlockHeaders) error
	MineBlock(transactions entity.Transactions) (*entity.Block, error)

	Subscribe(channel chan<- Event) Subscription
}

type blockchain struct {
	entity.ChainRepository
	utxo.Engine
	miner mining.Miner
	feed  *Feed
}

func Open(repository entity.ChainRepository, miner mining.Miner) (Blockchain, error) {
	bc := &blockchain{
		ChainRepository: repository,
		Engine:          utxo.NewUtxoEngine(utxo.NewUtxoRepository(), repository),
		miner:           miner,
		feed:            NewFeed(),
	}

	err := engine.ResetGenesis(bc)
	if err != nil {
		return nil, err
	}

	block, err := bc.BestBlock()
	if err != nil {
		return nil, err
	}

	_, err = bc.Engine.UpdateIndex(block)
	if err != nil {
		return nil, err
	}
	return bc, nil
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
	err := bc.ChainRepository.SaveHeader(header)
	if err != nil {
		return err
	}

	bc.feed.Send(HeaderSaved)
	return nil
}

func (bc *blockchain) SaveBlock(block *entity.Block) error {
	// TODO verify block
	err := bc.ChainRepository.SaveBlock(block)
	if err != nil {
		return err
	}

	bc.feed.Send(BlockSaved)
	bc.Engine.UpdateIndex(block)
	return nil
}

func (bc *blockchain) MineBlock(transactions entity.Transactions) (*entity.Block, error) {
	return engine.MineBlock(transactions, bc.miner, bc.ChainRepository)
}

func (bc *blockchain) Subscribe(channel chan<- Event) Subscription {
	return bc.feed.Subscribe(channel)
}
