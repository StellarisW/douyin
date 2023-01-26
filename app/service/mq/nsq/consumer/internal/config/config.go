package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type ConsumerConf struct {
	Topic   string
	Channel string
}

type Config struct {
	service.ServiceConf

	UserSysRpcClientConf  zrpc.RpcClientConf
	VideoSysRpcClientConf zrpc.RpcClientConf
	ChatSysRpcClientConf  zrpc.RpcClientConf
}
