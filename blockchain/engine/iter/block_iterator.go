package iter

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type Iterator struct {
	current *entity.Block
	reader  entity.ChainReader
}

func (it *Iterator) next() (*Iterator, error) {
	b, err := it.reader.BlockByHash(it.current.PreviousHash())
	if err != nil {
		return nil, err
	}

	return &Iterator{current: b, reader: it.reader}, nil
}

func (it *Iterator) hasNext() bool {
	return !it.current.IsGenesisBlock()
}

func head(reader entity.ChainReader) (*Iterator, error) {
	head, err := reader.BestBlock()
	if err != nil {
		return nil, err
	}

	return &Iterator{
		current: head,
		reader:  reader,
	}, nil
}

func ForEachBlock(reader entity.ChainReader, consume func(*entity.Block)) error {
	it, err := head(reader)
	if err != nil {
		return err
	}

	for {
		consume(it.current)

		if !it.hasNext() {
			break
		}

		it, err = it.next()
		if err != nil {
			return err
		}
	}
	return nil
}
