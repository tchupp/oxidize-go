package storage

import (
	"github.com/hashicorp/golang-lru"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

const (
	blockCacheSize  = 128
	headerCacheSize = 128
)

type chainCache struct {
	entity.ChainRepository

	blockCache  *lru.Cache
	headerCache *lru.Cache
}

func WrapWithCache(repository entity.ChainRepository) entity.ChainRepository {
	blockCache, _ := lru.New(blockCacheSize)
	headerCache, _ := lru.New(headerCacheSize)

	return &chainCache{
		ChainRepository: repository,
		blockCache:      blockCache,
		headerCache:     headerCache,
	}
}

func (c *chainCache) BlockByHash(hash *entity.Hash) (*entity.Block, error) {
	if block, ok := c.blockCache.Get(hash); ok {
		return block.(*entity.Block), nil
	}

	block, err := c.ChainRepository.BlockByHash(hash)
	if err == nil {
		c.blockCache.Add(hash, block)
	}
	return block, err
}

func (c *chainCache) BlockByIndex(index uint64) (*entity.Block, error) {
	if block, ok := c.blockCache.Get(index); ok {
		return block.(*entity.Block), nil
	}

	block, err := c.ChainRepository.BlockByIndex(index)
	if err == nil {
		c.blockCache.Add(index, block)
	}
	return block, err
}

func (c *chainCache) HeaderByHash(hash *entity.Hash) (*entity.BlockHeader, error) {
	if header, ok := c.headerCache.Get(hash); ok {
		return header.(*entity.BlockHeader), nil
	}

	header, err := c.ChainRepository.HeaderByHash(hash)
	if err == nil {
		c.headerCache.Add(hash, header)
	}
	return header, err
}

func (c *chainCache) HeaderByIndex(index uint64) (*entity.BlockHeader, error) {
	if header, ok := c.headerCache.Get(index); ok {
		return header.(*entity.BlockHeader), nil
	}

	header, err := c.ChainRepository.HeaderByIndex(index)
	if err == nil {
		c.headerCache.Add(index, header)
	}
	return header, err
}
