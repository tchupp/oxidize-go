package main

import (
	"github.com/tclchiam/block_n_go/blockchain"
	"log"
	"fmt"
	"strconv"
)

func main() {
	bc, err := blockchain.Open("reactions")
	if err != nil {
		log.Panic(err)
	}

	bc, err = bc.NewBlock("Send 4 BTC to Theo")
	if err != nil {
		log.Panic(err)
	}

	bc, err = bc.NewBlock("Send 18 BTC to Theo")
	if err != nil {
		log.Panic(err)
	}

	block := bc.Head()

	for block != nil {
		fmt.Printf("Index: %x\n", block.Index)
		fmt.Printf("PreviousHash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Is valid: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Println()

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
