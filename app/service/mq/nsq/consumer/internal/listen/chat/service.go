package chat

import (
	"douyin/app/common/mq/nsq"
	"douyin/app/service/mq/nsq/consumer/internal/svc"
	"douyin/app/service/mq/nsq/internal/consts"
	"github.com/zeromicro/go-zero/core/service"
)

func NewService(svcCtx *svc.ServiceContext) ([]service.Service, error) {
	chatConsumerService, err := nsq.NewConsumerService(
		consts.TopicChat,
		consts.ChannelChat,
		&Handler{
			ChatSysRpcClient: svcCtx.ChatSysRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	return []service.Service{
		chatConsumerService,
	}, nil
}
