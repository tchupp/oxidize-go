package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	nodeCmd "github.com/tclchiam/oxidize-go/node/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "oxyctrl",
		Short: "Oxidize control CLI",
		Long:  "CLI for Oxidize protocol",
	}
)

func init() {
	rootCmd.AddCommand(nodeCmd.NodeCmd)

	viper.SetDefault("data.dir", "~/.oxy/data")
	viper.SetDefault("node.addr", "localhost:8080")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}
}
