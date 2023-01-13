package logic

import (
	"context"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishReq) (resp *types.PublishRes, err error) {
	// todo: add your logic here and delete this line

	return
}
