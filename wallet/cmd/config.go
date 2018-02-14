package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("wallet.node.addr", "localhost:8080")
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

func (*config) nodeAddress() string {
	return viper.GetString("wallet.node.addr")
}
