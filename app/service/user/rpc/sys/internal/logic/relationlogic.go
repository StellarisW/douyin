package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/internal/model/consts"
	"douyin/app/service/user/rpc/sys/internal/model/relation"

	"douyin/app/service/user/rpc/sys/internal/svc"
	"douyin/app/service/user/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationLogic {
	return &RelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RelationLogic) Relation(in *pb.RelationReq) (*pb.RelationRes, error) {
	erx := l.svcCtx.RelationModel.Relation(l.ctx, in.SrcUserId, in.DstUserId, in.ActionType)
	if erx != nil {
		if erx.Code() == relation.ErrIdAlreadyFollow || erx.Code() == relation.ErrIdAlreadyUnfollow {
			return &pb.RelationRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicRelation,
					relation.ErrIdOprRelation,
					erx.Code(),
				),
				StatusMsg: erx.Error(),
			}, nil
		}
		return &pb.RelationRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.RelationRes{
		StatusCode: 0,
		StatusMsg:  "relation successfully",
	}, nil
}
