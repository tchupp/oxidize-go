package storage

import "github.com/tclchiam/block_n_go/blockchain/entity"

type chainMetrics struct {
	entity.ChainRepository
}

func WrapWithMetrics(repository entity.ChainRepository) entity.ChainRepository {
	// TODO metrics stuff
	return &chainMetrics{
		ChainRepository: repository,
	}
}
