package proofofwork

import (
	"math"
	"runtime"
	"log"

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

func (miner *miner) MineBlock(header *entity.BlockHeader) (*entity.BlockSolution) {
	solutions := make(chan *entity.BlockSolution)
	nonces := make(chan uint64, miner.workerCount)
	defer close(nonces)

	go func() {
		for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
			go worker(header, nonces, solutions)
		}
	}()

	for nonce := uint64(0); nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			return solution
		default:
			nonces <- nonce
		}
	}

	log.Panic(MaxNonceOverflowError)
	return nil
}

func worker(header *entity.BlockHeader, nonces <-chan uint64, solutions chan<- *entity.BlockSolution) {
	for nonce := range nonces {
		hash := mining.CalculateHash(header, nonce)

		if mining.HashValid(hash) {
			solutions <- &entity.BlockSolution{Nonce: nonce, Hash: &hash}
		}
	}
}
