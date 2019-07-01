package common

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type Auth interface {
	Authenticate() bool
}

type AuthRPC struct {
	client *rpc.Client
}

func (a *AuthRPC) Authenticate() bool {
	var resp bool
	err := a.client.Call("Plugin.Authenticate", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type AuthRPCServer struct {
	Impl Auth
}

func (s *AuthRPCServer) Authenticate(args interface{}, resp *bool) error {
	*resp = s.Impl.Authenticate()
	return nil
}

type AuthPlugin struct {
	Impl Auth
}

func (p *AuthPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &AuthRPCServer{Impl: p.Impl}, nil
}

func (AuthPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &AuthRPC{client: c}, nil
}
