package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
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
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetFollowListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprGetFollowList,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.GetFollowListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprGetFollowList,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.GetFollowList(l.ctx, &pb.GetFollowListReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.GetFollowListRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprGetFollowList,
				relation.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.GetFollowListRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetFollowListRes{
		StatusCode: 0,
		StatusMsg:  "get follow list successfully",
		UserList:   rpcRes.UserList,
	}, nil
}
