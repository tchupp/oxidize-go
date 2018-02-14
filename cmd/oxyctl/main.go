package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	nodeCmd "github.com/tclchiam/oxidize-go/node/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "oxyctl",
		Short: "Oxidize control CLI",
		Long:  "CLI for Oxidize protocol",
	}
)

func init() {
	rootCmd.AddCommand(nodeCmd.NodeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}
}
