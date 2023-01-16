package main

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/mq/nsq/consumer/internal/config"
	"douyin/app/service/mq/nsq/consumer/internal/listen"
	"douyin/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"go.uber.org/zap"
)

const serviceName = "mq.nsq.consumer"

var c config.Config

func main() {
	serviceFullName := utils.GetServiceFullName(serviceName)

	// 初始化日志管理器
	err := log.InitLogger(serviceFullName + "-data")
	if err != nil {
		panic("initialize logger failed")
	}

	logWriter := log.GetLogxZapWriter(log.Logger)

	logx.MustSetup(log.GetLogxConfig(serviceFullName, "error"))
	logx.SetWriter(logWriter)

	// 初始化配置管理器
	err = apollo.InitClient()
	if err != nil {
		log.Logger.Fatal("initialize Apollo Client failed.", zap.Error(err))
	}

	// 初始化消息队列设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = apollo.Common().UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
	if err != nil {
		log.Logger.Fatal(errx.UnmarshalServiceConfig, zap.Error(err))
	}

	// 初始化 log, trace
	err = c.SetUp()
	if err != nil {
		log.Logger.Fatal("initialize go-zero internal log service failed", zap.Error(err))
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	consumerServices, err := listen.NewServices(c)
	if err != nil {
		log.Logger.Fatal(errx.InitNsqConsumerService, zap.Error(err))
	}

	for _, consumerService := range consumerServices {
		serviceGroup.Add(consumerService)
	}

	serviceGroup.Start()
}
