package cmd

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/encoding"
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
	repository := buildRepository()
	nodeWallet := buildWallet()
	beneficiary := getBeneficiary(nodeWallet).Address()
	bc := buildBlockchain(repository, beneficiary)

	handler := interrupt.NewHandler()
	n := buildNode(handler, bc)

	n.Serve()

	handler.WaitForInterrupt()
}

func buildRepository() entity.ChainRepository {
	config := nodeConfig()
	repository := boltdb.Builder(config.nodeName(), encoding.BlockProtoEncoder()).
		WithPath(config.dataDirectory()).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()

	return repository
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

func buildBlockchain(repository entity.ChainRepository, beneficiary *identity.Address) blockchain.Blockchain {
	bc, err := blockchain.Open(repository, proofofwork.NewDefaultMiner(beneficiary))
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
	handler.AddInterruptCallback(n.Shutdown)

	return n
}
