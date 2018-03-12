package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bytes"

	"github.com/tclchiam/oxidize-go/conv"
	"github.com/tclchiam/oxidize-go/identity"
)

const subsidy = 10
const secretLength = 32

type (
	Transaction struct {
		ID      *Hash
		Inputs  []*SignedInput
		Outputs Outputs
		Secret  []byte
	}

	Transactions []*Transaction

	OutputReference struct {
		ID     *Hash
		Output *Output
	}
)

func (tx *Transaction) IsReward() bool {
	return len(tx.Inputs) == 0
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("-Transaction %s:", tx.ID))
	lines = append(lines, fmt.Sprintf("  Is Reward: %s", strconv.FormatBool(tx.IsReward())))
	lines = append(lines, fmt.Sprintf("  Secret:    %x", tx.Secret))

	lines = append(lines, "  Inputs: ")
	for _, input := range tx.Inputs {
		lines = append(lines, input.string("    "))
	}

	lines = append(lines, "  Outputs: ")
	for _, output := range tx.Outputs {
		lines = append(lines, output.string("    "))
	}

	return strings.Join(lines, "\n")
}

func NewRewardTx(beneficiary *identity.Address) *Transaction {
	var inputs []*SignedInput
	outputs := []*Output{NewOutput(subsidy, beneficiary)}

	return NewTx(inputs, outputs)
}

func NewTx(inputs []*SignedInput, outputs []*Output) *Transaction {
	secret := generateSecret()

	return &Transaction{
		ID:      calculateTransactionId(inputs, outputs, secret),
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

func calculateTransactionId(inputs []*SignedInput, outputs []*Output, secret []byte) *Hash {
	var rawTxContents [][]byte
	for _, input := range inputs {
		rawTxContents = append(rawTxContents, input.OutputReference.ID.Slice())
		rawTxContents = append(rawTxContents, conv.U32ToBytes(input.OutputReference.Output.Index))
		rawTxContents = append(rawTxContents, conv.U64ToBytes(input.OutputReference.Output.Value))
		rawTxContents = append(rawTxContents, input.OutputReference.Output.PublicKeyHash)
		rawTxContents = append(rawTxContents, input.PublicKey.Serialize())
		rawTxContents = append(rawTxContents, input.Signature.Serialize())
	}

	for _, output := range outputs {
		rawTxContents = append(rawTxContents, conv.U32ToBytes(output.Index))
		rawTxContents = append(rawTxContents, conv.U64ToBytes(output.Value))
		rawTxContents = append(rawTxContents, output.PublicKeyHash)
	}

	rawTxContents = append(rawTxContents, secret)

	hash := Hash(sha256.Sum256(bytes.Join(rawTxContents, []byte(nil))))
	return &hash
}

func (txs Transactions) Len() int                            { return len(txs) }
func (txs Transactions) Swap(i, j int)                       { txs[i], txs[j] = txs[j], txs[i] }
func (txs Transactions) Less(i, j int) bool                  { return txs[i].ID.Cmp(txs[j].ID) == -1 }
func (txs Transactions) Add(tx ...*Transaction) Transactions { return append(txs, tx...) }

func (txs Transactions) Sort() Transactions {
	copies := append(Transactions(nil), txs...)
	sort.Sort(copies)
	return copies
}

func (txs Transactions) IsEqual(other Transactions) bool {
	if txs == nil && other == nil {
		return true
	}
	if txs == nil || other == nil {
		return false
	}
	if len(txs) != len(other) {
		return false
	}

	other = other.Sort()
	self := txs.Sort()
	for i := 0; i < self.Len(); i++ {
		if !self[i].ID.IsEqual(other[i].ID) {
			return false
		}
	}

	return true
}

func (txs Transactions) Filter(predicate func(tx *Transaction) bool) Transactions {
	filtered := Transactions{}
	for _, tx := range txs {
		if predicate(tx) {
			filtered = filtered.Add(tx)
		}
	}
	return filtered
}

func (txs Transactions) Reduce(res interface{}, apply func(res interface{}, tx *Transaction) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for _, tx := range txs {
			res = apply(res, tx)
		}
		c <- res
	}()
	return <-c
}
