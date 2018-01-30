package memdb

import (
	"sync"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type headerMemoryRepository struct {
	head *entity.BlockHeader
	db   map[*entity.Hash]*entity.BlockHeader
	lock sync.RWMutex
}

func NewHeaderRepository() entity.HeaderRepository {
	return &headerMemoryRepository{db: make(map[*entity.Hash]*entity.BlockHeader)}
}

func (r *headerMemoryRepository) BestHeader() (head *entity.BlockHeader, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.head != nil {
		return r.head, nil
	}
	return nil, nil
}

func (r *headerMemoryRepository) Header(hash *entity.Hash) (*entity.BlockHeader, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if header, ok := r.db[hash]; ok {
		return header, nil
	}
	return nil, nil
}

func (r *headerMemoryRepository) SaveHeader(header *entity.BlockHeader) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.db[header.Hash] = header
	if r.head == nil {
		r.head = header
	} else if r.head.Index < header.Index {
		r.head = header
	}

	return nil
}

func (r *headerMemoryRepository) Close() error {
	return nil
}
