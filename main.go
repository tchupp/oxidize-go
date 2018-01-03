package main

import (
	"fmt"
	"github.com/tclchiam/block_n_go/blockchain"
	"log"
	"github.com/tclchiam/block_n_go/bolt_impl"
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/proofofwork"
)

func main() {
	owner := wallet.NewWallet()
	receiver := wallet.NewWallet()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n\n", owner.GetAddress(), receiver.GetAddress())

	repository, err := bolt_impl.NewRepository(blockchainName)
	if err != nil {
		log.Panic(err)
	}
	defer repository.Close()
	defer bolt_impl.DeleteBlockchain(blockchainName)

	bc, err := blockchain.Open(repository, proofofwork.NewDefaultMiner(), owner.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	err = bc.Send(owner, receiver, owner, 7)
	if err != nil {
		log.Panic(err)
	}
	err = bc.Send(receiver, owner, owner, 4)
	if err != nil {
		log.Panic(err)
	}

	err = bc.ForEachBlock(func(block *blockchain.Block) {
		fmt.Println(block)
	})
	if err != nil {
		log.Panic(err)
	}

	balance, err := bc.ReadBalance(owner.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Balance of '%s': %d\n\n", owner.GetAddress(), balance)

	balance, err = bc.ReadBalance(receiver.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Balance of '%s': %d\n\n", receiver.GetAddress(), balance)
}
