package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("wallet.node.rpc.addr", "localhost:8081")
	viper.SetDefault("wallet.data.dir", "~/.oxy/wallet/data")
}

var cfg = &config{}

type config struct {
}

func getWalletConfig() *config {
	return cfg
}

func (c *config) dataDirectory() string {
	dataDir, _ := homedir.Expand(viper.GetString("wallet.data.dir"))
	return dataDir
}

func (*config) nodeRPCAddress() string {
	return viper.GetString("wallet.node.rpc.addr")
}
