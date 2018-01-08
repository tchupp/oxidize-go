package wallet

import (
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
	"github.com/mr-tron/base58/base58"
	"github.com/tclchiam/block_n_go/crypto"
)

const (
	version               = byte(0x00)
	addressChecksumLength = 4
)

type Wallet struct {
	PrivateKey *crypto.PrivateKey
	PublicKey  *crypto.PublicKey
}

func NewWallet() *Wallet {
	privateKey := crypto.NewP256PrivateKey()
	return &Wallet{PrivateKey: privateKey, PublicKey: privateKey.PubKey()}
}

func (w *Wallet) GetAddress() string {
	publicKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, publicKeyHash...)
	checksum := checksum(versionedPayload)

	rawAddress := append(versionedPayload, checksum...)
	return base58.Encode(rawAddress)
}

func AddressToPublicKeyHash(address string) ([]byte, error) {
	rawAddress, err := base58.Decode(address)
	if err != nil {
		return nil, err
	}
	return rawAddress[1: len(rawAddress)-addressChecksumLength], nil
}

func HashPubKey(publicKey *crypto.PublicKey) []byte {
	publicSHA256 := sha256.Sum256(publicKey.Serialize())

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
