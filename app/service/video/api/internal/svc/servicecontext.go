package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/log"
	"douyin/app/service/video/api/internal/config"
	"douyin/app/service/video/rpc/sys/sys"
	"github.com/StellarisW/go-sensitive"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config config.Config

	JWTAuthMiddleware rest.Middleware
	CORSMiddleware    rest.Middleware

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

	err := filterManager.GetStore().LoadDictPath("./manifest/config/dic/default_dict.txt")
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

		SysRpcClient: sys.NewSys(
			zrpc.MustNewClient(
				c.SysRpcClientConf,
				zrpc.WithDialOption(
					grpc.WithDefaultCallOptions(
						grpc.MaxCallRecvMsgSize(256<<32), // 最大上传 256MB 视频
					),
				),
			),
		),

		Filter: filterManager,
	}
}
