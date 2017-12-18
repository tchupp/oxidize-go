package blockchain

import (
	"fmt"

	"github.com/tclchiam/block_n_go/tx"
	"encoding/hex"
	"log"
)

func (bc *Blockchain) buildExpenseTransaction(sender, receiver string, expense int) (*tx.Transaction, error) {
	unspentOutputs, err := bc.findUnspentOutputs(sender)
	if err != nil {
		return nil, err
	}

	balance := calculateBalance(unspentOutputs)
	if balance < expense {
		return nil, fmt.Errorf("account '%s' does not have enough to send: balance '%d', sending '%d'", sender, balance, expense)
	}

	liquidBalance := 0
	takeMinimumToMeetExpense := func(_ string, output *tx.Output) bool {
		take := liquidBalance < expense
		if take {
			liquidBalance += output.Value
		}
		return take
	}

	buildInputs := func(res interface{}, transactionId string, output *tx.Output) interface{} {
		b, err := hex.DecodeString(transactionId)
		if err != nil {
			log.Panic(err)
		}

		input := tx.NewInput(b, output.Id, sender)
		return res.(tx.Inputs).Add(input)
	}

	inputs := unspentOutputs.
		Filter(takeMinimumToMeetExpense).
		Reduce(tx.NewInputs(nil), buildInputs).(tx.Inputs)

	outputs := tx.EmptyOutputs().
		Add(tx.NewOutput(expense, receiver)).
		Add(tx.NewOutput(liquidBalance-expense, sender))

	return tx.NewTx(inputs, outputs), nil
}
