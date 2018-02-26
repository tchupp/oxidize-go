package utxo

import (
	"sync"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type Repository interface {
	SaveSpendableOutput(*entity.Hash, *entity.Output) error
	SaveSpendableOutputs(*entity.Hash, []*entity.Output) error

	RemoveSpendableOutput(*entity.Hash, *entity.Output) error
	RemoveSpendableOutputs(*entity.Hash, []*entity.Output) error

	SaveSpentOutput(*entity.Hash, *entity.Output) error

	SpendableOutputs() (*OutputSet, error)
}

type utxoMemRepository struct {
	lock             sync.RWMutex
	spendableOutputs *OutputSet
	spentOutputs     *OutputSet
}

func NewUtxoRepository() Repository {
	return &utxoMemRepository{
		spendableOutputs: NewOutputSet(),
		spentOutputs:     NewOutputSet(),
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

	panic("implement me")
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

func (r *utxoMemRepository) SpendableOutputs() (*OutputSet, error) {
	return r.spendableOutputs.Copy(), nil
}
