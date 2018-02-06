package storage

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type Builder interface {
	WithCache() Builder
	WithLogger() Builder
	WithMetrics() Builder
	Build() entity.ChainRepository
}
