package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	chatSys "douyin/app/service/chat/rpc/sys/sys"
	"douyin/app/service/mq/nsq/consumer/internal/config"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	Viper *viper.Viper

	ChatSysRpcClient chatSys.Sys
}

func NewServiceContext(c config.Config) *ServiceContext {
	v, err := apollo.Common().GetViper("mq.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetViper, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		Viper: v,

		ChatSysRpcClient: chatSys.NewSys(zrpc.MustNewClient(c.ChatSysRpcClientConf)),
	}
}
