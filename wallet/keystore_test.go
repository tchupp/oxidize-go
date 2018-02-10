package wallet

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tclchiam/oxidize-go/identity"
)

func TestKeyStore_SaveAccount(t *testing.T) {
	randomIdentity := identity.RandomIdentity()
	keyStore := NewKeyStore(os.TempDir())

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
	keyStore := NewKeyStore(os.TempDir())

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

	assert.Subset(t, identities, []*identity.Identity{randomIdentity1, randomIdentity2})
}
