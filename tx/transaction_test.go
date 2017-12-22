package tx

import (
	"testing"
)

func TestTransaction_FindUnspentOutput(t *testing.T) {
	const to = "Theo"

	t.Run("One", func(t *testing.T) {
		transaction := NewCoinbaseTx(to)

		unspentOutputs := transaction.FindOutputsForAddress(to)
		count := unspentOutputs.Reduce(0, func(res interface{}, transactionId string, output *Output) interface{} {
			return res.(int) + 1
		})

		if count != 1 {
			t.Fatalf("Expected %d unspent output, was: %d", 1, count)
		}
	})
}

func TestTransaction_FindSpentOutput(t *testing.T) {
	const to = "Theo"

	t.Run("One", func(t *testing.T) {
		transaction := NewCoinbaseTx(to)

		spentOutputs := transaction.FindSpentOutputs(to)

		if len(spentOutputs) != 0 {
			t.Fatalf("Expected %d spent output, was: %d", 0, len(spentOutputs))
		}
	})
}
