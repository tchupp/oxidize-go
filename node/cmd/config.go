package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfg = &config{}

func init() {
	viper.SetDefault("node.name", "test_net")
	viper.SetDefault("node.http.port", 8080)
	viper.SetDefault("node.rpc.port", 8081)

	viper.SetDefault("node.data.dir", "~/.oxy/node/data")
}

type config struct {
}

func nodeConfig() *config {
	return cfg
}

func (c *config) dataDirectory() string {
	dataDir, _ := homedir.Expand(viper.GetString("node.data.dir"))
	return dataDir
}

func (*config) nodeName() string {
	return viper.GetString("node.name")
}

func (*config) nodeRPCPort() int {
	return viper.GetInt("node.rpc.port")
}

func (*config) nodeHTTPPort() int {
	return viper.GetInt("node.http.port")
}
