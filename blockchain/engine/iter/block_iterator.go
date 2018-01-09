package iter

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/storage"
)

type Iterator struct {
	current         *entity.Block
	blockRepository storage.BlockRepository
}

func (it *Iterator) next() (*Iterator, error) {
	b, err := it.blockRepository.Block(it.current.PreviousHash)
	if err != nil {
		return nil, err
	}

	return &Iterator{current: b, blockRepository: it.blockRepository}, nil
}

func (it *Iterator) hasNext() bool {
	return !it.current.IsGenesisBlock()
}

func head(blockRepository storage.BlockRepository) (*Iterator, error) {
	head, err := blockRepository.Head()
	if err != nil {
		return nil, err
	}

	return &Iterator{
		current:         head,
		blockRepository: blockRepository,
	}, nil
}

func ForEachBlock(blockRepository storage.BlockRepository, consume func(*entity.Block)) error {
	it, err := head(blockRepository)
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
