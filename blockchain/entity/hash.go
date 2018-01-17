package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

const (
	hashLength         = sha256.Size
	maxHexStringLength = hashLength * 2
)

type Hash [hashLength]byte

var EmptyHash = Hash{}

func (hash Hash) String() string { return hex.EncodeToString(hash[:]) }
func (hash Hash) Slice() []byte  { return hash[:] }
func (hash Hash) IsEmpty() bool  { return hash == EmptyHash }

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
	if len(newHash) != hashLength {
		return nil, fmt.Errorf("invalid hash length of %v, want %v", len(newHash), hashLength)
	}
	copy(hash[:], newHash)

	return &hash, nil
}

func NewHashFromString(hash string) (*Hash, error) {
	if len(hash) > maxHexStringLength {
		hash = hash[:maxHexStringLength]
	}
	if len(hash)%2 != 0 {
		hash = "0" + hash
	}

	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	for len(bytes) < hashLength {
		bytes = append([]byte{0}, bytes...)
	}

	return NewHash(bytes)
}

func NewHashOrPanic(newHash string) *Hash {
	hash, err := NewHashFromString(newHash)
	if err != nil {
		log.Panic(err)
	}

	return hash
}
