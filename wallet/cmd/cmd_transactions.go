package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	WalletCmd.AddCommand(viewTransactionsCommand)
}

var viewTransactionsCommand = &cobra.Command{
	Use:     "transactions",
	Aliases: []string{"txs"},
	Short:   "View wallet transactions",
	Long:    "View the transactions of the wallet",
	Run:     runViewTransactionsCommand,
}

var runViewTransactionsCommand = func(cmd *cobra.Command, args []string) {
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
		color.Cyan("%s\n", account.Address())

		for i, tx := range account.Transactions() {
			if tx.Spender().IsEqual(account.Address()) {
				color.Red("%d| To: '%s',   Amount: %d\n", i, tx.Receiver(), tx.Amount())
			}

			if tx.Receiver().IsEqual(account.Address()) {
				if tx.Spender() == nil {
					color.Green("%d| Reward - Amount: %d\n", i, tx.Amount())
				} else {
					color.Green("%d| From: '%s' - Amount: %d\n", i, tx.Spender(), tx.Amount())
				}
			}
		}
	}
}
