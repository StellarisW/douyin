package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/internal/model/consts"
	"douyin/app/service/user/rpc/sys/internal/model/sign"

	"douyin/app/service/user/rpc/sys/internal/svc"
	"douyin/app/service/user/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterRes, error) {
	userId, token, erx := l.svcCtx.SignModel.Register(l.ctx, in.Username, in.Password)
	if erx != nil {
		if erx.Code() == sign.ErrIdUsernameExists {
			return &pb.RegisterRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicSign,
					sign.ErrIdOprRegister,
					sign.ErrIdUsernameExists,
				),
				StatusMsg: sign.ErrUsernameExists,
			}, nil
		}

		return &pb.RegisterRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicSign,
				sign.ErrIdOprRegister,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.RegisterRes{
		StatusCode: 0,
		StatusMsg:  "register successfully",
		UserId:     userId,
		Token:      token,
	}, nil
}
