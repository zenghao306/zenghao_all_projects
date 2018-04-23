package rpc

import (
	"hdt_app_go/appmeta/rpc/svr"
)

/*
var (
	RpcClient *RPCClient
)

type RPCClient struct {
	Register *client.RegisterRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	register, err := client.NewRegisterRPCCli()
	if err != nil {
		Log.Err(err)
		return
	}
	c = &RPCClient{
		Register: register,
	}
	return
}
*/
func InitRpc() {
	svr.NewRpcSvr()
}
