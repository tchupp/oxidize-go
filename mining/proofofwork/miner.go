package proofofwork

import (
	"math"
	"runtime"
	"log"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/block"
)

const (
	maxNonce = math.MaxInt64
)

var (
	defaultWorkerCount = uint(runtime.NumCPU())
)

type Miner struct {
	workerCount uint
}

func NewMiner(workerCount uint) *Miner {
	return &Miner{workerCount: workerCount}
}

func NewDefaultMiner() *Miner {
	return NewMiner(defaultWorkerCount)
}

func (miner *Miner) MineBlock(header *block.Header) (*block.Block) {
	solutions := make(chan *block.Solution)
	nonces := make(chan int, miner.workerCount)
	defer close(nonces)

	for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
		go worker(header, nonces, solutions)
	}

	for nonce := 0; nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			return block.NewBlock(header, solution)
		default:
			nonces <- nonce
		}
	}

	log.Panic(blockchain.MaxNonceOverflowError)
	return nil
}

func worker(header *block.Header, nonces <-chan int, solutions chan<- *block.Solution) {
	for nonce := range nonces {
		hash := block.CalculateHash(header, nonce)

		if block.HashValid(hash) {
			solutions <- &block.Solution{nonce, hash}
		}
	}
}
