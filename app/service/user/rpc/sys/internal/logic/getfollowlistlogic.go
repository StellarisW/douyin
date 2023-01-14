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

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowListLogic) GetFollowList(in *pb.GetFollowListReq) (*pb.GetFollowListRes, error) {
	list, erx := l.svcCtx.RelationModel.GetFollowList(l.ctx, in.SrcUserId, in.DstUserId)
	if erx != nil {
		return &pb.GetFollowListRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprGetFollowList,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetFollowListRes{
		StatusCode: 0,
		StatusMsg:  "get follow list successfully",
		UserList:   list,
	}, nil
}
