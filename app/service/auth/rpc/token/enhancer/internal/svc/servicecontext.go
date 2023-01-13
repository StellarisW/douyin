package svc

import (
	"douyin/app/service/auth/rpc/token/enhancer/internal/config"
	"douyin/app/service/auth/rpc/token/store/tokenstore"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	TokenStoreRpcClient tokenstore.TokenStore
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		TokenStoreRpcClient: tokenstore.NewTokenStore(zrpc.MustNewClient(c.TokenStoreRpcClientConf)),
	}
}
