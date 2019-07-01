package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jun06t/go-sample/go-plugin-rpc/common"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/auth"),
		Logger:          logger,
		SyncStdout:      os.Stdout,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := rpcClient.Dispense(common.AuthPluginName)
	if err != nil {
		log.Fatal(err)
	}

	auth := raw.(common.Auth)
	if auth.Authenticate() {
		fmt.Println("success!!")
	} else {
		fmt.Println("fail!!")
	}
}

var pluginMap = map[string]plugin.Plugin{
	common.AuthPluginName: &common.AuthPlugin{},
}
