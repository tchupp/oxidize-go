package main

import (
	"fmt"
	"github.com/tclchiam/block_n_go/blockchain"
	"log"
	"strconv"
	"github.com/tclchiam/block_n_go/tx"
)

func main() {
	const address = "123456"

	bc, err := blockchain.Open("reactions", address)
	if err != nil {
		log.Panic(err)
	}

	transaction := tx.NewCoinbaseTx(address, "Send 4 BTC to Theo")
	transactions := []*tx.Transaction{transaction}

	bc, err = bc.NewBlock(transactions)
	if err != nil {
		log.Panic(err)
	}

	transaction = tx.NewCoinbaseTx(address, "Send 18 BTC to Theo")
	transactions = []*tx.Transaction{transaction}

	bc, err = bc.NewBlock(transactions)
	if err != nil {
		log.Panic(err)
	}

	block := bc.Head()

	for {
		fmt.Printf("============ Block ============\n")
		fmt.Printf("Index: %x\n", block.Index())
		fmt.Printf("Hash: %x\n", block.Hash())
		fmt.Printf("PreviousHash: %x\n", block.PreviousHash())
		fmt.Printf("Transactions:\n")
		for _, transaction := range block.Transactions() {
			fmt.Println(transaction.String())
		}
		fmt.Printf("Nonce: %d\n", block.Nonce())
		fmt.Printf("Is valid: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Println()

		if len(block.PreviousHash()) == 0 {
			break
		}

		block, err = block.Next()
		if err != nil {
			log.Panic(err)
		}
	}

	err = bc.Delete()
	if err != nil {
		log.Panic(err)
	}
}
