package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/log"
	"douyin/app/service/chat/api/internal/config"
	"douyin/app/service/chat/rpc/sys/sys"
	"github.com/StellarisW/go-sensitive"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	CORSMiddleware    rest.Middleware
	JWTAuthMiddleware rest.Middleware

	SysRpcClient sys.Sys

	Filter *sensitive.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	filterManager := sensitive.NewFilter(
		sensitive.StoreOption{
			Type: sensitive.StoreMemory,
		},
		sensitive.FilterOption{
			Type: sensitive.FilterDfa,
		})

	err := filterManager.GetStore().LoadDictPath("./manifest/config/dict/default_dict.txt")
	if err != nil {
		log.Logger.Fatal("initialize filter failed", zap.Error(err))
	}

	corsMiddleware, err := apollo.Middleware().NewCORSMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize corsMiddleware failed.", zap.Error(err))
	}

	JWTAuthMiddleware, err := apollo.Middleware().NewJWTAuthMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize JWTAuthMiddleware failed.", zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		CORSMiddleware:    corsMiddleware.Handle,
		JWTAuthMiddleware: JWTAuthMiddleware.Handle,

		SysRpcClient: sys.NewSys(zrpc.MustNewClient(c.SysRpcClientConf)),

		Filter: filterManager,
	}
}
