package proofofwork

import (
	"log"
	"math"
	"runtime"
	"time"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
)

const (
	maxNonce = math.MaxUint64
)

var (
	defaultWorkerCount = uint(runtime.NumCPU())
)

type miner struct {
	workerCount uint
}

func NewMiner(workerCount uint) mining.Miner {
	return &miner{workerCount: workerCount}
}

func NewDefaultMiner() mining.Miner {
	return NewMiner(defaultWorkerCount)
}

func (miner *miner) MineBlock(parent *entity.Block, transactions entity.Transactions) (*entity.Block) {
	return miner.mineBlock(parent, transactions, uint64(time.Now().Unix()))
}

func (miner *miner) mineBlock(parent *entity.Block, transactions entity.Transactions, now uint64) (*entity.Block) {
	transactionsHash := mining.CalculateTransactionsHash(transactions)

	work := &mining.BlockHashingInput{
		Index:            parent.Index() + 1,
		PreviousHash:     parent.Hash(),
		Timestamp:        now,
		TransactionsHash: transactionsHash,
	}

	solutions := make(chan *entity.BlockSolution)
	nonces := make(chan uint64, miner.workerCount)
	defer close(nonces)

	go func() {
		for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
			go worker(work, nonces, solutions)
		}
	}()

	for nonce := uint64(0); nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			header := entity.NewBlockHeader(work.Index, work.PreviousHash, work.TransactionsHash, work.Timestamp, solution.Nonce, solution.Hash)
			return entity.NewBlock(header, transactions)
		default:
			nonces <- nonce
		}
	}

	log.Panic(MaxNonceOverflowError)
	return nil
}

func worker(work *mining.BlockHashingInput, nonces <-chan uint64, solutions chan<- *entity.BlockSolution) {
	for nonce := range nonces {
		hash := mining.CalculateBlockHash(work, nonce)

		if mining.HashValid(hash) {
			solutions <- &entity.BlockSolution{Nonce: nonce, Hash: hash}
		}
	}
}
