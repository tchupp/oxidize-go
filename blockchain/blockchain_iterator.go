package blockchain

type Iterator struct {
	current    *Block
	repository Repository
}

func (it *Iterator) Next() (*Iterator, error) {
	block, err := it.repository.Block(it.current.PreviousHash)
	if err != nil {
		return nil, err
	}

	return &Iterator{current: block, repository: it.repository}, nil
}

func (it *Iterator) HasNext() bool {
	return !it.current.IsGenesisBlock()
}

func (bc *Blockchain) Head() (*Iterator, error) {
	head, err := bc.repository.Head()
	if err != nil {
		return nil, err
	}

	return &Iterator{
		current:    head,
		repository: bc.repository,
	}, nil
}

func (bc *Blockchain) ForEachBlock(consume func(*Block)) (err error) {
	block, err := bc.Head()
	if err != nil {
		return err
	}

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
