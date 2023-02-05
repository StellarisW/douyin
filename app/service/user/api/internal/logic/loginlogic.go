package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/api/internal/consts"
	"douyin/app/service/user/api/internal/consts/sign"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/rpc/sys/pb"
	"go.uber.org/zap"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	if ok, err := l.svcCtx.Regexp.UsernameReg.MatchString(req.Username); !ok || err != nil {
		return &types.LoginRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				sign.ErrIdInvalidUsername,
			),
			StatusMsg: sign.ErrInvalidUsername,
		}, nil
	}

	if ok, err := l.svcCtx.Regexp.PasswordReg.MatchString(req.Password); !ok || err != nil {
		return &types.LoginRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				sign.ErrIdInvalidUsername,
			),
			StatusMsg: sign.ErrInvalidPassword,
		}, nil
	}

	isExist, err := l.svcCtx.Rdb.SIsMember(l.ctx,
		user.RdbKeyRegisterSet,
		req.Username).Result()
	if err != nil {
		log.Logger.Error(errx.RedisGet, zap.Error(err))
		return &types.LoginRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				sign.ErrIdRedisGet,
			),
			StatusMsg: errx.Internal,
		}, nil
	}
	if !isExist {
		return &types.LoginRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				sign.ErrIdUsernameNotExist,
			),
			StatusMsg: sign.ErrUsernameNotExist,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.Login(l.ctx, &pb.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.LoginRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprLogin,
				sign.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.LoginRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.LoginRes{
		StatusCode: 0,
		StatusMsg:  "login successfully",
		UserId:     rpcRes.UserId,
		Token:      rpcRes.Token,
	}, nil
}
