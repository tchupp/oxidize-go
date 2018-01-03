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

func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

func (hash *Hash) Slice() []byte {
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
	ret := new(Hash)
	err := decode(ret, hash)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func decode(dst *Hash, src string) error {
	if len(src) > MaxHashStringSize {
		return ErrHashStrSize
	}

	var srcBytes []byte
	if len(src)%2 == 0 {
		srcBytes = []byte(src)
	} else {
		srcBytes = make([]byte, 1+len(src))
		srcBytes[0] = '0'
		copy(srcBytes[1:], src)
	}

	var reversedHash Hash
	_, err := hex.Decode(reversedHash[HashSize-hex.DecodedLen(len(srcBytes)):], srcBytes)
	if err != nil {
		return err
	}

	for i, b := range reversedHash[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = reversedHash[HashSize-1-i], b
	}

	return nil
}
