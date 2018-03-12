package storage

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
)

type ChainRepositoryBuilder interface {
	WithPath(string) ChainRepositoryBuilder
	WithCache() ChainRepositoryBuilder
	WithLogger() ChainRepositoryBuilder
	WithMetrics() ChainRepositoryBuilder
	Build() entity.ChainRepository
}

type UtxoRepositoryBuilder interface {
	WithPath(string) UtxoRepositoryBuilder
	WithCache() UtxoRepositoryBuilder
	WithLogger() UtxoRepositoryBuilder
	WithMetrics() UtxoRepositoryBuilder
	Build() utxo.Repository
}
