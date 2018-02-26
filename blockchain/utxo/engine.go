package utxo

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type BlockIndex struct {
	hash  *entity.Hash
	index uint64
}

func (i *BlockIndex) Hash() *entity.Hash {
	if i == nil {
		return nil
	}
	return i.hash
}

func (i *BlockIndex) Index() uint64 {
	if i == nil {
		return 0
	}
	return i.index
}

type Engine interface {
	UpdateIndex(block *entity.Block) (*BlockIndex, error)
	FirstSpendableIndex() (*BlockIndex, error)
	FirstSpendableIndexForAddress(*identity.Address) (*BlockIndex, error)

	SpendableOutputs(*identity.Address) (*OutputSet, error)
	IsSpendable(*entity.Hash, *entity.Output) (bool, error)
}

type engine struct {
	chainReader entity.ChainReader
	repository  Repository
}

func NewUtxoEngine(repository Repository, chainReader entity.ChainReader) Engine {
	return &engine{
		repository:  repository,
		chainReader: chainReader,
	}
}

func (e *engine) UpdateIndex(block *entity.Block) (*BlockIndex, error) {
	err := e.processToBlock(block)
	if err != nil {
		return nil, err
	}

	return e.updateIndex(block)
}

func (e *engine) FirstSpendableIndex() (*BlockIndex, error) {
	panic("implement me")
}

func (e *engine) FirstSpendableIndexForAddress(*identity.Address) (*BlockIndex, error) {
	panic("implement me")
}

func (e *engine) SpendableOutputs(address *identity.Address) (*OutputSet, error) {
	outputs, err := e.repository.SpendableOutputs()
	if err != nil {
		return nil, err
	}

	return outputs.FilterByOutput(func(output *entity.Output) bool { return output.ReceivedBy(address) }), nil
}

func (e *engine) IsSpendable(txId *entity.Hash, output *entity.Output) (bool, error) {
	panic("implement me")
}

func (e *engine) processToBlock(block *entity.Block) error {
	blockIndex, err := e.repository.BlockIndex()
	if err != nil {
		return err
	}

	for blockIndex.Index()+1 < block.Index() {
		b, err := e.chainReader.BlockByIndex(blockIndex.Index())
		if err != nil {
			return err
		}

		blockIndex, err = e.updateIndex(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *engine) updateIndex(block *entity.Block) (*BlockIndex, error) {
	for _, tx := range block.Transactions() {
		if err := e.repository.SaveSpendableOutputs(tx.ID, tx.Outputs); err != nil {
			return nil, err
		}

		for _, input := range tx.Inputs {
			if err := e.repository.RemoveSpendableOutput(input.OutputReference.ID, input.OutputReference.Output); err != nil {
				return nil, err
			}
		}
	}

	blockIndex := &BlockIndex{hash: block.Hash(), index: block.Index()}
	return blockIndex, e.repository.SaveBlockIndex(blockIndex)
}
