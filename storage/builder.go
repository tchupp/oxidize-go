package storage

import "github.com/tclchiam/oxidize-go/blockchain/entity"

type Builder interface {
	WithPath(string) Builder
	WithCache() Builder
	WithLogger() Builder
	WithMetrics() Builder
	Build() entity.ChainRepository
}
