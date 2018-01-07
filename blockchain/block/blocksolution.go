package block

import "github.com/tclchiam/block_n_go/blockchain/chainhash"

type Solution struct {
	Nonce int
	Hash  chainhash.Hash
}
