package main

import (
	"fmt"
	"log"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/storage/boltdb"
	"github.com/tclchiam/block_n_go/wallet"
)

func main() {
	owner := wallet.NewWallet()
	receiver := wallet.NewWallet()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n\n", owner.GetAddress(), receiver.GetAddress())

	blockRepository, err := boltdb.NewBlockRepository(blockchainName, encoding.NewBlockGobEncoder())
	if err != nil {
		log.Panic(err)
	}
	defer blockRepository.Close()
	defer boltdb.DeleteBlockchain(blockchainName)

	bc, err := blockchain.Open(blockRepository, proofofwork.NewDefaultMiner(), owner.GetAddress())
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

	err = bc.ForEachBlock(func(block *entity.Block) {
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
