package tx

import (
	"fmt"
	"crypto/sha256"
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

const subsidy = 10
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Transaction struct {
	ID   []byte
	Vin  []Input
	Vout []Output
}

func NewGenesisCoinbaseTx(to string) *Transaction {
	return NewCoinbaseTx(to, genesisCoinbaseData)
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := Input{Txid: []byte(nil), Vout: -1, ScriptSig: data}
	txout := Output{subsidy, to}
	tx := Transaction{nil, []Input{txin}, []Output{txout}}
	tx.ID = tx.Hash()

	return &tx
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		//lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.ScriptSig))
		//lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		//lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}