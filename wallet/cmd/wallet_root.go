package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tclchiam/oxidize-go/wallet"
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
	wallet := buildWallet()

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

func buildWallet() wallet.Wallet {
	config := getWalletConfig()

	walletDataDir := config.walletDataDirectory()
	keyStore := wallet.NewKeyStore(walletDataDir)

	return wallet.NewWallet(keyStore)
}
