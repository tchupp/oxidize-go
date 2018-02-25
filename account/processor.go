package account

import (
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/identity"
)

type blockProcessor struct {
	bc   blockchain.Blockchain
	repo *accountRepo
	quit chan struct{}
}

func NewBlockProcessor(bc blockchain.Blockchain, repo *accountRepo) *blockProcessor {
	return &blockProcessor{
		bc:   bc,
		quit: make(chan struct{}),
		repo: repo,
	}
}

func (p *blockProcessor) Start() {
	go p.eventLoop()
}

func (p *blockProcessor) eventLoop() {
	c := make(chan blockchain.Event, 10)
	sub := p.bc.Subscribe(c)
	defer sub.Unsubscribe()

	log.Debugf("starting indexing...")
	currentIndex := p.processNewBlocks(0)
	log.Debugf("caught up with index '%d', waiting for event...", currentIndex)

	for {
		select {
		case <-p.quit:
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
			currentIndex = p.processNewBlocks(currentIndex)
			log.Debugf("caught up with index '%d', waiting for event...", currentIndex)
		}
	}
}

func (p *blockProcessor) processNewBlocks(currentIndex uint64) uint64 {
	bestBlock, err := p.bc.BestBlock()
	for ; err == nil && currentIndex <= bestBlock.Index(); currentIndex += 1 {
		block, err := p.bc.BlockByIndex(currentIndex)
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

		for _, tx := range block.Transactions() {
			var spenderAddress *identity.Address
			if len(tx.Inputs) > 0 {
				spenderAddress = identity.FromPublicKey(tx.Inputs[0].PublicKey)
			}

			for _, out := range tx.Outputs {
				accountTx := &Transaction{
					amount:   out.Value,
					spender:  spenderAddress,
					receiver: identity.FromPublicKeyHash(out.PublicKeyHash),
				}

				p.repo.SaveTx(accountTx.spender, accountTx)
				p.repo.SaveTx(accountTx.receiver, accountTx)
			}
		}
	}

	return currentIndex
}

func (p *blockProcessor) Close() error {
	p.quit <- struct{}{}
	return nil
}
