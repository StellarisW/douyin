package service

import (
	"douyin/app/common/errx"
	"douyin/app/service/chat/api/internal/ws/internal/config"
	"douyin/app/service/chat/api/internal/ws/internal/service/listen"
	"douyin/app/service/chat/api/internal/ws/internal/svc"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
)

func New(c config.Config) ([]service.Service, error) {
	var services []service.Service

	svcContext := svc.NewServiceContext(c)

	listenServices, err := listen.NewService(svcContext.Viper.GetString(""))
	if err != nil {
		return nil, fmt.Errorf("%s, err: %v", errx.InitService, err)
	}

	services = append(services, listenServices)

	return services, nil
}
