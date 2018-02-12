package cmd

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type config struct {
}

func getNodeConfig() *config {
	return &config{}
}

func (c *config) nodeDataDirectory() string {
	return filepath.Join(c.getDataDir(), "node")
}

func (*config) getDataDir() string {
	dataDir, _ := homedir.Expand(viper.GetString("data.dir"))
	return dataDir
}
