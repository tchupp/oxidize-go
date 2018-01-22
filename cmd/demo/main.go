package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/storage/boltdb"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/mining"
)

func main() {
	owner := identity.RandomIdentity()
	receiver := identity.RandomIdentity()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n", owner, receiver)

	miner := proofofwork.NewDefaultMiner(owner)
	blockRepository, err := boltdb.NewBlockRepository(blockchainName, encoding.BlockProtoEncoder())
	if err != nil {
		log.Panic(err)
	}
	defer boltdb.DeleteBlockchain(blockchainName)

	genesisBlock := buildGenesisBlock(miner)
	if err = blockRepository.SaveBlock(genesisBlock); err != nil {
		log.Panic(err)
	}

	bc, err := blockchain.Open(blockRepository, miner)
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

func buildGenesisBlock(miner mining.Miner) *entity.Block {
	header := entity.NewBlockHeader(math.MaxUint64, nil, nil, 0, 0, &entity.EmptyHash)
	parent := entity.NewBlock(header, entity.Transactions{})

	return miner.MineBlock(parent, entity.Transactions{})
}
