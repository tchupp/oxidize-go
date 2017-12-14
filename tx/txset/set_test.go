package txset

import "testing"

func TestTransactionSet(t *testing.T) {
	transactionSet := New()

	transactionSet.Add("tx1", 0)
	if transactionSet.Contains("tx1", 0) == true {
		t.Fatalf("Expected transaction set to not contain, tx Id: %s, with output: %d ", "tx1", 0)
	}

	transactionSet = transactionSet.Add("tx2", 2)
	if transactionSet.Contains("tx2", 2) == false {
		t.Fatalf("Expected transaction set to not contain, tx Id: %s, with output: %d ", "tx2", 2)
	}
}
