package tx

import (
	"testing"
	"github.com/tclchiam/block_n_go/tx/txset"
)

func TestTransaction_FindUnspentOutput(t *testing.T) {
	const to = "Theo"

	t.Run("One", func(t *testing.T) {
		transaction := NewCoinbaseTx(to, "")
		spentTransactions := txset.New()

		unspentOutput := transaction.FindUnspentOutput(spentTransactions, to)

		if len(unspentOutput) != 1 {
			t.Fatalf("Expected %d unspet output, was: %d", 1, len(unspentOutput))
		}
	})
}
