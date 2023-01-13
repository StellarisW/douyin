package main

import (
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"go.uber.org/zap"

	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/utils"

	"douyin/app/service/user/api/internal/config"
	"douyin/app/service/user/api/internal/handler"
	"douyin/app/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

const serviceName = "user.api"

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

	// 启动微服务服务器
	server := rest.MustNewServer(c.RestConf)

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	log.Logger.Info("starting server...", zap.String("host", c.Host), zap.Int("port", c.Port))
	server.Start()
}
