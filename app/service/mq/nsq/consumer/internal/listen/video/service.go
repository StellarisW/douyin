package video

import (
	"douyin/app/common/mq/nsq"
	"douyin/app/service/mq/nsq/consumer/internal/svc"
	"douyin/app/service/mq/nsq/internal/consts"
	"github.com/zeromicro/go-zero/core/service"
)

func NewService(svcCtx *svc.ServiceContext) ([]service.Service, error) {
	videoFavoriteConsumerService, err := nsq.NewConsumerService(
		consts.TopicVideoFavorite,
		consts.ChannelVideoFavorite,
		&FavoriteHandler{
			VideoSysRpcClient: svcCtx.VideoSysRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	videoCommentConsumerService, err := nsq.NewConsumerService(
		consts.TopicVideoComment,
		consts.ChannelVideoComment,
		&FavoriteHandler{
			VideoSysRpcClient: svcCtx.VideoSysRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	return []service.Service{
		videoFavoriteConsumerService,
		videoCommentConsumerService,
	}, nil
}
