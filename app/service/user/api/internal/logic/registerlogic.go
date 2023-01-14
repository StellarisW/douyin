package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/user/api/internal/consts"
	"douyin/app/service/user/api/internal/consts/sign"
	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"
	"douyin/app/service/user/internal/sys"
	"douyin/app/service/user/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterRes, err error) {
	if ok, err := l.svcCtx.Regexp.UsernameReg.MatchString(req.Username); !ok || err != nil {
		return &types.RegisterRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprRegister,
				sign.ErrIdInvalidUsername,
			),
			StatusMsg: sign.ErrInvalidUsername,
		}, nil
	}

	if ok, err := l.svcCtx.Regexp.PasswordReg.MatchString(req.Password); !ok || err != nil {
		return &types.RegisterRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicSign,
				sign.ErrIdOprRegister,
				sign.ErrIdInvalidUsername,
			),
			StatusMsg: sign.ErrInvalidPassword,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.Register(l.ctx, &pb.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})
	if rpcRes.StatusCode != 0 {
		return &types.RegisterRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.RegisterRes{
		StatusCode: 0,
		StatusMsg:  "register successfully",
		UserId:     rpcRes.UserId,
		Token:      rpcRes.Token,
	}, nil
}
