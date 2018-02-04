package storage

import "github.com/tclchiam/block_n_go/blockchain/entity"

type Builder interface {
	WithCache() Builder
	WithLogger() Builder
	WithMetrics() Builder
	Build() entity.ChainRepository
}
