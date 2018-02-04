package mining

import (
	"crypto/sha256"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func calculateHash(b []byte) entity.Hash {
	return entity.Hash(sha256.Sum256(b))
}
