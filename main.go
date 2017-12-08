package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain().
		AddBlock("Send 4 BTC to Theo").
		AddBlock("Send 18 BTC to Theo")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Is valid: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Println()
	}
}
