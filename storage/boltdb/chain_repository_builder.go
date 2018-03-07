package boltdb

import (
	"log"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/storage"
)

type chainRepoBuilder struct {
	path         string
	name         string
	blockEncoder entity.BlockEncoder
	withCache    bool
	withLogger   bool
	withMetrics  bool
}

func ChainBuilder(name string) storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{name: name, blockEncoder: encoding.BlockProtoEncoder(), path: "./"}
}

func (b *chainRepoBuilder) WithBlockEncoder(blockEncoder entity.BlockEncoder) storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{
		path:         b.path,
		name:         b.name,
		blockEncoder: blockEncoder,
		withCache:    true,
		withLogger:   b.withLogger,
		withMetrics:  b.withMetrics,
	}
}

func (b *chainRepoBuilder) WithPath(path string) storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{
		path:         path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    true,
		withLogger:   b.withLogger,
		withMetrics:  b.withMetrics,
	}
}

func (b *chainRepoBuilder) WithCache() storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    true,
		withLogger:   b.withLogger,
		withMetrics:  b.withMetrics,
	}
}

func (b *chainRepoBuilder) WithLogger() storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    b.withCache,
		withLogger:   true,
		withMetrics:  b.withMetrics,
	}
}

func (b *chainRepoBuilder) WithMetrics() storage.ChainRepositoryBuilder {
	return &chainRepoBuilder{
		path:         b.path,
		name:         b.name,
		blockEncoder: b.blockEncoder,
		withCache:    b.withCache,
		withLogger:   b.withLogger,
		withMetrics:  true,
	}
}

func (b *chainRepoBuilder) Build() entity.ChainRepository {
	repository, err := NewChainRepository(b.path, b.name, b.blockEncoder)
	if err != nil {
		log.Panic(err)
	}

	if b.withCache {
		repository = storage.WrapWithCache(repository)
	}
	if b.withLogger {
		repository = storage.WrapChainWithLogger(repository)
	}
	if b.withMetrics {
		repository = storage.WrapWithMetrics(repository)
	}

	return repository
}
