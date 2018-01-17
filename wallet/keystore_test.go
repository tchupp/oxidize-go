package wallet

import (
	"os"
	"testing"
	"github.com/tclchiam/block_n_go/identity"
)

func TestKeyStore_SaveAccount(t *testing.T) {
	randomIdentity := identity.RandomIdentity()
	keyStore := NewKeyStore(os.TempDir())

	err := keyStore.SaveIdentity(randomIdentity)
	if err != nil {
		t.Fatalf("saving identity: %s", err)
	}

	readIdentity, err := keyStore.GetIdentity(randomIdentity.Address())
	if err != nil {
		t.Fatalf("reading identity: %s", err)
	}

	if !readIdentity.IsEqual(randomIdentity) {
		t.Errorf("read different identity than expected. \nwrote: '%x', \n read: '%x'", randomIdentity, readIdentity)
	}
}
