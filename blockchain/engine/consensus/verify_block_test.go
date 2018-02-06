package consensus

import (
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func TestVerifyBlock(t *testing.T) {
	err := VerifyBlock(entity.DefaultGenesisBlock())

	if err != nil {
		t.Errorf("verifying block: %s", err)
	}
}
