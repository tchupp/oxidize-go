package wallet

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/tclchiam/oxidize-go/crypto"
	"github.com/tclchiam/oxidize-go/identity"
)

const ecPrivateKeyType = "EC PRIVATE KEY"
const keyFileTemplate = "address-%s.pem"

var keyFileRegex = regexp.MustCompile("address-([a-zA-Z0-9]+)\\.pem")

type KeyStore struct {
	path string
}

func NewKeyStore(path string) *KeyStore {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Clean(path), 0700)
	}
	return &KeyStore{path: path}
}

func (store *KeyStore) Identities() (identity.Identities, error) {
	infos, err := ioutil.ReadDir(store.path)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, info := range infos {
		if keyFileRegex.MatchString(info.Name()) {
			addr := keyFileRegex.ReplaceAllString(info.Name(), "${1}")
			addrs = append(addrs, addr)
		}
	}

	var result *multierror.Error
	var ids identity.Identities
	for _, addr := range addrs {
		id, err := store.Identity(addr)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		ids = append(ids, id)
	}

	return ids, result.ErrorOrNil()
}

func (store *KeyStore) Identity(address string) (*identity.Identity, error) {
	filename := buildPemFilename(store.path, address)

	block, err := readBlockFromFile(filename)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return identity.NewIdentity(crypto.FromPrivateKey(privateKey)), nil
}

func (store *KeyStore) SaveIdentity(identity *identity.Identity) error {
	filename := buildPemFilename(store.path, identity.Address().Serialize())

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
	return filepath.Join(path, fmt.Sprintf(keyFileTemplate, address))
}
