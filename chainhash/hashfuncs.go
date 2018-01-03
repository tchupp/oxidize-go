package chainhash

import "crypto/sha256"

func CalculateHash(b []byte) Hash {
	return Hash(sha256.Sum256(b))
}
