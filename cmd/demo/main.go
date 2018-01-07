package main

import (
	"fmt"
	"log"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/block"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/storage/boltdb"
	"github.com/tclchiam/block_n_go/wallet"
)

func main() {
	owner := wallet.NewWallet()
	receiver := wallet.NewWallet()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n\n", owner.GetAddress(), receiver.GetAddress())

	repository, err := boltdb.NewRepository(blockchainName)
	if err != nil {
		log.Panic(err)
	}
	defer repository.Close()
	defer boltdb.DeleteBlockchain(blockchainName)

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

	err = bc.ForEachBlock(func(block *block.Block) {
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