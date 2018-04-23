package rpc

import (
	. "hdt_app_go/gateway_app/log"
	"hdt_app_go/gateway_app/rpc/client"
)

var (
	RpcClient *RPCClient
)

type RPCClient struct {
	Register *client.RegisterRPCCli
}

func NewRPCClient(name string) (c *RPCClient, err error) {
	register, err := client.NewRegisterRPCCli(name)
	if err != nil {
		Log.Err(err)
		return
	}
	c = &RPCClient{
		Register: register,
	}
	return
}

func InitRpc(name string) {
	c, err := NewRPCClient(name)
	if err != nil {
		Log.Err(err)
		return
	}
	RpcClient = c
}
