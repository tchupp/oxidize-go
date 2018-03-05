package cmd

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/rpc"
	"github.com/tclchiam/oxidize-go/storage/boltdb"
	"github.com/tclchiam/oxidize-go/wallet"
)

var (
	NodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Oxidize Node Sub-command",
		Long:  "Node CLI for Oxidize protocol",
		Run:   showNodeSummaryCommand,
	}
)

var showNodeSummaryCommand = func(cmd *cobra.Command, args []string) {
	chainRepository := buildChainRepository()
	utxoRepository := buildUtxoRepository()
	nodeWallet := buildWallet()
	beneficiary := getBeneficiary(nodeWallet).Address()
	bc := buildBlockchain(chainRepository, utxoRepository, beneficiary)

	handler := interrupt.NewHandler()
	n := buildNode(handler, bc)

	n.Serve()

	handler.WaitForInterrupt()
}

func buildChainRepository() entity.ChainRepository {
	config := nodeConfig()

	return boltdb.ChainBuilder(config.nodeName()).
		WithPath(config.dataDirectory()).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()
}

func buildUtxoRepository() utxo.Repository {
	config := nodeConfig()

	return boltdb.UtxoBuilder(config.nodeName()).
		WithPath(config.dataDirectory()).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()
}

func buildWallet() wallet.Wallet {
	nodeDataDir := nodeConfig().dataDirectory()

	keyStore := wallet.NewKeyStore(nodeDataDir)
	return wallet.NewWallet(keyStore, nil)
}

func getBeneficiary(nodeWallet wallet.Wallet) *identity.Identity {
	identities, err := nodeWallet.Identities()
	if err != nil {
		log.WithError(err).Panic("error reading identities")
	}

	if id := identities.FirstOrNil(); id != nil {
		return id
	}

	beneficiary, err := nodeWallet.NewIdentity()
	if err != nil {
		log.WithError(err).Panic("error creating new identity")
	}
	return beneficiary
}

func buildBlockchain(chainRepository entity.ChainRepository, utxoRepository utxo.Repository, beneficiary *identity.Address) blockchain.Blockchain {
	bc, err := blockchain.Open(chainRepository, utxoRepository, proofofwork.NewDefaultMiner(beneficiary))
	if err != nil {
		log.WithError(err).Panic("failed to open blockchain")
	}

	return bc
}

func buildNode(handler interrupt.Handler, bc blockchain.Blockchain) node.Node {
	config := nodeConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", config.nodePort()))
	if err != nil {
		log.WithError(err).Panic("failed to listen")
	}

	n := node.WrapWithLogger(node.NewNode(bc, rpc.NewServer(lis)))
	handler.AddInterruptCallback(func() { n.Close() })

	return n
}
