package account

import (
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

type BlockProcessor interface {
	Process(block *entity.Block)
}

type BlockProcessors []BlockProcessor

func (ps BlockProcessors) Process(block *entity.Block) {
	for _, p := range ps {
		p.Process(block)
	}
}

type chainIndexer struct {
	bc         blockchain.Blockchain
	processors BlockProcessors
	quit       chan struct{}
}

func NewChainIndexer(bc blockchain.Blockchain, repo *accountRepo) *chainIndexer {
	indexer := &chainIndexer{
		bc:         bc,
		processors: BlockProcessors{NewAccountUpdater(repo)},
		quit:       make(chan struct{}),
	}

	go indexer.indexingLoop()

	return indexer
}

func (indexer *chainIndexer) indexingLoop() {
	c := make(chan blockchain.Event, 10)
	sub := indexer.bc.Subscribe(c)
	defer sub.Unsubscribe()

	log.Debugf("starting indexing...")
	currentIndex := indexer.handleNewBlocks(0)
	log.Debugf("caught up with index '%d', waiting for event...", currentIndex)

	for {
		select {
		case <-indexer.quit:
			log.Debug("quitting indexing")
			return
		case event, ok := <-c:
			if !ok {
				log.Debug("quitting indexing")
				return
			}
			if event != blockchain.BlockSaved {
				break
			}

			log.Debugf("resuming indexing with index '%d'...", currentIndex)
			currentIndex = indexer.handleNewBlocks(currentIndex)
			log.Debugf("caught up with index '%d', waiting for event...", currentIndex)
		}
	}
}

func (indexer *chainIndexer) handleNewBlocks(currentIndex uint64) uint64 {
	bestBlock, err := indexer.bc.BestBlock()
	for ; err == nil && currentIndex <= bestBlock.Index(); currentIndex += 1 {
		block, err := indexer.bc.BlockByIndex(currentIndex)
		if err != nil {
			log.WithField("index", currentIndex).
				WithError(err).
				Info("error reading block")
			return currentIndex
		}
		if block == nil {
			log.WithField("index", currentIndex).
				Infof("block is nil")
			return currentIndex
		}

		indexer.processors.Process(block)
	}

	return currentIndex
}

func (indexer *chainIndexer) Close() error {
	indexer.quit <- struct{}{}
	return nil
}
