package entity

import (
	"bytes"
	"encoding/hex"
	"testing"
	"fmt"
)

// mainNetGenesisHash is the hash of the first block in the block chain for the
// main network (genesis block).
var mainNetGenesisHash = Hash([HashLength]byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0xd6, 0x68,
	0x9c, 0x08, 0x5a, 0xe1, 0x65, 0x83, 0x1e, 0x93,
	0x4f, 0xf7, 0x63, 0xae, 0x46, 0xa2, 0xa6, 0xc1,
	0x72, 0xb3, 0xf1, 0xb6, 0x0a, 0x8c, 0xe2, 0x6f,
})

func TestHash(t *testing.T) {
	// Hash of block 234439.
	blockHashStr := "014a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef"
	blockHash := NewHashOrPanic(blockHashStr)

	// Hash of block 234440 as byte slice.
	buf := []byte{
		0x79, 0xa6, 0x1a, 0xdb, 0xc6, 0xe5, 0xa2, 0xe1,
		0x39, 0xd2, 0x71, 0x3a, 0x54, 0x6e, 0xc7, 0xc8,
		0x75, 0x63, 0x2e, 0x75, 0xf1, 0xdf, 0x9c, 0x3f,
		0xa6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	hash, err := NewHash(buf)
	if err != nil {
		t.Errorf("NewHash: unexpected error %v", err)
	}

	// Ensure proper size.
	if len(hash) != HashLength {
		t.Errorf("NewHash: hash length mismatch - got: %v, want: %v", len(hash), HashLength)
	}

	// Ensure contents match.
	if !bytes.Equal(hash.Slice(), buf) {
		t.Errorf("NewHash: hash contents mismatch - got: %v, want: %v", hash.Slice(), buf)
	}

	// Ensure contents of hash of block 234440 don't match 234439.
	if hash.IsEqual(blockHash) {
		t.Errorf("IsEqual: hash contents should not match - got: %v, want: %v", hash, blockHash)
	}

	// Set hash from byte slice and ensure contents match.
	hash, err = NewHash(blockHash.Slice())
	if err != nil {
		t.Errorf("setBytes: %v", err)
	}
	if !hash.IsEqual(blockHash) {
		t.Errorf("IsEqual: hash contents mismatch - got: %v, want: %v", hash, blockHash)
	}

	// Ensure nil hashes are handled properly.
	if !(*Hash)(nil).IsEqual(nil) {
		t.Error("IsEqual: nil hashes should match")
	}
	if hash.IsEqual(nil) {
		t.Error("IsEqual: non-nil hash matches nil hash")
	}

	// Invalid size for setBytes.
	hash, err = NewHash([]byte{0x00})
	if err == nil {
		t.Errorf("setBytes: failed to received expected err - got: nil")
	}

	// Invalid size for NewHash.
	invalidHash := make([]byte, HashLength+1)
	_, err = NewHash(invalidHash)
	if err == nil {
		t.Errorf("NewHash: failed to received expected err - got: nil")
	}
}

// TestHashString  tests the stringized output for hashes.
func TestHashString(t *testing.T) {
	// Block 100000 hash.
	wantStr := "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	hash := Hash([HashLength]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xba, 0x27,
		0xaa, 0x20, 0x0b, 0x1c, 0xec, 0xaa, 0xd4, 0x78,
		0xd2, 0xb0, 0x04, 0x32, 0x34, 0x6c, 0x3f, 0x1f,
		0x39, 0x86, 0xda, 0x1a, 0xfd, 0x33, 0xe5, 0x06,
	})

	hashStr := hash.String()
	if hashStr != wantStr {
		t.Errorf("String: wrong hash string - got %v, want %v", hashStr, wantStr)
	}
}

// TestNewHashFromString executes tests against the NewHashFromString function.
func TestNewHashFromString(t *testing.T) {
	tests := []struct {
		in   string
		want Hash
		err  error
	}{
		// Genesis hash.
		{
			"000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			mainNetGenesisHash,
			nil,
		},

		// Genesis hash with stripped leading zeros.
		{
			"19d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			mainNetGenesisHash,
			nil,
		},

		// Empty string.
		{
			"",
			Hash{},
			nil,
		},

		// Single digit hash.
		{
			"1",
			Hash([HashLength]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
			}),
			nil,
		},

		// Block 203707 with stripped leading zeros.
		{
			"3264bc2ac36a60840790ba1d475d01367e7c723da941069e9dc",
			Hash([HashLength]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x26,
				0x4b, 0xc2, 0xac, 0x36, 0xa6, 0x08, 0x40, 0x79,
				0x0b, 0xa1, 0xd4, 0x75, 0xd0, 0x13, 0x67, 0xe7,
				0xc7, 0x23, 0xda, 0x94, 0x10, 0x69, 0xe9, 0xdc,
			}),
			nil,
		},

		// Hash string that is too long.
		{
			"012345678901234567890123456789012345678901234567890123456789012345",
			[HashLength]byte{
				0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45,
				0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01,
				0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67,
				0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23,
			},
			nil,
		},

		// Hash string that is contains non-hex chars.
		{
			"abcdefg",
			Hash{},
			hex.InvalidByteError('g'),
		},
	}

	t.Logf("Running %d tests", len(tests))

	for index, testParams := range tests {
		if err := hashFromStrTestSuite(testParams.in, testParams.want, testParams.err, index); err != nil {
			t.Error(err)
		}
	}
}

func hashFromStrTestSuite(input string, expectedResult Hash, expectedError error, index int) error {
	const unwantedErrStr = "NewHashFromString #%d failed to detect expected error - got: %v want: %v"
	const unexpectedResultStr = "NewHashFromString #%d got: %v want: %v"

	result, actualError := NewHashFromString(input)

	if actualError != expectedError {
		return fmt.Errorf(unwantedErrStr, index, actualError, expectedError)
	} else if actualError != nil {
		return nil
	} else if !expectedResult.IsEqual(result) {
		return fmt.Errorf(unexpectedResultStr, index, result, &expectedResult)
	}

	return nil
}
