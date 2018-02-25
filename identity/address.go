package identity

import (
	"bytes"
	"crypto/sha256"

	"github.com/mr-tron/base58/base58"
	"github.com/tclchiam/oxidize-go/crypto"
)

const (
	version        = byte(0x00)
	checksumLength = 4
)

type Address struct {
	version       byte
	publicKeyHash []byte
	checksum      []byte
}

func FromPublicKey(publicKey *crypto.PublicKey) *Address {
	return FromPublicKeyHash(publicKey.Hash())
}

func FromPublicKeyHash(publicKeyHash []byte) *Address {
	checksum := checksum(publicKeyHash)

	return &Address{
		version:       version,
		publicKeyHash: publicKeyHash,
		checksum:      checksum,
	}
}

func NewAddress(version byte, publicKeyHash []byte, checksum []byte) *Address {
	return &Address{
		version:       version,
		publicKeyHash: publicKeyHash,
		checksum:      checksum,
	}
}

func DeserializeAddress(data string) (*Address, error) {
	b, err := base58.Decode(data)
	if err != nil {
		return nil, err
	}

	version := b[0]
	publicKeyHash := b[1 : len(b)-checksumLength]
	checksum := b[len(b)-checksumLength:]

	return &Address{
		version:       version,
		publicKeyHash: publicKeyHash,
		checksum:      checksum,
	}, nil
}

func (a *Address) Serialize() string {
	input := [][]byte{
		{a.version},
		a.publicKeyHash,
		a.checksum,
	}
	return base58.Encode(bytes.Join(input, []byte{}))
}

func (a *Address) String() string        { return a.Serialize() }
func (a *Address) Version() byte         { return a.version }
func (a *Address) PublicKeyHash() []byte { return a.publicKeyHash }
func (a *Address) Checksum() []byte      { return a.checksum }

func (a *Address) IsEqual(other *Address) bool {
	if a == nil && other == nil {
		return true
	}
	if a == nil || other == nil {
		return false
	}

	return a.Serialize() == other.Serialize()
}

func checksum(publicKeyHash []byte) []byte {
	payload := append([]byte{version}, publicKeyHash...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checksumLength]
}
