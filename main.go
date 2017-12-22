package main

import (
	"fmt"
	"github.com/tclchiam/block_n_go/blockchain"
	"log"
	"strconv"
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/bolt_impl"
	"github.com/tclchiam/block_n_go/wallet"
)

func main() {
	owner := wallet.NewWallet().GetAddress()
	receiver := wallet.NewWallet().GetAddress()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n", owner, receiver)

	repository, err := bolt_impl.NewRepository(blockchainName)
	if err != nil {
		log.Panic(err)
	}
	defer repository.Close()
	defer bolt_impl.DeleteBlockchain(blockchainName)

	bc, err := blockchain.Open(repository, owner)
	if err != nil {
		log.Panic(err)
	}

	bc, err = bc.Send(owner, receiver, 3)
	if err != nil {
		log.Panic(err)
	}

	err = bc.ForEachBlock(printBlock)
	if err != nil {
		log.Panic(err)
	}

	balance, err := bc.ReadBalance(owner)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Balance of '%s': %d\n\n", owner, balance)

	balance, err = bc.ReadBalance(receiver)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Balance of '%s': %d\n\n", receiver, balance)
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
