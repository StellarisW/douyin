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

type GetFollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowerListLogic) GetFollowerList(req *types.GetFollowerListReq) (resp *types.GetFollowerListRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetFollowerListRes{
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
		return &types.GetFollowerListRes{
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

	rpcRes, _ := l.svcCtx.SysRpcClient.GetFollowerList(l.ctx, &pb.GetFollowerListReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes.StatusCode != 0 {
		return &types.GetFollowerListRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetFollowerListRes{
		StatusCode: 0,
		StatusMsg:  "get follower list successfully",
		UserList:   rpcRes.UserList,
	}, nil
}
