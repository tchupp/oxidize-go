package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/storage/boltdb"
)

func main() {
	owner := identity.RandomIdentity()
	receiver := identity.RandomIdentity()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n", owner, receiver)

	miner := proofofwork.NewDefaultMiner(owner)

	repository := boltdb.Builder(blockchainName, encoding.BlockProtoEncoder()).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()
	defer repository.Close()
	defer boltdb.DeleteBlockchain(blockchainName)

	genesisBlock := miner.MineBlock(&entity.GenesisParentHeader, entity.Transactions{})
	if err := repository.SaveBlock(genesisBlock); err != nil {
		log.Panic(err)
	}

	bc, err := blockchain.Open(repository, miner)
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
