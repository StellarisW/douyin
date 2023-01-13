package main

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"

	"douyin/app/service/auth/rpc/token/store/internal/config"
	"douyin/app/service/auth/rpc/token/store/internal/server"
	"douyin/app/service/auth/rpc/token/store/internal/svc"
	"douyin/app/service/auth/rpc/token/store/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "auth.rpc.auth.tokenstore"

var c config.Config

func main() {
	serviceFullName := utils.GetServiceFullName(serviceName)

	// 初始化日志管理器
	err := log.InitLogger(serviceFullName + "-data")
	if err != nil {
		panic("initialize logger failed")
	}

	logWriter := log.GetLogxZapWriter(log.Logger)

	logx.MustSetup(log.GetLogxConfig(utils.GetServiceFullName(serviceName), "error"))
	logx.SetWriter(logWriter)

	// 初始化配置管理器
	err = apollo.InitClient()
	if err != nil {
		log.Logger.Fatal(errx.InitAgolloClient, zap.Error(err))
	}

	// 初始化微服务设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = apollo.Common().UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
	if err != nil {
		log.Logger.Fatal(errx.UnmarshalServiceConfig, zap.Error(err))
	}

	ctx := svc.NewServiceContext(c)

	// 启动微服务服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterTokenStoreServer(grpcServer, server.NewTokenStoreServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 注册服务到consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)

	defer s.Stop()

	log.Logger.Info("starting rpc server...", zap.String("addr", c.ListenOn))
	s.Start()
}
