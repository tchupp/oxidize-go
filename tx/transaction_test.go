package tx

import (
	"testing"

	"github.com/tclchiam/block_n_go/wallet"
)

func TestTransaction_FindUnspentOutput(t *testing.T) {
	address := wallet.NewWallet().GetAddress()

	t.Run("One", func(t *testing.T) {
		transaction := NewCoinbaseTx(address)

		unspentOutputs := transaction.FindOutputsForAddress(address)
		if unspentOutputs.Len() != 1 {
			t.Fatalf("Expected %d unspent output, was: %d", 1, unspentOutputs.Len())
		}
	})
}

func TestTransaction_FindSpentOutput(t *testing.T) {
	address := wallet.NewWallet().GetAddress()

	t.Run("One", func(t *testing.T) {
		transaction := NewCoinbaseTx(address)

		spentOutputs := transaction.FindSpentOutputs(address)

		if len(spentOutputs) != 0 {
			t.Fatalf("Expected %d spent output, was: %d", 0, len(spentOutputs))
		}
	})
}
