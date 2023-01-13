package logic

import (
	"context"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowListLogic) GetFollowList(req *types.GetFollowListReq) (resp *types.GetFollowListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
