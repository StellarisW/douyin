package logic

import (
	"context"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPublishListLogic) GetPublishList(req *types.GetPublishListReq) (resp *types.GetPublishListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
