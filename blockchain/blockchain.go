package blockchain

import (
	"fmt"
	"os"
	"github.com/tclchiam/block_n_go/tx"
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

func Open(nodeName string, address string) (*Blockchain, error) {
	db, err := openDB(nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	head, err := open(db, []byte(blockBucketName), address)
	if err != nil {
		return nil, err
	}

	return &Blockchain{head: head, nodeName: nodeName}, nil
}

func (bc *Blockchain) NewBlock(transactions []*tx.Transaction) (*Blockchain, error) {
	nodeName := bc.nodeName

	db, err := openDB(nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	head, err := newBlock(db, []byte(blockBucketName), transactions)
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

func (bc *Blockchain) ForEachBlock(consume func(*Block)) (err error) {
	block := bc.Head()

	for {
		consume(block.current)

		if !block.HasNext() {
			break
		}

		block, err = block.Next()
		if err != nil {
			return err
		}
	}
	return nil
}
