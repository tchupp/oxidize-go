package proofofwork

import (
	"math"
	"runtime"
	"log"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
)

const (
	maxNonce = math.MaxInt64
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

func (miner *miner) MineBlock(header *entity.BlockHeader) (*entity.Block) {
	solutions := make(chan *entity.BlockSolution)
	nonces := make(chan int, miner.workerCount)
	defer close(nonces)

	go func() {
		for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
			go worker(header, nonces, solutions)
		}
	}()

	for nonce := 0; nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			return entity.NewBlock(header, solution)
		default:
			nonces <- nonce
		}
	}

	log.Panic(blockchain.MaxNonceOverflowError)
	return nil
}

func worker(header *entity.BlockHeader, nonces <-chan int, solutions chan<- *entity.BlockSolution) {
	for nonce := range nonces {
		hash := mining.CalculateHash(header, nonce)

		if mining.HashValid(hash) {
			solutions <- &entity.BlockSolution{nonce, hash}
		}
	}
}
