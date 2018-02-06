package identity

import (
	"bytes"
	"crypto/sha256"

	"github.com/mr-tron/base58/base58"
	"github.com/tclchiam/oxidize-go/crypto"
	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x00)
	checksumLength = 4
)

type Identity struct {
	version       byte
	publicKeyHash []byte
	checksum      []byte

	publicKey  *crypto.PublicKey
	privateKey *crypto.PrivateKey
}

func RandomIdentity() *Identity {
	publicKey := crypto.NewP256PrivateKey()
	return NewIdentity(publicKey)
}

func NewIdentity(privateKey *crypto.PrivateKey) *Identity {
	publicKey := privateKey.PubKey()
	publicKeyHash := hashPublicKey(publicKey)
	checksum := checksum(publicKeyHash)

	return &Identity{
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

func (a *Identity) Address() string {
	input := [][]byte{
		{version},
		a.publicKeyHash,
		a.checksum,
	}
	return base58.Encode(bytes.Join(input, []byte{}))
}

func (a *Identity) Sign(data []byte) (*crypto.Signature, error) {
	return a.privateKey.Sign(data)
}

func (a *Identity) String() string                 { return a.Address() }
func (a *Identity) Version() byte                  { return a.version }
func (a *Identity) PrivateKey() *crypto.PrivateKey { return a.privateKey }
func (a *Identity) PublicKey() *crypto.PublicKey   { return a.publicKey }
func (a *Identity) PublicKeyHash() []byte          { return a.publicKeyHash }
func (a *Identity) Checksum() []byte               { return a.checksum }
func (a *Identity) IsEqual(other *Identity) bool   { return a.Address() == other.Address() }
