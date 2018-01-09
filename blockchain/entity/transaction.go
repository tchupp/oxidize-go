package entity

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

const subsidy = 10
const secretLength = 32

type (
	Transaction struct {
		ID      TransactionId
		Inputs  []*SignedInput
		Outputs []*Output
		Secret  []byte
	}

	OutputReference struct {
		ID     TransactionId
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

func NewGenesisCoinbaseTx(ownerAddress string, encoder TransactionEncoder) *Transaction {
	return NewCoinbaseTx(ownerAddress, encoder)
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

type TransactionId [sha256.Size]byte

func (txId TransactionId) String() string {
	return hex.EncodeToString(txId[:])
}

func calculateTransactionId(inputs []*SignedInput, outputs []*Output, secret []byte, encoder TransactionEncoder) TransactionId {
	return sha256.Sum256(serializeTxData(inputs, outputs, secret, encoder))
}

func serializeTxData(inputs []*SignedInput, outputs []*Output, secret []byte, encoder TransactionEncoder) []byte {
	transaction := &Transaction{
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
