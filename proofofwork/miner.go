package proofofwork

import (
	"math"
	"runtime"
	"log"
	"github.com/tclchiam/block_n_go/blockchain"
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
	return &Miner{workerCount: defaultWorkerCount}
}

func (miner *Miner) MineBlock(header *blockchain.BlockHeader) (*blockchain.Block) {
	solutions := make(chan *blockchain.BlockSolution)
	nonces := make(chan int, miner.workerCount)
	defer close(nonces)

	for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
		go worker(header, nonces, solutions)
	}

	for nonce := 0; nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			return blockchain.NewBlock(header, solution)
		default:
			nonces <- nonce
		}
	}

	log.Panic(blockchain.MaxNonceOverflowError)
	return nil
}

func worker(header *blockchain.BlockHeader, nonces <-chan int, solutions chan<- *blockchain.BlockSolution) {
	for nonce := range nonces {
		hash := blockchain.CalculateBlockHash(header, nonce)

		if blockchain.BlockHashValid(hash) {
			solutions <- &blockchain.BlockSolution{nonce, hash}
		}
	}
}
