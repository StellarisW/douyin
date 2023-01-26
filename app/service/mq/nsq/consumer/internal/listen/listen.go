package listen

import (
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/consumer/internal/config"
	"douyin/app/service/mq/nsq/consumer/internal/listen/chat"
	"douyin/app/service/mq/nsq/consumer/internal/listen/user"
	"douyin/app/service/mq/nsq/consumer/internal/listen/video"
	"douyin/app/service/mq/nsq/consumer/internal/svc"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
)

func NewServices(c config.Config) ([]service.Service, error) {
	var services []service.Service

	svcContext := svc.NewServiceContext(c)

	userServices, err := user.NewService(svcContext)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %v", errx.InitNsqService, err)
	}

	services = append(services, userServices...)

	videoServices, err := video.NewService(svcContext)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %v", errx.InitNsqService, err)
	}

	services = append(services, videoServices...)

	chatServices, err := chat.NewService(svcContext)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %v", errx.InitNsqService, err)
	}

	services = append(services, chatServices...)

	return services, nil
}
