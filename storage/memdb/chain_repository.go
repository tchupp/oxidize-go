package memdb

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"sync"
)

type chainMemoryRepository struct {
	blockDB      map[string]*entity.Block
	headerDB     map[string]*entity.BlockHeader
	blockHashDB  map[uint64]*entity.Hash
	headerHashDB map[uint64]*entity.Hash
	lock         sync.RWMutex
}

func NewChainRepository() entity.ChainRepository {
	return &chainMemoryRepository{
		blockDB:      make(map[string]*entity.Block),
		headerDB:     make(map[string]*entity.BlockHeader),
		blockHashDB:  make(map[uint64]*entity.Hash),
		headerHashDB: make(map[uint64]*entity.Hash),
	}
}

func NewBlockRepository() entity.BlockRepository {
	return NewChainRepository()
}

func NewHeaderRepository() entity.HeaderRepository {
	return NewChainRepository()
}

func (r *chainMemoryRepository) BestBlock() (head *entity.Block, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	bestIndex := bestIndex(r.blockHashDB)
	if hash, ok := r.blockHashDB[bestIndex]; ok {
		return r.blockDB[hash.String()], nil
	}
	return nil, nil
}

func bestIndex(hashDB map[uint64]*entity.Hash) uint64 {
	bestIndex := uint64(0)
	for index := range hashDB {
		if bestIndex < index {
			bestIndex = index
		}
	}
	return bestIndex
}

func (r *chainMemoryRepository) BlockByHash(hash *entity.Hash) (*entity.Block, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if block, ok := r.blockDB[hash.String()]; ok {
		return block, nil
	}
	return nil, nil
}

func (r *chainMemoryRepository) BlockByIndex(index uint64) (*entity.Block, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if hash, ok := r.blockHashDB[index]; ok {
		if block, ok := r.blockDB[hash.String()]; ok {
			return block, nil
		}
	}
	return nil, nil
}

func (r *chainMemoryRepository) SaveBlock(block *entity.Block) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.blockDB[block.Hash().String()] = block
	r.headerDB[block.Hash().String()] = block.Header()
	r.blockHashDB[block.Index()] = block.Hash()
	r.headerHashDB[block.Index()] = block.Hash()

	return nil
}

func (r *chainMemoryRepository) BestHeader() (head *entity.BlockHeader, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	bestIndex := bestIndex(r.headerHashDB)
	if hash, ok := r.headerHashDB[bestIndex]; ok {
		return r.headerDB[hash.String()], nil
	}
	return nil, nil
}

func (r *chainMemoryRepository) HeaderByHash(hash *entity.Hash) (*entity.BlockHeader, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if header, ok := r.headerDB[hash.String()]; ok {
		return header, nil
	}
	return nil, nil
}

func (r *chainMemoryRepository) HeaderByIndex(index uint64) (*entity.BlockHeader, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if hash, ok := r.headerHashDB[index]; ok {
		if header, ok := r.headerDB[hash.String()]; ok {
			return header, nil
		}
	}
	return nil, nil
}

func (r *chainMemoryRepository) SaveHeader(header *entity.BlockHeader) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.headerDB[header.Hash.String()] = header
	r.headerHashDB[header.Index] = header.Hash

	return nil
}

func (r *chainMemoryRepository) Close() error {
	return nil
}
