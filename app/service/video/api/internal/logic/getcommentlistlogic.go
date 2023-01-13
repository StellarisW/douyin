package logic

import (
	"context"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentListLogic) GetCommentList(req *types.GetCommentListReq) (resp *types.GetCommentListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
