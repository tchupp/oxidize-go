package boltdb

import (
	"log"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/storage"
)

type builder struct {
	repository entity.ChainRepository
}

func Builder(name string, blockEncoder entity.BlockEncoder) storage.Builder {
	repository, err := NewChainRepository(name, blockEncoder)
	if err != nil {
		log.Panic(err)
	}

	return &builder{repository: repository}
}

func (b builder) WithCache() storage.Builder {
	return &builder{repository: storage.WrapWithCache(b.repository)}
}

func (b builder) WithLogger() storage.Builder {
	return &builder{repository: storage.WrapWithLogger(b.repository)}
}

func (b builder) WithMetrics() storage.Builder {
	return &builder{repository: storage.WrapWithMetrics(b.repository)}
}

func (b builder) Build() entity.ChainRepository {
	return b.repository
}
