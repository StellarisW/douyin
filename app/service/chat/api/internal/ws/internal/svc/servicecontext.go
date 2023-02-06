package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/chat/api/internal/ws/internal/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	Viper *viper.Viper
}

func NewServiceContext(c config.Config) *ServiceContext {
	v, err := apollo.Common().GetViper("chat.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetViper, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		Viper: v,
	}
}
