package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	WalletCmd.AddCommand(checkBalanceCommand)
}

var checkBalanceCommand = &cobra.Command{
	Use:  "balance",
	Long: "Read the balance of the wallet",
	Run:  runReadBalanceCommand,
}

var runReadBalanceCommand = func(cmd *cobra.Command, args []string) {
	wallet := buildWallet()

	balance, err := wallet.Balance()
	if err != nil {
		color.Red("error reading balance: %s\n", err)
		return
	}

	fmt.Print("Balance: ")
	color.White("%d\n", balance)
}
