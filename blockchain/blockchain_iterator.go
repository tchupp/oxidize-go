package blockchain

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

func (bc *Blockchain) head() (*Iterator, error) {
	head, err := bc.blockRepository.Head()
	if err != nil {
		return nil, err
	}

	return &Iterator{
		current:         head,
		blockRepository: bc.blockRepository,
	}, nil
}

func (bc *Blockchain) ForEachBlock(consume func(*entity.Block)) (err error) {
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
