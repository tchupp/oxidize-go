package cmd

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfg = &config{}

func init() {
	viper.SetDefault("node.name", "test_net")
	viper.SetDefault("node.port", 8080)
}

type config struct {
}

func nodeConfig() *config {
	return cfg
}

func (c *config) dataDirectory() string {
	return filepath.Join(c.getDataDir(), "node")
}

func (*config) nodeName() string {
	return viper.GetString("node.name")
}

func (*config) nodePort() int {
	return viper.GetInt("node.port")
}

func (*config) getDataDir() string {
	dataDir, _ := homedir.Expand(viper.GetString("data.dir"))
	return dataDir
}
