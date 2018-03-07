package storage

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
)

type loggingUtxoStorageDecorator struct {
	utxo.Repository
}

func WrapUtxoWithLogger(repository utxo.Repository) utxo.Repository {
	return &loggingUtxoStorageDecorator{Repository: repository}
}

func (d *loggingUtxoStorageDecorator) SaveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	utxoLogger := log.WithField("transaction", txId)

	if err := d.Repository.SaveSpendableOutput(txId, output); err != nil {
		utxoLogger.WithError(err).Warn("failed to save spendable output")
	}

	utxoLogger.Debug("saved spendable output")
	return nil
}

func (d *loggingUtxoStorageDecorator) SaveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	utxoLogger := log.WithField("transaction", txId)

	if err := d.Repository.SaveSpendableOutputs(txId, outputs); err != nil {
		utxoLogger.WithError(err).Warn("failed to save spendable outputs")
	}

	utxoLogger.Debug("saved spendable outputs")
	return nil
}

func (d *loggingUtxoStorageDecorator) RemoveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	utxoLogger := log.WithField("transaction", txId)

	if err := d.Repository.RemoveSpendableOutput(txId, output); err != nil {
		utxoLogger.WithError(err).Warn("failed to remove spendable output")
	}

	utxoLogger.Debug("removed spendable output")
	return nil
}

func (d *loggingUtxoStorageDecorator) RemoveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	utxoLogger := log.WithField("transaction", txId)

	if err := d.Repository.RemoveSpendableOutputs(txId, outputs); err != nil {
		utxoLogger.WithError(err).Warn("failed to remove spendable outputs")
	}

	utxoLogger.Debug("removed spendable outputs")
	return nil
}

func (d *loggingUtxoStorageDecorator) SaveSpentOutput(txId *entity.Hash, output *entity.Output) error {
	return d.Repository.SaveSpentOutput(txId, output)
}

func (d *loggingUtxoStorageDecorator) SaveBlockIndex(blockIndex *utxo.BlockIndex) error {
	return d.Repository.SaveBlockIndex(blockIndex)
}
