package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	walletCmd "github.com/tclchiam/oxidize-go/wallet/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "oxy",
		Short: "Oxidize client CLI",
		Long:  "CLI for Oxidize protocol",
	}
)

func init() {
	rootCmd.AddCommand(walletCmd.WalletCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}
}
