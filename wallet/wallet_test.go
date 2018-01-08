package wallet

import (
	"testing"
	"bytes"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func TestHashPubKey(t *testing.T) {
	wallet := NewWallet()
	walletPublicKeyHash := HashPubKey(wallet.PublicKey)

	publicSHA256 := sha256.Sum256(wallet.PublicKey.Serialize())

	hashImpl := ripemd160.New()
	_, err := hashImpl.Write(publicSHA256[:])
	if err != nil {
		t.Errorf("error writing ripemd160 hash: %s", err)
	}

	expectedHash := hashImpl.Sum(nil)

	if bytes.Compare(walletPublicKeyHash, expectedHash) != 0 {
		t.Errorf("public key hash from wallet '%x' did not match decoded public key hash from address '%x'", walletPublicKeyHash, expectedHash)
	}
}

func TestAddressToPublicKeyHash(t *testing.T) {
	wallet := NewWallet()
	walletPublicKeyHash := HashPubKey(wallet.PublicKey)

	address := wallet.GetAddress()
	addressPublicKeyHash, err := AddressToPublicKeyHash(address)
	if err != nil {
		t.Errorf("error converting address to public key hash: %s", err)
	}

	if bytes.Compare(walletPublicKeyHash, addressPublicKeyHash) != 0 {
		t.Errorf("public key hash from wallet '%x' did not match decoded public key hash from address '%x'", walletPublicKeyHash, addressPublicKeyHash)
	}
}
