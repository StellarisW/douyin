package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/middleware"
	"douyin/app/service/user/api/internal/consts"
	"douyin/app/service/user/api/internal/consts/relation"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/pb"
	"strconv"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListReq) (resp *types.GetFriendListRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetFriendListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.GetFriendListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.GetFriendList(l.ctx, &pb.GetFriendListReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes.StatusCode != 0 {
		return &types.GetFriendListRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetFriendListRes{
		StatusCode: 0,
		StatusMsg:  "get friend list successfully",
		UserList:   rpcRes.UserList,
	}, nil
}
