package blockchain

import "math/big"

type blockSolution struct {
	nonce int
	hash  []byte
}

func miner(block *Block, nonces <-chan int, solutions chan<- *blockSolution) {
	for nonce := range nonces {
		hash := hashBlock(block, nonce)

		if new(big.Int).SetBytes(hash[:]).Cmp(target) == -1 {
			solutions <- &blockSolution{nonce, hash[:]}
		}
	}
}
