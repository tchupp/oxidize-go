package identity_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/mr-tron/base58/base58"
	"github.com/tclchiam/oxidize-go/identity"
)

func TestIdentity_Address(t *testing.T) {
	id := identity.RandomIdentity()
	publicKey := id.PublicKey()

	input := [][]byte{
		{0x00},
		publicKey.Hash(),
		checksum(publicKey.Hash()),
	}

	expectedBase58 := base58.Encode(bytes.Join(input, []byte{}))
	if expectedBase58 != id.Address().Serialize() {
		t.Errorf("Expected base58 did not equal actual. Got: '%s', wanted: '%s'", id.Address(), expectedBase58)
	}
}

func checksum(publicKeyHash []byte) []byte {
	payload := append([]byte{0x00}, publicKeyHash...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]
}
