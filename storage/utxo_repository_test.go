package storage_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/storage/boltdb"
	"github.com/tclchiam/oxidize-go/storage/memdb"
)

func TestRepository_BlockIndexCanBeReadBack(t *testing.T) {
	suite := func(repository utxo.Repository, t *testing.T) {
		tx := entity.NewRewardTx(identity.RandomIdentity().Address())
		expectedIndex := utxo.NewBlockIndex(tx.ID, 7)
		err := repository.SaveBlockIndex(expectedIndex)
		if err != nil {
			t.Fatalf("failed to save block index: %s", err)
		}

		blockIndex, err := repository.BlockIndex()
		if err != nil {
			t.Fatalf("failed to read block index: %s", err)
		}

		if !blockIndex.IsEqual(expectedIndex) {
			t.Errorf("unexpected block index. Got - %s, wanted - %s", blockIndex, expectedIndex)
		}

		assert.NoError(t, repository.Close())
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewUtxoRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.UtxoBuilder("BlockIndexCanBeReadBack").
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer boltdb.DeleteUtxo("BlockIndexCanBeReadBack")

		suite(repository, t)
	})
}

func TestRepository_SpendableOutputsCanBeEmpty(t *testing.T) {
	suite := func(repository utxo.Repository, t *testing.T) {
		outputs, err := repository.SpendableOutputs()
		if err != nil {
			t.Fatalf("failed to read spendable outputs: %s", err)
		}
		if outputs == nil {
			t.Fatalf("spendable outputs is nil")
		}

		expectedOutputs := utxo.NewOutputSet()
		if !reflect.DeepEqual(expectedOutputs, outputs) {
			t.Errorf("unexpected outputs. Got - %s, wanted - %s", outputs, expectedOutputs)
		}

		assert.NoError(t, repository.Close())
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewUtxoRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.UtxoBuilder("SpendableOutputsCanBeEmpty").
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer boltdb.DeleteUtxo("SpendableOutputsCanBeEmpty")

		suite(repository, t)
	})
}

func TestRepository_CanRetrieveSpendableOutputs(t *testing.T) {
	suite := func(repository utxo.Repository, t *testing.T) {
		address := identity.RandomIdentity().Address()
		tx := entity.NewRewardTx(address)

		err := repository.SaveSpendableOutputs(tx.ID, tx.Outputs)
		if err != nil {
			t.Fatalf("failed to save spendable outputs: %s", err)
		}

		outputs, err := repository.SpendableOutputs()
		if err != nil {
			t.Fatalf("failed to read spendable outputs: %s", err)
		}
		if outputs == nil {
			t.Fatalf("spendable outputs is nil")
		}

		expectedOutputs := utxo.NewOutputSet().AddMany(tx.ID, tx.Outputs)
		if !reflect.DeepEqual(expectedOutputs, outputs) {
			t.Errorf("unexpected outputs. Got - %s, wanted - %s", outputs, expectedOutputs)
		}

		assert.NoError(t, repository.Close())
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewUtxoRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.UtxoBuilder("CanRetrieveSpendableOutputs").
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer boltdb.DeleteUtxo("CanRetrieveSpendableOutputs")

		suite(repository, t)
	})
}

func TestRepository_CanRemoveSpendableOutputs(t *testing.T) {
	suite := func(repository utxo.Repository, t *testing.T) {
		address := identity.RandomIdentity().Address()
		output := entity.NewOutput(13, address)

		tx := entity.NewTx(nil, []*entity.Output{output})

		err := repository.SaveSpendableOutput(tx.ID, output)
		if err != nil {
			t.Fatalf("failed to save spendable outputs: %s", err)
		}
		err = repository.RemoveSpendableOutput(tx.ID, output)
		if err != nil {
			t.Fatalf("failed to save spendable outputs: %s", err)
		}

		outputs, err := repository.SpendableOutputs()
		if err != nil {
			t.Fatalf("failed to read spendable outputs: %s", err)
		}

		expectedEmptyOutputs := utxo.NewOutputSet()
		if !reflect.DeepEqual(expectedEmptyOutputs, outputs) {
			t.Errorf("unexpected outputs. Got - %s, wanted - %s", outputs, expectedEmptyOutputs)
		}

		assert.NoError(t, repository.Close())
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewUtxoRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.UtxoBuilder("CanRemoveSpendableOutputs").
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer boltdb.DeleteUtxo("CanRemoveSpendableOutputs")

		suite(repository, t)
	})
}
