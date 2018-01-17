package consensus

import (
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func TestVerifyHeader(t *testing.T) {
	err := VerifyHeader(entity.DefaultGenesisBlock().Header())
	if err != nil {
		t.Errorf("verifying header: %s", err)
	}
}
