package identity

import (
	"bytes"
	"fmt"

	"github.com/mr-tron/base58/base58"
)

const addressLength = 34

type Address struct {
	version       byte
	publicKeyHash []byte
	checksum      []byte
}

func DeserializeAddress(data string) (*Address, error) {
	if len(data) != addressLength {
		return nil, fmt.Errorf("unexpected address size: %d", len(data))
	}

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

func (a *Address) String() string              { return a.Serialize() }
func (a *Address) Version() byte               { return a.version }
func (a *Address) PublicKeyHash() []byte       { return a.publicKeyHash }
func (a *Address) Checksum() []byte            { return a.checksum }
func (a *Address) IsEqual(other *Address) bool { return a.String() == other.String() }
