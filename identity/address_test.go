package identity

import (
	"reflect"
	"testing"

	"github.com/tclchiam/oxidize-go/crypto"
)

func TestDeserializeAddress(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *Address
		wantErr bool
	}{
		{
			name: "parses correctly",
			data: "12GJBjXZr8DvjYwBgeXWjP4pSkhfyXUXT7",
			want: &Address{
				version:       byte(0x00),
				publicKeyHash: []byte{13, 220, 181, 232, 21, 237, 117, 95, 95, 103, 189, 221, 213, 173, 42, 66, 22, 138, 41, 112},
				checksum:      []byte{174, 246, 175, 186},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeserializeAddress(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("version: %b,", got.version)
				t.Logf("publicKeyHash: %v,", got.publicKeyHash)
				t.Logf("checksum: %v\n", got.checksum)
				t.Errorf("DeserializeAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddress_Serialize(t *testing.T) {
	identity := NewIdentity(crypto.NewP256PrivateKey())
	serializedAddress := identity.Address().Serialize()

	address, err := DeserializeAddress(serializedAddress)
	if err != nil {
		t.Fatalf("failed to deserialize address: %s", err)
	}

	if !identity.Address().IsEqual(address) {
		t.Errorf("expected addresses to be equal. wanted %s, got %s", identity.Address(), address)
	}

}
