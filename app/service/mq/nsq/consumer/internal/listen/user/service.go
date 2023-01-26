package user

import (
	"douyin/app/common/mq/nsq"
	"douyin/app/service/mq/nsq/consumer/internal/svc"
	"douyin/app/service/mq/nsq/internal/consts"
	"github.com/zeromicro/go-zero/core/service"
)

func NewService(svcCtx *svc.ServiceContext) ([]service.Service, error) {
	userConsumerService, err := nsq.NewConsumerService(
		consts.TopicUserRelation,
		consts.ChannelUserRelation,
		&RelationHandler{
			UserSysRpcClient: svcCtx.UserSysRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	return []service.Service{
		userConsumerService,
	}, nil
}
