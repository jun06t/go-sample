package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jun06t/go-sample/go-plugin-rpc/common"
)

const (
	password = "hoge"
)

type PasswordAuth struct {
	logger hclog.Logger
}

func (p *PasswordAuth) Authenticate() bool {
	return true
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	auth := &PasswordAuth{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		common.AuthPluginName: &common.AuthPlugin{Impl: auth},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins:         pluginMap,
	})
}
