package identity

import (
	"bytes"
	"crypto/sha256"

	"github.com/mr-tron/base58/base58"
	"github.com/tclchiam/block_n_go/crypto"
	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x00)
	checksumLength = 4
)

type Address struct {
	version       byte
	publicKeyHash []byte
	checksum      []byte

	publicKey  *crypto.PublicKey
	privateKey *crypto.PrivateKey
}

func RandomAddress() *Address {
	publicKey := crypto.NewP256PrivateKey()
	return NewAddress(publicKey)
}

func NewAddress(privateKey *crypto.PrivateKey) *Address {
	publicKey := privateKey.PubKey()
	publicKeyHash := hashPublicKey(publicKey)
	checksum := checksum(publicKeyHash)

	return &Address{
		version:       version,
		publicKeyHash: publicKeyHash,
		checksum:      checksum,
		publicKey:     publicKey,
		privateKey:    privateKey,
	}
}

func hashPublicKey(publicKey *crypto.PublicKey) []byte {
	publicSHA256 := sha256.Sum256(publicKey.Serialize())

	hashImpl := ripemd160.New()
	hashImpl.Write(publicSHA256[:])
	hash := hashImpl.Sum(nil)
	return hash
}

func checksum(publicKeyHash []byte) []byte {
	payload := append([]byte{version}, publicKeyHash...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checksumLength]
}

func (a *Address) Base58() string {
	input := [][]byte{
		{version},
		a.publicKeyHash,
		a.checksum,
	}
	return base58.Encode(bytes.Join(input, []byte{}))
}

func (a *Address) String() string        { return a.Base58() }
func (a *Address) Version() byte         { return a.version }
func (a *Address) PublicKeyHash() []byte { return a.publicKeyHash }
func (a *Address) Checksum() []byte      { return a.checksum }
