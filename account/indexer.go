package account

import (
	"fmt"
	"sync"

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

type IndexerStatus uint

const (
	Starting IndexerStatus = iota
	Idle
	Syncing
	Done
)

func (s IndexerStatus) String() string {
	switch s {
	case Starting:
		return "Starting"
	case Idle:
		return "Idle"
	case Syncing:
		return "Syncing"
	case Done:
		return "Done"
	default:
		return ""
	}
}

type chainIndexer struct {
	bc         blockchain.Blockchain
	processors BlockProcessors
	lock       sync.RWMutex
	status     IndexerStatus
	quit       chan struct{}
}

func NewChainIndexer(bc blockchain.Blockchain, processors ...BlockProcessor) *chainIndexer {
	indexer := &chainIndexer{
		bc:         bc,
		processors: processors,
		status:     Starting,
		quit:       make(chan struct{}),
	}

	go indexer.indexingLoop()

	return indexer
}

func (indexer *chainIndexer) Status() IndexerStatus {
	indexer.lock.RLock()
	defer indexer.lock.RUnlock()

	return indexer.status
}

func (indexer *chainIndexer) updateStatus(status IndexerStatus, msg string) {
	indexer.lock.Lock()
	defer indexer.lock.Unlock()

	indexer.status = status
	log.WithField("state", status).Debug(msg)
}

func (indexer *chainIndexer) indexingLoop() {
	c := make(chan blockchain.Event, 10)
	sub := indexer.bc.Subscribe(c)
	defer sub.Unsubscribe()

	indexer.updateStatus(Syncing, "starting indexing...")
	currentIndex := indexer.handleNewBlocks(0)
	indexer.updateStatus(Idle, fmt.Sprintf("caught up with index '%d', waiting for event...", currentIndex))

	for {
		select {
		case <-indexer.quit:
			indexer.updateStatus(Done, "quitting indexing")
			return
		case event, ok := <-c:
			indexer.updateStatus(Syncing, fmt.Sprintf("resuming indexing with index '%d'...", currentIndex))

			if !ok {
				indexer.updateStatus(Done, "quitting indexing")
				return
			}
			if event == blockchain.BlockSaved {
				currentIndex = indexer.handleNewBlocks(currentIndex)
			}
			indexer.updateStatus(Idle, fmt.Sprintf("caught up with index '%d', waiting for event...", currentIndex))
		}
	}
}

func (indexer *chainIndexer) handleNewBlocks(currentIndex uint64) uint64 {
	bestBlock, err := indexer.bc.BestBlock()
	if err != nil {
		return currentIndex
	}

	for ; currentIndex <= bestBlock.Index(); currentIndex += 1 {
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
