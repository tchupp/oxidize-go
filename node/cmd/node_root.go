package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
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
	nodeName := "test_node"
	repository := boltdb.Builder(nodeName, encoding.BlockProtoEncoder()).
		WithCache().
		WithMetrics().
		WithLogger().
		Build()

	nodeWallet := buildWallet()
	beneficiary := getBeneficiary(nodeWallet)

	bc, err := blockchain.Open(repository, proofofwork.NewDefaultMiner(beneficiary.Address()))
	if err != nil {
		log.WithError(err).Panic("failed to open blockchain")
	}
	defer boltdb.DeleteBlockchain(nodeName)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.WithError(err).Panic("failed to listen")
	}

	server := rpc.NewServer(lis)
	baseNode := node.WrapWithLogger(node.NewNode(bc, server))

	baseNode.Serve()

	interruptReceived := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			interruptReceived <- true
		}
	}()
	<-interruptReceived

	baseNode.Shutdown()
}

func getBeneficiary(nodeWallet wallet.Wallet) *identity.Identity {
	identities, err := nodeWallet.Identities()
	if err != nil {
		log.WithError(err).Panic("error reading identities")
	}
	if len(identities) != 0 {
		return identities[0]
	}
	beneficiary, err := nodeWallet.NewIdentity()
	if err != nil {
		log.WithError(err).Panic("error creating new identity")
	}
	return beneficiary

}

func buildWallet() wallet.Wallet {
	config := getNodeConfig()
	nodeDataDir := config.nodeDataDirectory()

	keyStore := wallet.NewKeyStore(nodeDataDir)
	return wallet.NewWallet(keyStore, nil)
}
