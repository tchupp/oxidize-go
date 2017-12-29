package tx

import "crypto/sha256"

func Hash(tx *Transaction) []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(serialize(&txCopy))

	return hash[:]
}
