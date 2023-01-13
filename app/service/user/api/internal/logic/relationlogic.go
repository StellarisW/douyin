package logic

import (
	"context"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationLogic {
	return &RelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationLogic) Relation(req *types.RelationReq) (resp *types.RelationRes, err error) {
	// todo: add your logic here and delete this line

	return
}
