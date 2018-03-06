package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	WalletCmd.AddCommand(checkBalanceCommand)
}

var checkBalanceCommand = &cobra.Command{
	Use:     "balance",
	Aliases: []string{"b"},
	Short:   "View wallet balance",
	Long:    "View the balance of the wallet",
	Run:     runReadBalanceCommand,
}

var runReadBalanceCommand = func(cmd *cobra.Command, args []string) {
	wallet, err := buildWallet()
	if err != nil {
		color.Red("error building wallet: %s\n", err)
		return
	}

	accounts, err := wallet.Accounts()
	if err != nil {
		color.Red("error reading balance: %s\n", err)
		return
	}

	color.White("Account: \n")
	for _, account := range accounts {
		color.Cyan("%s: %d\n", account.Address(), account.Spendable())
	}
}
