package block

import (
	"fmt"
	"time"
	"strings"

	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
)

type Header struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*tx.Transaction
}

func NewGenesisBlockHeader(address string) *Header {
	transaction := tx.NewGenesisCoinbaseTx(address)

	return NewBlockHeader(0, chainhash.EmptyHash, []*tx.Transaction{transaction})
}

func NewBlockHeader(index int, previousHash chainhash.Hash, transactions []*tx.Transaction) *Header {
	return &Header{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
}

func (header *Header) ForEachTransaction(consume func(*tx.Transaction)) {
	for _, transaction := range header.Transactions {
		consume(transaction)
	}
}

func (header Header) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ Header ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", header.Index))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", header.PreviousHash.Slice()))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", header.Timestamp))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	header.ForEachTransaction(func(transaction *tx.Transaction) {
		lines = append(lines, transaction.String())
	})

	return strings.Join(lines, "\n")
}
