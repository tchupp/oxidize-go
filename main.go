package main

import (
	"fmt"
	"github.com/tclchiam/block_n_go/blockchain"
	"log"
	"strconv"
	"github.com/tclchiam/block_n_go/tx"
)

func main() {
	const address = "Theo"

	bc, err := blockchain.Open("reactions", address)
	if err != nil {
		log.Panic(err)
	}

	transaction := tx.NewCoinbaseTx(address, "Send 4 BTC to Theo")
	transactions := []*tx.Transaction{transaction}

	bc, err = bc.MineBlock(transactions)
	if err != nil {
		log.Panic(err)
	}

	transaction = tx.NewCoinbaseTx(address, "Send 18 BTC to Theo")
	transactions = []*tx.Transaction{transaction}

	bc, err = bc.MineBlock(transactions)
	if err != nil {
		log.Panic(err)
	}

	err = bc.ForEachBlock(printBlock)
	if err != nil {
		log.Panic(err)
	}

	unspentOutputs, err := bc.FindUnspentTransactionOutputs(address)
	if err != nil {
		log.Panic(err)
	}

	balance := 0
	for output := range unspentOutputs {
		balance += output.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)

	err = bc.Delete()
	if err != nil {
		log.Panic(err)
	}
}

func printBlock(block *blockchain.Block) {
	fmt.Printf("============ Block ============\n")
	fmt.Printf("Index: %x\n", block.Index)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("PreviousHash: %x\n", block.PreviousHash)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Is valid: %s\n", strconv.FormatBool(block.Validate()))
	fmt.Printf("Transactions:\n")
	block.ForEachTransaction(func(transaction *tx.Transaction) {
		fmt.Println(transaction.String())
	})
	fmt.Println()
}
