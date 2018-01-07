package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/block"
	"github.com/tclchiam/block_n_go/storage"
)

type Iterator struct {
	current *block.Block
	reader  storage.BlockReader
}

func (it *Iterator) next() (*Iterator, error) {
	b, err := it.reader.Block(it.current.PreviousHash)
	if err != nil {
		return nil, err
	}

	return &Iterator{current: b, reader: it.reader}, nil
}

func (it *Iterator) hasNext() bool {
	return !it.current.IsGenesisBlock()
}

func (bc *Blockchain) head() (*Iterator, error) {
	head, err := bc.reader.Head()
	if err != nil {
		return nil, err
	}

	return &Iterator{
		current: head,
		reader:  bc.reader,
	}, nil
}

func (bc *Blockchain) ForEachBlock(consume func(*block.Block)) (err error) {
	b, err := bc.head()
	if err != nil {
		return err
	}

	for {
		consume(b.current)

		if !b.hasNext() {
			break
		}

		b, err = b.next()
		if err != nil {
			return err
		}
	}
	return nil
}
