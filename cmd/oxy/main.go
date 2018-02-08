package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	walletCmd "github.com/tclchiam/oxidize-go/wallet/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "oxy",
		Short: "Oxidize CLI",
		Long:  "CLI for Oxidize protocol",
	}
)

func init() {
	rootCmd.AddCommand(walletCmd.WalletCmd)

	viper.SetDefault("data.dir", "~/.oxy/data")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
