package memdb

import (
	"sync"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
)

type utxoMemRepository struct {
	lock             sync.RWMutex
	spendableOutputs *utxo.OutputSet
	spentOutputs     *utxo.OutputSet
	blockIndex       *utxo.BlockIndex
}

func NewUtxoRepository() utxo.Repository {
	return &utxoMemRepository{
		spendableOutputs: utxo.NewOutputSet(),
		spentOutputs:     utxo.NewOutputSet(),
	}
}

func (r *utxoMemRepository) SaveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.spendableOutputs = r.spendableOutputs.Add(txId, output)

	return nil
}

func (r *utxoMemRepository) SaveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, output := range outputs {
		r.spendableOutputs = r.spendableOutputs.Add(txId, output)
	}

	return nil
}

func (r *utxoMemRepository) RemoveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.spendableOutputs = r.spendableOutputs.Remove(txId, output)

	return nil
}

func (r *utxoMemRepository) RemoveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, output := range outputs {
		r.spendableOutputs = r.spendableOutputs.Remove(txId, output)
	}

	return nil
}

func (r *utxoMemRepository) SaveSpentOutput(txId *entity.Hash, output *entity.Output) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.spentOutputs = r.spentOutputs.Add(txId, output)

	return nil
}

func (r *utxoMemRepository) SpendableOutputs() (*utxo.OutputSet, error) {
	return r.spendableOutputs.Copy(), nil
}

func (r *utxoMemRepository) BlockIndex() (*utxo.BlockIndex, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.blockIndex, nil
}

func (r *utxoMemRepository) SaveBlockIndex(blockIndex *utxo.BlockIndex) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.blockIndex = blockIndex
	return nil
}
