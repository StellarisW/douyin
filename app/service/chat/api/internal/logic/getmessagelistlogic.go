package logic

import (
	"context"

	"douyin/app/service/chat/api/internal/svc"
	"douyin/app/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.GetMessageListReq) (resp *types.GetMessageListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
