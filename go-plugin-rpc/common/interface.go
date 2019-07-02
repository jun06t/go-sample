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

const AuthPluginName = "authPlugin"

type Auth interface {
	Authenticate() bool
}

type AuthRPCClient struct {
	client *rpc.Client
}

func (a *AuthRPCClient) Authenticate() bool {
	var resp bool
	// using net/rpc client's synchronous call.
	err := a.client.Call("Plugin.Authenticate", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type AuthRPCServer struct {
	Impl Auth
}

// Authenticate conforms to the requirements of the net/rpc server method.
// https://golang.org/pkg/net/rpc/
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
	return &AuthRPCClient{client: c}, nil
}
