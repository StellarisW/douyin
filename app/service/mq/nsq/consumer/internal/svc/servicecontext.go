package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	chatSys "douyin/app/service/chat/rpc/sys/sys"
	"douyin/app/service/mq/nsq/consumer/internal/config"
	userSys "douyin/app/service/user/rpc/sys/sys"
	videoSys "douyin/app/service/video/rpc/sys/sys"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	Viper *viper.Viper

	UserSysRpcClient  userSys.Sys
	VideoSysRpcClient videoSys.Sys
	ChatSysRpcClient  chatSys.Sys
}

func NewServiceContext(c config.Config) *ServiceContext {
	v, err := apollo.Common().GetViper("mq.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetViper, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		Viper: v,

		UserSysRpcClient:  userSys.NewSys(zrpc.MustNewClient(c.UserSysRpcClientConf)),
		VideoSysRpcClient: videoSys.NewSys(zrpc.MustNewClient(c.VideoSysRpcClientConf)),
		ChatSysRpcClient:  chatSys.NewSys(zrpc.MustNewClient(c.ChatSysRpcClientConf)),
	}
}
