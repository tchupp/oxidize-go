package chainhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const HashSize = sha256.Size

const MaxHashStringSize = HashSize * 2

var ErrHashStrSize = fmt.Errorf("max hash string length is %v bytes", MaxHashStringSize)

type Hash [HashSize]byte

var EmptyHash = Hash([HashSize]byte{0})

func (hash Hash) String() string {
	return hex.EncodeToString(hash[:])
}

func (hash Hash) Slice() []byte {
	return hash[:]
}

func (hash *Hash) IsEqual(target *Hash) bool {
	if hash == nil && target == nil {
		return true
	}
	if hash == nil || target == nil {
		return false
	}
	return *hash == *target
}

func NewHash(newHash []byte) (*Hash, error) {
	var hash Hash
	if len(newHash) != HashSize {
		return nil, fmt.Errorf("invalid hash length of %v, want %v", len(newHash), HashSize)
	}
	copy(hash[:], newHash)

	return &hash, nil
}

func NewHashFromStr(hash string) (*Hash, error) {
	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	for len(bytes) < HashSize {
		bytes = append([]byte{0}, bytes...)
	}

	return NewHash(bytes)
}
