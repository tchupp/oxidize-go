package wallet

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tclchiam/oxidize-go/identity"
)

func TestKeyStore_SaveIdentity(t *testing.T) {
	randomIdentity := identity.RandomIdentity()
	keyStore := NewKeyStore(makeKeystoreDir())
	defer os.RemoveAll(keyStore.path)

	err := keyStore.SaveIdentity(randomIdentity)
	if err != nil {
		t.Fatalf("saving identity: %s", err)
	}

	readIdentity, err := keyStore.Identity(randomIdentity.Address().Serialize())
	if err != nil {
		t.Fatalf("reading identity: %s", err)
	}

	if !readIdentity.IsEqual(randomIdentity) {
		t.Errorf("read different identity than expected. \nwrote: '%x', \n read: '%x'", randomIdentity, readIdentity)
	}
}

func TestKeyStore_ListIdentity(t *testing.T) {
	randomIdentity1 := identity.RandomIdentity()
	randomIdentity2 := identity.RandomIdentity()
	keyStore := NewKeyStore(makeKeystoreDir())
	defer os.RemoveAll(keyStore.path)

	err := keyStore.SaveIdentity(randomIdentity1)
	if err != nil {
		t.Fatalf("saving identity: %s", err)
	}
	err = keyStore.SaveIdentity(randomIdentity2)
	if err != nil {
		t.Fatalf("saving identity: %s", err)
	}

	identities, err := keyStore.Identities()
	if err != nil {
		t.Fatalf("reading identity list: %s", err)
	}

	if len(identities) != 2 {
		t.Errorf("unexpected number of identities. got - %d, wanted - %d", len(identities), 2)
	}
	for _, id := range identities {
		if !id.IsEqual(randomIdentity1) && !id.IsEqual(randomIdentity2) {
			t.Errorf("unexpected identity. got - %s. not %s or %s", id, randomIdentity1, randomIdentity2)
		}
	}
}

func makeKeystoreDir() string {
	secret := make([]byte, 10)
	rand.Read(secret)
	return filepath.Join(os.TempDir(), fmt.Sprintf("%x", secret))
}
