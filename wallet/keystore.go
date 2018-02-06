package wallet

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tclchiam/oxidize-go/crypto"
	"github.com/tclchiam/oxidize-go/identity"
)

const ecPrivateKeyType = "EC PRIVATE KEY"

type KeyStore struct {
	path string
}

func NewKeyStore(path string) *KeyStore {
	return &KeyStore{path: path}
}

func (store *KeyStore) GetIdentity(address string) (*identity.Identity, error) {
	filename := buildPemFilename(store.path, address)

	block, err := readBlockFromFile(filename)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return identity.NewIdentity((*crypto.PrivateKey)(privateKey)), nil
}

func (store *KeyStore) SaveIdentity(identity *identity.Identity) error {
	filename := buildPemFilename(store.path, identity.Address())

	bytes, err := x509.MarshalECPrivateKey(identity.PrivateKey().ToECDSA())
	if err != nil {
		return err
	}
	content := pem.EncodeToMemory(&pem.Block{
		Type:  ecPrivateKeyType,
		Bytes: bytes,
	})

	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return err
	}

	file, err := ioutil.TempFile(filepath.Dir(filename), fmt.Sprintf(".%s.tmp", filepath.Base(filename)))
	if err != nil {
		return err
	}

	if _, err := file.Write(content); err != nil {
		file.Close()
		os.Remove(file.Name())
		return err
	}
	file.Close()
	return os.Rename(file.Name(), filename)
}

func readBlockFromFile(filename string) (*pem.Block, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("file '%s' did not contain a pem encoded key", filename)
	}
	return block, nil
}

func buildPemFilename(path, address string) string {
	return filepath.Join(path, fmt.Sprintf("address-%s.pem", address))
}
