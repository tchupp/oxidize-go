package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "blockchain_%s.db"
const blockBucketName = "blocks"

type Blockchain struct {
	head     *Block
	nodeName string
}

func (bc *Blockchain) Head() *Iterator {
	return &Iterator{
		current:    bc.head,
		nodeName:   bc.nodeName,
		bucketName: []byte(blockBucketName),
	}
}

func Open(nodeName string) (*Blockchain, error) {
	db, err := openDB(nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	head, err := open(db, []byte(blockBucketName))
	if err != nil {
		return nil, err
	}

	return &Blockchain{head: head, nodeName: nodeName}, nil
}

func (bc *Blockchain) NewBlock(data string) (*Blockchain, error) {
	nodeName := bc.nodeName

	db, err := openDB(nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	head, err := newBlock(db, []byte(blockBucketName), data)
	if err != nil {
		return nil, err
	}

	return &Blockchain{head: head, nodeName: nodeName}, nil
}

func (bc *Blockchain) Delete() error {
	dbFile := fmt.Sprintf(dbFile, bc.nodeName)
	err := os.Remove(dbFile)
	if err != nil {
		return fmt.Errorf("deleting blockchain file: %s", err)
	}
	return nil
}

func openDB(nodeName string) (*bolt.DB, error) {
	dbFile := fmt.Sprintf(dbFile, nodeName)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening db: %s", err)
	}
	return db, err
}
