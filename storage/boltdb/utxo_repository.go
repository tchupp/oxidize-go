package boltdb

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/wire"
)

const (
	utxoDbFile = "utxo_%s_v1.db"

	txsBucketName        = "txs"
	blockIndexBucketName = "blockIndex"
)

type utxoBoltRepository struct {
	db *bolt.DB
}

func NewUtxoRepository(path, name string) (utxo.Repository, error) {
	db, err := openDB(path, fmt.Sprintf(utxoDbFile, name))
	if err != nil {
		return nil, err
	}

	if err = createBucket(db, txsBucketName); err != nil {
		return nil, err
	}
	if err = createBucket(db, blockIndexBucketName); err != nil {
		return nil, err
	}

	return &utxoBoltRepository{
		db: db,
	}, nil
}

func DeleteUtxo(name string) error {
	return DeleteBoltFile("./", fmt.Sprintf(utxoDbFile, name))
}

func (r *utxoBoltRepository) SaveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, txsBucketName)
		if err != nil {
			return err
		}

		transaction := &entity.Transaction{
			ID:      txId,
			Inputs:  []*entity.SignedInput{},
			Outputs: []*entity.Output{},
			Secret:  []byte{},
		}
		if transactionBytes := bucket.Get(txId.Slice()); len(transactionBytes) != 0 {
			transaction, err = wire.DecodeTransaction(transactionBytes)
			if err != nil {
				return err
			}
		}

		transaction.Outputs = append(transaction.Outputs, output)
		transactionBytes, err := wire.EncodeTransaction(transaction)
		if err != nil {
			return err
		}

		return bucket.Put(txId.Slice(), transactionBytes)
	})
}

func (r *utxoBoltRepository) SaveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, txsBucketName)
		if err != nil {
			return err
		}

		transaction := &entity.Transaction{
			ID:      txId,
			Inputs:  []*entity.SignedInput{},
			Outputs: []*entity.Output{},
			Secret:  []byte{},
		}
		if transactionBytes := bucket.Get(txId.Slice()); len(transactionBytes) != 0 {
			transaction, err = wire.DecodeTransaction(transactionBytes)
			if err != nil {
				return err
			}
		}

		for _, output := range outputs {
			if !transaction.Outputs.Contains(output) {
				transaction.Outputs = append(transaction.Outputs, output)
			}
		}
		transactionBytes, err := wire.EncodeTransaction(transaction)
		if err != nil {
			return err
		}

		return bucket.Put(txId.Slice(), transactionBytes)
	})
}

func (r *utxoBoltRepository) RemoveSpendableOutput(txId *entity.Hash, output *entity.Output) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, txsBucketName)
		if err != nil {
			return err
		}

		transactionBytes := bucket.Get(txId.Slice())
		if len(transactionBytes) == 0 {
			return nil
		}

		transaction, err := wire.DecodeTransaction(transactionBytes)
		if err != nil {
			return err
		}

		transaction.Outputs = transaction.Outputs.Remove(output)
		transactionBytes, err = wire.EncodeTransaction(transaction)
		if err != nil {
			return err
		}

		return bucket.Put(txId.Slice(), transactionBytes)
	})
}

func (r *utxoBoltRepository) RemoveSpendableOutputs(txId *entity.Hash, outputs []*entity.Output) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, txsBucketName)
		if err != nil {
			return err
		}

		transactionBytes := bucket.Get(txId.Slice())
		if len(transactionBytes) == 0 {
			return nil
		}

		transaction, err := wire.DecodeTransaction(transactionBytes)
		if err != nil {
			return err
		}

		for _, output := range outputs {
			transaction.Outputs = transaction.Outputs.Remove(output)
		}
		transactionBytes, err = wire.EncodeTransaction(transaction)
		if err != nil {
			return err
		}

		return bucket.Put(txId.Slice(), transactionBytes)
	})
}

func (r *utxoBoltRepository) SaveSpentOutput(*entity.Hash, *entity.Output) error {
	return nil
}

func (r *utxoBoltRepository) SpendableOutputs() (outputs *utxo.OutputSet, err error) {
	outputs = utxo.NewOutputSet()
	err = r.db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, txsBucketName)
		if err != nil {
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			transaction, err := wire.DecodeTransaction(v)
			if err != nil {
				return err
			}

			outputs = outputs.AddMany(transaction.ID, transaction.Outputs)
			return nil
		})
	})
	return outputs, err
}

func (r *utxoBoltRepository) BlockIndex() (blockIndex *utxo.BlockIndex, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blockIndexBucketName)
		if err != nil {
			return err
		}

		index, err := bestIndex(bucket)
		if err != nil {
			return err
		}
		if index == nil {
			blockIndex = utxo.NewBlockIndex(nil, 0)
			return nil
		}

		hash, err := entity.NewHash(bucket.Get(index))
		if err != nil {
			return err
		}

		blockIndex = utxo.NewBlockIndex(hash, binary.BigEndian.Uint64(index))
		return nil
	})

	return blockIndex, err
}

func (r *utxoBoltRepository) SaveBlockIndex(index *utxo.BlockIndex) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blockIndexBucketName)
		if err != nil {
			return err
		}

		return bucket.Put(toByte(index.Index()), index.Hash().Slice())
	})
}

func (r *utxoBoltRepository) Close() error {
	return nil
}
