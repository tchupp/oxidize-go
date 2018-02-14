package boltdb

import (
	"log"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/storage"
)

type builder struct {
	path         string
	name         string
	blockEncoder entity.BlockEncoder
	withCache    bool
	withLogger   bool
	withMetrics  bool
}

func Builder(name string, blockEncoder entity.BlockEncoder) storage.Builder {
	return &builder{name: name, blockEncoder: blockEncoder, path: "./"}
}

func (b *builder) WithPath(path string) storage.Builder {
	return &builder{
		path:         path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    true,
		withLogger:   b.withLogger,
		withMetrics:  b.withMetrics,
	}
}

func (b *builder) WithCache() storage.Builder {
	return &builder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    true,
		withLogger:   b.withLogger,
		withMetrics:  b.withMetrics,
	}
}

func (b *builder) WithLogger() storage.Builder {
	return &builder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    b.withCache,
		withLogger:   true,
		withMetrics:  b.withMetrics,
	}
}

func (b *builder) WithMetrics() storage.Builder {
	return &builder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    b.withCache,
		withLogger:   b.withLogger,
		withMetrics:  true,
	}
}

func (b *builder) Build() entity.ChainRepository {
	repository, err := NewChainRepository(b.path, b.name, b.blockEncoder)
	if err != nil {
		log.Panic(err)
	}

	if b.withCache {
		repository = storage.WrapWithCache(repository)
	}
	if b.withLogger {
		repository = storage.WrapWithLogger(repository)
	}
	if b.withMetrics {
		repository = storage.WrapWithMetrics(repository)
	}

	return repository
}
