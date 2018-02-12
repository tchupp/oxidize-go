package cmd

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type config struct {
}

func getWalletConfig() *config {
	return &config{}
}

func (c *config) walletDataDirectory() string {
	return filepath.Join(c.getDataDir(), "wallet")
}

func (*config) nodeAddress() string {
	return viper.GetString("node.addr")
}

func (*config) getDataDir() string {
	dataDir, _ := homedir.Expand(viper.GetString("data.dir"))
	return dataDir
}
