package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/iter"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/storage/boltdb"
)

func main() {
	owner := identity.RandomIdentity()
	receiver := identity.RandomIdentity()
	const blockchainName = "reactions"

	fmt.Printf("Owner: '%s', receiver: '%s'\n", owner, receiver)

	miner := proofofwork.NewDefaultMiner(owner.Address())

	chainRepository := boltdb.ChainBuilder(blockchainName).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()
	defer chainRepository.Close()
	defer boltdb.DeleteBlockchain(blockchainName)

	utxoRepository := boltdb.UtxoBuilder(blockchainName).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()
	defer utxoRepository.Close()
	defer boltdb.DeleteUtxo(blockchainName)

	genesisBlock := miner.MineBlock(&entity.GenesisParentHeader, entity.Transactions{})
	if err := chainRepository.SaveBlock(genesisBlock); err != nil {
		log.Panic(err)
	}

	bc, err := blockchain.Open(chainRepository, utxoRepository, miner)
	if err != nil {
		log.Panic(err)
	}
	accountEngine := account.NewEngine(bc)

	err = accountEngine.Send(owner, receiver.Address(), 7)
	if err != nil {
		log.Panic(err)
	}
	err = accountEngine.Send(receiver, owner.Address(), 4)
	if err != nil {
		log.Panic(err)
	}

	err = iter.ForEachBlock(chainRepository, func(block *entity.Block) {
		fmt.Println(block)
	})
	if err != nil {
		log.Panic(err)
	}

	ownerAccount, err := accountEngine.Account(owner.Address())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Account: %s\n\n", ownerAccount)

	receiverAccount, err := accountEngine.Account(receiver.Address())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Account: %s\n\n", receiverAccount)
}
