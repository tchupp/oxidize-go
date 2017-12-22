package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
	"github.com/mr-tron/base58/base58"
)

const (
	version               = byte(0x00)
	addressChecksumLength = 4
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

func (w *Wallet) GetAddress() string {
	publicKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, publicKeyHash...)
	checksum := checksum(versionedPayload)

	rawAddress := append(versionedPayload, checksum...)
	return base58.Encode(rawAddress)
}

func AddressToPublicKeyHash(address string) []byte {
	rawAddress, err := base58.Decode(address)
	if err != nil {
		log.Panic(err)
	}
	return rawAddress[1: len(rawAddress)-addressChecksumLength]
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	hashImpl := ripemd160.New()
	_, err := hashImpl.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	return hashImpl.Sum(nil)
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLength]
}
