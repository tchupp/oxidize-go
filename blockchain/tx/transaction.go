package tx

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

const subsidy = 10
const secretLength = 32

type (
	Transaction struct {
		ID        TransactionId
		TxInputs  []*SignedInput
		TxOutputs []*Output
		Secret    []byte
	}

	OutputReference struct {
		ID     TransactionId
		Output *Output
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 0
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %s:", tx.ID))
	lines = append(lines, fmt.Sprintf("     Is Coinbase: %s", strconv.FormatBool(tx.IsCoinbase())))

	for _, input := range tx.TxInputs {
		lines = append(lines, input.String())
	}

	for _, output := range tx.TxOutputs {
		lines = append(lines, output.String())
	}

	return strings.Join(lines, "\n")
}

func NewGenesisCoinbaseTx(ownerAddress string) *Transaction {
	return NewCoinbaseTx(ownerAddress)
}

func NewCoinbaseTx(minerAddress string) *Transaction {
	var inputs []*SignedInput
	outputs := []*Output{NewOutput(subsidy, minerAddress)}

	return NewTx(inputs, outputs)
}

func NewTx(inputs []*SignedInput, outputs []*Output) *Transaction {
	secret := generateSecret()

	return &Transaction{
		ID:        calculateTransactionId(inputs, outputs, secret),
		TxInputs:  inputs,
		TxOutputs: outputs,
		Secret:    secret,
	}
}

func generateSecret() []byte {
	secret := make([]byte, secretLength)
	rand.Read(secret)
	return secret
}
