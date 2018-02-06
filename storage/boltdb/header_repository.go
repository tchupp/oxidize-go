package boltdb

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func bestHeaderHash(tx *bolt.Tx) (hash *entity.Hash, err error) {
	bucket, err := bucket(tx, headersBucketName)
	if err != nil {
		return nil, err
	}

	index, err := bestIndex(bucket)
	if err != nil {
		return nil, err
	}

	rawHash := bucket.Get(index)
	return entity.NewHash(rawHash)
}

func headerByHash(tx *bolt.Tx, encoder entity.BlockEncoder, hash *entity.Hash) (*entity.BlockHeader, error) {
	bucket, err := bucket(tx, headersBucketName)
	if err != nil {
		return nil, err
	}

	header, err := readHeader(bucket, hash, encoder)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func readHeader(bucket *bolt.Bucket, hash *entity.Hash, encoder entity.BlockEncoder) (*entity.BlockHeader, error) {
	latestHeaderData := bucket.Get(hash.Slice())
	if latestHeaderData == nil || len(latestHeaderData) == 0 {
		return nil, nil
	}

	header, err := encoder.DecodeHeader(latestHeaderData)
	if err != nil {
		return nil, err
	}

	return header, err
}

func saveHeader(tx *bolt.Tx, encoder entity.BlockEncoder, header *entity.BlockHeader) error {
	bucket, err := bucket(tx, headersBucketName)
	if err != nil {
		return err
	}

	headerData, err := encoder.EncodeHeader(header)
	if err != nil {
		return fmt.Errorf("serializing header: %s", err)
	}

	err = bucket.Put(header.Hash.Slice(), headerData)
	if err != nil {
		return fmt.Errorf("writing header: %s", err)
	}

	err = bucket.Put(toByte(header.Index), header.Hash.Slice())
	if err != nil {
		return fmt.Errorf("writing header hash: %s", err)
	}

	return nil
}

func toByte(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
