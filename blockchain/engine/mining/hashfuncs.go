package mining

import (
	"crypto/sha256"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func calculateHash(b []byte) entity.Hash {
	return entity.Hash(sha256.Sum256(b))
}
