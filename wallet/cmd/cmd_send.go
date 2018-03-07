package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tclchiam/oxidize-go/identity"
)

var (
	sendTransactionReceiver string
	sendTransactionAmount   int64
)

func init() {
	WalletCmd.AddCommand(sendTransactionCommand)

	sendTransactionCommand.Flags().StringVarP(&sendTransactionReceiver, "receiver", "r", "", "Send to this address (required)")
	sendTransactionCommand.MarkFlagRequired("receiver")

	sendTransactionCommand.Flags().Int64VarP(&sendTransactionAmount, "amount", "a", 0, "Amount to send (required)")
	sendTransactionCommand.MarkFlagRequired("amount")
}

var sendTransactionCommand = &cobra.Command{
	Use:   "send",
	Short: "Send to an address",
	Long:  "Send oxygen to an address",
	Run:   runSendTransactionCommand,
}

var runSendTransactionCommand = func(cmd *cobra.Command, args []string) {
	wallet, err := buildWallet()
	if err != nil {
		color.Red("error building wallet: %s\n", err)
		return
	}

	receiverAddress, err := identity.DeserializeAddress(sendTransactionReceiver)
	if err != nil {
		color.Red("invalid receiving address '%s': %s\n", sendTransactionReceiver, err)
		return
	}

	newIdentity, err := wallet.NewIdentity()
	if err != nil {
		color.Red("error creating new identity for change: %s\n", err)
		return
	}

	if err := wallet.Send(receiverAddress, newIdentity.Address(), sendTransactionAmount); err != nil {
		color.Red("%s\n", err)
	}
}
