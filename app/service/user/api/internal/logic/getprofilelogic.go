package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"douyin/app/service/user/api/internal/consts"
	"douyin/app/service/user/api/internal/consts/profile"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/pb"
	"strconv"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req *types.GetProfileReq) (resp *types.GetProfileRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetProfileRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicProfile,
				profile.ErrIdOprGetProfile,
				profile.ErrIdParseInt,
			),
			StatusMsg: profile.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.GetProfileRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicProfile,
				profile.ErrIdOprGetProfile,
				profile.ErrIdParseInt,
			),
			StatusMsg: profile.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.GetProfile(l.ctx, &pb.GetProfileReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.GetProfileRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicProfile,
				profile.ErrIdOprGetProfile,
				profile.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.GetProfileRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetProfileRes{
		StatusCode: 0,
		StatusMsg:  "get profile successfully",
		User: types.Profile{
			Id:            rpcRes.User.Id,
			Name:          rpcRes.User.Name,
			FollowCount:   rpcRes.User.FollowCount,
			FollowerCount: rpcRes.User.FollowerCount,
			IsFollow:      rpcRes.User.IsFollow,
		},
	}, nil
}
