package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/internal/model/consts"
	"douyin/app/service/user/rpc/sys/internal/model/profile"

	"douyin/app/service/user/rpc/sys/internal/svc"
	"douyin/app/service/user/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProfileLogic) GetProfile(in *pb.GetProfileReq) (*pb.GetProfileRes, error) {
	userInfo, erx := l.svcCtx.ProfileModel.GetProfile(l.ctx, in.SrcUserId, in.DstUserId)
	if erx != nil {
		if erx.Code() == profile.ErrIdUserNotFound {
			return &pb.GetProfileRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicProfile,
					profile.ErrIdOprGetProfile,
					profile.ErrIdUserNotFound,
				),
				StatusMsg: profile.ErrUserNotFound,
			}, nil
		}
		return &pb.GetProfileRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicProfile,
				profile.ErrIdOprGetProfile,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetProfileRes{
		StatusCode: 0,
		StatusMsg:  "get profile successfully",
		User:       userInfo,
	}, nil
}
