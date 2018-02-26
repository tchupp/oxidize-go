package utxo

import (
	"sync"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type BlockIndex struct {
	hash  *entity.Hash
	index uint64
}

type Engine interface {
	UpdateIndex(block *entity.Block) (*BlockIndex, error)
	FirstSpendableIndex() (*BlockIndex, error)
	FirstSpendableIndexForAddress(*identity.Address) (*BlockIndex, error)

	SpendableOutputs(*identity.Address) (*OutputSet, error)
	IsSpendable(*BlockIndex, *entity.Output) (bool, error)
}

type engine struct {
	lock sync.RWMutex

	blockIndex   *BlockIndex
	chainReader  entity.ChainReader
	txRepository Repository
}

func NewUtxoEngine(txRepository Repository, chainReader entity.ChainReader) Engine {
	return &engine{
		txRepository: txRepository,
		chainReader:  chainReader,
		blockIndex:   &BlockIndex{index: 0},
	}
}

func (e *engine) UpdateIndex(block *entity.Block) (*BlockIndex, error) {
	err := e.processToBlock(block)
	if err != nil {
		return e.blockIndex, err
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
	e.lock.RLock()
	defer e.lock.RUnlock()

	outputs, err := e.txRepository.SpendableOutputs()
	if err != nil {
		return nil, err
	}

	return outputs.FilterByOutput(func(output *entity.Output) bool { return output.ReceivedBy(address) }), nil
}

func (e *engine) IsSpendable(*BlockIndex, *entity.Output) (bool, error) {
	panic("implement me")
}

func (e *engine) processToBlock(block *entity.Block) error {
	for e.blockIndex.index+1 < block.Index() {
		b, err := e.chainReader.BlockByIndex(e.blockIndex.index)
		if err != nil {
			return err
		}

		_, err = e.updateIndex(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *engine) updateIndex(block *entity.Block) (*BlockIndex, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, tx := range block.Transactions() {
		for _, output := range tx.Outputs {
			if err := e.txRepository.SaveSpendableOutput(tx.ID, output); err != nil {
				return e.blockIndex, err
			}
		}
		for _, input := range tx.Inputs {
			if err := e.txRepository.RemoveSpendableOutput(input.OutputReference.ID, input.OutputReference.Output); err != nil {
				return e.blockIndex, err
			}
		}
	}

	e.blockIndex = &BlockIndex{hash: block.Hash(), index: block.Index()}
	return e.blockIndex, nil
}
