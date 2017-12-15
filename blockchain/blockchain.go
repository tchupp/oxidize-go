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
