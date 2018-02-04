package engine_test

import (
	"testing"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/engine"
	"github.com/tclchiam/block_n_go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

var (
	miner = proofofwork.NewDefaultMiner(identity.RandomIdentity())

	header1 = miner.MineBlock(&entity.GenesisParentHeader, entity.Transactions{}).Header()
	header2 = miner.MineBlock(header1, entity.Transactions{}).Header()
	header3 = miner.MineBlock(header2, entity.Transactions{}).Header()

	badHeader = &entity.BlockHeader{
		Index:            3,
		PreviousHash:     entity.NewHashOrPanic("0000d7ed4c5c6cd34828d07d43da441eab32dca9abc352298fdbb0f8d887ee2e"),
		Timestamp:        1517243842,
		TransactionsHash: entity.NewHashOrPanic("d8c4558738d6cf0d5ef3069bd335888b0ca9c7391e6dd7a07298743e8f3b7759"),
		Nonce:            22243,
		Hash:             entity.NewHashOrPanic("0000701b473d9ea60442bd4af2d77efa64657be3ca0932e02d192ac7ceddb5cd"),
	}
)

func TestSaveHeaders(t *testing.T) {
	type args struct {
		headers entity.BlockHeaders
		bc      blockchain.Blockchain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "none",
			args:    args{bc: buildBlockchain(t), headers: entity.BlockHeaders{}},
			wantErr: false,
		},
		{
			name:    "one",
			args:    args{bc: buildBlockchain(t), headers: entity.BlockHeaders{header1}},
			wantErr: false,
		},
		{
			name:    "many",
			args:    args{bc: buildBlockchain(t), headers: entity.BlockHeaders{header1, header2, header3}},
			wantErr: false,
		},
		{
			name:    "out of order",
			args:    args{bc: buildBlockchain(t), headers: entity.BlockHeaders{header3, header1, header2}},
			wantErr: false,
		},
		{
			name:    "bad header",
			args:    args{bc: buildBlockchain(t), headers: entity.BlockHeaders{badHeader}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("SaveHeaders() name = %s", tt.name)
			if err := engine.SaveHeaders(tt.args.headers, tt.args.bc); (err != nil) != tt.wantErr {
				t.Errorf("SaveHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func buildBlockchain(t *testing.T) blockchain.Blockchain {
	bc, err := blockchain.Open(memdb.NewChainRepository(), nil)
	if err != nil {
		t.Fatalf("failed to open test blockchain: %s", err)
	}
	return bc
}
