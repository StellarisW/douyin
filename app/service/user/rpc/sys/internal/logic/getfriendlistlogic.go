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

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *pb.GetFriendListReq) (*pb.GetFriendListRes, error) {
	list, erx := l.svcCtx.RelationModel.GetFriendList(l.ctx, in.SrcUserId, in.DstUserId)
	if erx != nil {
		return &pb.GetFriendListRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprGetFriendList,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetFriendListRes{
		StatusCode: 0,
		StatusMsg:  "get friend list successfully",
		UserList:   list,
	}, nil
}
