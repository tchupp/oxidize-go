package mining

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func TestFindDifficulty(t *testing.T) {
	type args struct {
		hash *entity.Hash
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "difficulty 0",
			args: args{hash: entity.NewHashOrPanic("69c6ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: 0,
		},
		{
			name: "difficulty 1",
			args: args{hash: entity.NewHashOrPanic("09c6ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: 1,
		},
		{
			name: "difficulty 4",
			args: args{hash: entity.NewHashOrPanic("0000ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: 4,
		},
		{
			name: "difficulty 8",
			args: args{hash: entity.NewHashOrPanic("000000008ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindDifficulty(tt.args.hash); got != tt.want {
				t.Errorf("FindDifficulty(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHasDifficulty(t *testing.T) {
	type args struct {
		hash       *entity.Hash
		difficulty uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has difficulty 0",
			args: args{difficulty: 0, hash: entity.NewHashOrPanic("69c6ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: true,
		},
		{
			name: "has difficulty 1",
			args: args{difficulty: 1, hash: entity.NewHashOrPanic("09c6ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: true,
		},
		{
			name: "has difficulty 4",
			args: args{difficulty: 4, hash: entity.NewHashOrPanic("0000ba688ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: true,
		},
		{
			name: "has difficulty 8",
			args: args{difficulty: 8, hash: entity.NewHashOrPanic("000000008ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: true,
		},
		{
			name: "has difficulty 8, not 6",
			args: args{difficulty: 6, hash: entity.NewHashOrPanic("00009f7a8ee9eb482aad1bf310681da13cd0690bd631beefc959de8ac8440bc8")},
			want: false,
		},
	}

	fmt.Printf("int256: %s\n", hex.EncodeToString(maxUint256.Bytes()))
	fmt.Printf("int255:   %s\n", hex.EncodeToString(new(big.Int).Rsh(maxUint256, 4).Bytes()))
	fmt.Printf("int254:   %s\n", hex.EncodeToString(new(big.Int).Rsh(maxUint256, 8).Bytes()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDifficulty(tt.args.hash, tt.args.difficulty); got != tt.want {
				t.Errorf("HasDifficulty(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
