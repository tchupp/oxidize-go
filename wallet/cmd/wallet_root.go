package cmd

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tclchiam/oxidize-go/wallet"
	"github.com/tclchiam/oxidize-go/wallet/rpc"
)

var (
	WalletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Oxidize Wallet Sub-command",
		Long:  "Wallet CLI for Oxidize protocol",
		Run:   showWalletSummaryCommand,
	}
)

var showWalletSummaryCommand = func(cmd *cobra.Command, args []string) {
	wallet, err := buildWallet()
	if err != nil {
		color.Red("error building wallet: %s\n", err)
		return
	}

	identities, err := wallet.Identities()
	if err != nil {
		color.Red("error reading addresses: %s\n", err)
		return
	}

	color.White("Wallet:")
	for i, identity := range identities {
		color.Cyan("%d. %s\n", i+1, identity.Address())
	}
	return
}

func buildWallet() (wallet.Wallet, error) {
	config := getWalletConfig()

	walletDataDir := config.dataDirectory()
	keyStore := wallet.NewKeyStore(walletDataDir)

	conn, err := grpc.Dial(config.nodeAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("node '%s' is not up", config.nodeAddress())
	}
	client := rpc.NewWalletClient(conn)

	return wallet.NewWallet(keyStore, client), nil
}
