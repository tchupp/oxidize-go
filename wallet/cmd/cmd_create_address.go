package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	WalletCmd.AddCommand(createAddressCommand)
}

var createAddressCommand = &cobra.Command{
	Use:   "new",
	Short: "Create new address",
	Long:  "Create a new address for the wallet",
	Run:   runCreateAddressCommand,
}

var runCreateAddressCommand = func(cmd *cobra.Command, args []string) {
	wallet, err := buildWallet()
	if err != nil {
		color.Red("error building wallet: %s\n", err)
		return
	}

	newIdentity, err := wallet.NewIdentity()
	if err != nil {
		color.Red("error saving new address: %s\n", err)
		return
	}

	fmt.Print("Saved new address: ")
	color.White("%s\n", newIdentity)
}
