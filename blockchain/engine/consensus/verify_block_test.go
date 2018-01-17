package consensus

import (
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func TestVerifyBlock(t *testing.T) {
	err := VerifyBlock(entity.DefaultGenesisBlock())

	if err != nil {
		t.Errorf("verifying block: %s", err)
	}
}
