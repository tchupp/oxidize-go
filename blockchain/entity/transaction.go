package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const subsidy = 10
const secretLength = 32

type (
	Transaction struct {
		ID      *Hash
		Inputs  []*SignedInput
		Outputs []*Output
		Secret  []byte
	}

	Transactions []*Transaction

	OutputReference struct {
		ID     *Hash
		Output *Output
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 0
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %s:", tx.ID))
	lines = append(lines, fmt.Sprintf("     Is Coinbase: %s", strconv.FormatBool(tx.IsCoinbase())))
	lines = append(lines, fmt.Sprintf("     Secret:      %x", tx.Secret))

	for _, input := range tx.Inputs {
		lines = append(lines, input.String())
	}

	for _, output := range tx.Outputs {
		lines = append(lines, output.String())
	}

	return strings.Join(lines, "\n")
}

func NewCoinbaseTx(coinbaseAddress string, encoder TransactionEncoder) *Transaction {
	var inputs []*SignedInput
	outputs := []*Output{NewOutput(subsidy, coinbaseAddress)}

	return NewTx(inputs, outputs, encoder)
}

func NewTx(inputs []*SignedInput, outputs []*Output, encoder TransactionEncoder) *Transaction {
	secret := generateSecret()

	return &Transaction{
		ID:      calculateTransactionId(inputs, outputs, secret, encoder),
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  secret,
	}
}

func generateSecret() []byte {
	secret := make([]byte, secretLength)
	rand.Read(secret)
	return secret
}

func calculateTransactionId(inputs []*SignedInput, outputs []*Output, secret []byte, encoder TransactionEncoder) *Hash {
	hash := Hash(sha256.Sum256(serializeTxData(inputs, outputs, secret, encoder)))
	return &hash
}

func serializeTxData(inputs []*SignedInput, outputs []*Output, secret []byte, encoder TransactionEncoder) []byte {
	transaction := &Transaction{
		ID:      &EmptyHash,
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  secret,
	}

	encoded, err := encoder.EncodeTransaction(transaction)
	if err != nil {
		log.Panic(err)
	}
	return encoded
}
