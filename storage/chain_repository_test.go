package storage_test

import (
	"math"
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/storage/boltdb"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

var (
	blockEncoder = encoding.BlockProtoEncoder()
	miner        = proofofwork.NewDefaultMiner(identity.RandomIdentity())

	block1 = miner.MineBlock(&entity.BlockHeader{Index: math.MaxUint64, Hash: &entity.EmptyHash}, entity.Transactions{})
	block2 = miner.MineBlock(block1.Header(), entity.Transactions{})
)

func TestRepository_ReturnsNilIfBlockDoesNotExist(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		blockByHash, err := repository.BlockByHash(block1.Hash())
		if err != nil {
			t.Fatalf("failed to read block by hash: %s", err)
		}
		if blockByHash != nil {
			t.Errorf("expected block to be nil. Got - %s, wanted - %s", blockByHash, nil)
		}

		blockByIndex, err := repository.BlockByIndex(block1.Index())
		if err != nil {
			t.Fatalf("failed to read block by index: %s", err)
		}
		if blockByIndex != nil {
			t.Errorf("expected block to be nil. Got - %s, wanted - %s", blockByIndex, nil)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("ReturnsNilIfBlockDoesNotExist", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("ReturnsNilIfBlockDoesNotExist")

		suite(repository, t)
	})
}

func TestRepository_CanReadSavedBlocks(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		err := repository.SaveBlock(block1)
		if err != nil {
			t.Fatalf("failed to save block1: %s", err)
		}

		blockByHash, err := repository.BlockByHash(block1.Hash())
		if err != nil {
			t.Fatalf("failed to read block by hash: %s", err)
		}
		if !blockByHash.IsEqual(block1) {
			t.Errorf("expected read block to equal block1. Got - %s, wanted - %s", blockByHash, block1)
		}

		blockByIndex, err := repository.BlockByIndex(block1.Index())
		if err != nil {
			t.Fatalf("failed to read block by index: %s", err)
		}
		if !blockByIndex.IsEqual(block1) {
			t.Errorf("expected read block to equal block1. Got - %s, wanted - %s", blockByIndex, block1)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("CanReadSavedBlocks", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("CanReadSavedBlocks")

		suite(repository, t)
	})
}

func TestRepository_ReturnsNilIfHeaderDoesNotExist(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		headerByHash, err := repository.HeaderByHash(block1.Hash())
		if err != nil {
			t.Fatalf("failed to read header by hash: %s", err)
		}
		if headerByHash != nil {
			t.Errorf("expected header to be nil. Got - %s, wanted - %s", headerByHash, nil)
		}

		headerByIndex, err := repository.HeaderByIndex(block1.Index())
		if err != nil {
			t.Fatalf("failed to read header by index: %s", err)
		}
		if headerByIndex != nil {
			t.Errorf("expected header to be nil. Got - %s, wanted - %s", headerByIndex, nil)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("ReturnsNilIfHeaderDoesNotExist", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("ReturnsNilIfHeaderDoesNotExist")

		suite(repository, t)
	})
}

func TestRepository_CanReadSavedHeaders(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		header1 := block1.Header()

		err := repository.SaveBlock(block1)
		if err != nil {
			t.Fatalf("failed to save block1: %s", err)
		}

		headerByHash, err := repository.HeaderByHash(header1.Hash)
		if err != nil {
			t.Fatalf("failed to read header by hash: %s", err)
		}
		if !headerByHash.IsEqual(header1) {
			t.Errorf("expected headers to be equal. Got - %s, wanted - %s", headerByHash, header1)
		}

		headerByIndex, err := repository.HeaderByIndex(header1.Index)
		if err != nil {
			t.Fatalf("failed to read header by index: %s", err)
		}
		if !headerByIndex.IsEqual(header1) {
			t.Errorf("expected headers to be equal. Got - %s, wanted - %s", headerByIndex, header1)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("CanReadSavedHeaders", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("CanReadSavedHeaders")

		suite(repository, t)
	})
}

func TestRepository_SavingBlocksAlsoSavesHeader(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		err := repository.SaveBlock(block1)
		if err != nil {
			t.Fatalf("failed to save block1: %s", err)
		}

		header1 := block1.Header()

		headerByHash, err := repository.HeaderByHash(header1.Hash)
		if err != nil {
			t.Fatalf("failed to read header by hash: %s", err)
		}
		if !headerByHash.IsEqual(header1) {
			t.Errorf("expected headers to be equal. Got - %s, wanted - %s", headerByHash, header1)
		}

		headerByIndex, err := repository.HeaderByIndex(header1.Index)
		if err != nil {
			t.Fatalf("failed to read header by index: %s", err)
		}
		if !headerByIndex.IsEqual(header1) {
			t.Errorf("expected headers to be equal. Got - %s, wanted - %s", headerByIndex, header1)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("SavingBlocksAlsoSavesHeader", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("SavingBlocksAlsoSavesHeader")

		suite(repository, t)
	})
}

func TestRepository_BestBlockReturnsHighestBlock(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		err := repository.SaveBlock(block1)
		if err != nil {
			t.Fatalf("failed to save block1: %s", err)
		}

		err = repository.SaveBlock(block2)
		if err != nil {
			t.Fatalf("failed to save block2: %s", err)
		}

		bestBlock, err := repository.BestBlock()
		if err != nil {
			t.Fatalf("failed to read best block: %s", err)
		}

		if !bestBlock.IsEqual(block2) {
			t.Errorf("expected best block to equal block2. Got - %s, wanted - %s", bestBlock, block2)
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("BestBlockReturnsHighestBlock", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("BestBlockReturnsHighestBlock")

		suite(repository, t)
	})
}

func TestRepository_BestHeaderReturnsHighestIndexHeader(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		err := repository.SaveHeader(block1.Header())
		if err != nil {
			t.Fatalf("failed to save header1: %s", err)
		}

		err = repository.SaveHeader(block2.Header())
		if err != nil {
			t.Fatalf("failed to save header2: %s", err)
		}

		bestHeader, err := repository.BestHeader()
		if err != nil {
			t.Fatalf("failed to read best header: %s", err)
		}

		if !bestHeader.IsEqual(block2.Header()) {
			t.Errorf("expected best header to equal header2. Got - %s, wanted - %s", bestHeader, block2.Header())
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("BestHeaderReturnsHighestIndexHeader", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("BestHeaderReturnsHighestIndexHeader")

		suite(repository, t)
	})
}

func TestRepository_BestHeaderCanBeHigherThanBestBlock(t *testing.T) {
	suite := func(repository entity.ChainRepository, t *testing.T) {
		err := repository.SaveBlock(block1)
		if err != nil {
			t.Fatalf("failed to save block1: %s", err)
		}

		bestBlock1, err := repository.BestBlock()
		if err != nil {
			t.Fatalf("failed to read best header: %s", err)
		}
		if !bestBlock1.IsEqual(block1) {
			t.Errorf("expected best header to equal block1. Got - %s, wanted - %s", bestBlock1, block1)
		}

		bestHeader1, err := repository.BestHeader()
		if err != nil {
			t.Fatalf("failed to read best header: %s", err)
		}
		if !bestHeader1.IsEqual(block1.Header()) {
			t.Errorf("expected best header to equal header1. Got - %s, wanted - %s", bestHeader1, block1.Header())
		}

		err = repository.SaveHeader(block2.Header())
		if err != nil {
			t.Fatalf("failed to save header2: %s", err)
		}

		bestBlock1, err = repository.BestBlock()
		if err != nil {
			t.Fatalf("failed to read best header: %s", err)
		}
		if !bestBlock1.IsEqual(block1) {
			t.Errorf("expected best header to equal block1. Got - %s, wanted - %s", bestBlock1, block1)
		}

		bestHeader2, err := repository.BestHeader()
		if err != nil {
			t.Fatalf("failed to read best header: %s", err)
		}
		if !bestHeader2.IsEqual(block2.Header()) {
			t.Errorf("expected best header to equal header2. Got - %s, wanted - %s", bestHeader2, block2.Header())
		}
	}

	t.Run("memdb", func(t *testing.T) {
		repository := memdb.NewChainRepository()

		suite(repository, t)
	})
	t.Run("boltdb", func(t *testing.T) {
		repository := boltdb.Builder("BestHeaderCanBeHigherThanBestBlock", blockEncoder).
			WithCache().
			WithMetrics().
			WithLogger().
			Build()
		defer repository.Close()
		defer boltdb.DeleteBlockchain("BestHeaderCanBeHigherThanBestBlock")

		suite(repository, t)
	})
}
