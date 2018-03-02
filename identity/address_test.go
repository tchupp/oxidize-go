package identity_test

import (
	"reflect"
	"testing"

	"bytes"
	"crypto/sha256"

	"github.com/tclchiam/oxidize-go/crypto"
	"github.com/tclchiam/oxidize-go/identity"
	"golang.org/x/crypto/ripemd160"
)

func TestDeserializeAddress(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *identity.Address
		wantErr bool
	}{
		{
			name: "parses valid correctly",
			data: "12GJBjXZr8DvjYwBgeXWjP4pSkhfyXUXT7",
			want: identity.NewAddress(
				byte(0x00),
				[]byte{13, 220, 181, 232, 21, 237, 117, 95, 95, 103, 189, 221, 213, 173, 42, 66, 22, 138, 41, 112},
				[]byte{174, 246, 175, 186},
			),
			wantErr: false,
		},
		{
			name:    "parses empty correctly",
			data:    "",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := identity.DeserializeAddress(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("version: %b,", got.Version())
				t.Logf("publicKeyHash: %v,", got.PublicKeyHash())
				t.Logf("checksum: %v\n", got.Checksum())
				t.Errorf("DeserializeAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddress_Serialize(t *testing.T) {
	address := identity.RandomIdentity().Address()
	serializedAddress := address.Serialize()

	deserializedAddress, err := identity.DeserializeAddress(serializedAddress)
	if err != nil {
		t.Fatalf("failed to deserialize address: %s", err)
	}

	if !address.IsEqual(deserializedAddress) {
		t.Errorf("expected addresses to be equal. wanted %s, got %s", address, deserializedAddress)
	}
}

func TestAddress_Version(t *testing.T) {
	address := identity.RandomIdentity().Address()

	expectedVersion := byte(0x00)

	if expectedVersion != address.Version() {
		t.Errorf("Expected version did not equal actual. Got: '%b', wanted: '%b'", address.Version(), expectedVersion)
	}
}

func TestAddress_PublicKeyHash(t *testing.T) {
	privateKey := crypto.NewP256PrivateKey()
	publicKey := privateKey.PubKey()
	address := identity.NewIdentity(privateKey).Address()

	publicSHA256 := sha256.Sum256(publicKey.Serialize())

	hashImpl := ripemd160.New()
	hashImpl.Write(publicSHA256[:])
	expectedHash := hashImpl.Sum(nil)

	if len(expectedHash) != 20 {
		t.Errorf("Expected len did not equal actual. Got: %d, wanted: %d", len(expectedHash), 20)
	}
	if bytes.Compare(expectedHash, address.PublicKeyHash()) != 0 {
		t.Errorf("Expected hash did not equal actual. Got: '%s', wanted: '%s'", address, expectedHash)
	}
}

func TestAddress_Checksum(t *testing.T) {
	address := identity.RandomIdentity().Address()

	payload := append([]byte{address.Version()}, address.PublicKeyHash()...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	expectedChecksum := secondSHA[:4]

	if len(expectedChecksum) != 4 {
		t.Errorf("Expected len did not equal actual. Got: %d, wanted: %d", len(expectedChecksum), 4)
	}
	if bytes.Compare(expectedChecksum, address.Checksum()) != 0 {
		t.Errorf("Expected checksum did not equal actual. Got: '%s', wanted: '%s'", address.Checksum(), expectedChecksum)
	}
}
