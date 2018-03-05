package utxo

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type Repository interface {
	SaveSpendableOutput(*entity.Hash, *entity.Output) error
	SaveSpendableOutputs(*entity.Hash, []*entity.Output) error

	RemoveSpendableOutput(*entity.Hash, *entity.Output) error
	RemoveSpendableOutputs(*entity.Hash, []*entity.Output) error

	SaveSpentOutput(*entity.Hash, *entity.Output) error

	SpendableOutputs() (*OutputSet, error)

	BlockIndex() (*BlockIndex, error)
	SaveBlockIndex(*BlockIndex) error

	Close() error
}
