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

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginRes, error) {
	userId, token, erx := l.svcCtx.SignModel.Login(l.ctx, in.Username, in.Password)
	if erx != nil {
		if erx.Code() == sign.ErrIdWrongUsernameOrPassword {
			return &pb.LoginRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicSign,
					sign.ErrIdOprLogin,
					sign.ErrIdWrongUsernameOrPassword,
				),
				StatusMsg: sign.ErrWrongUsernameOrPassword,
			}, nil
		}

		return &pb.LoginRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.LoginRes{
		StatusCode: 0,
		StatusMsg:  "register successfully",
		UserId:     userId,
		Token:      token,
	}, nil
}
