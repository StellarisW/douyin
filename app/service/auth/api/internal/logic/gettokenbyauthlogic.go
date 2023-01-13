package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/auth/api/internal/consts"
	"douyin/app/service/auth/api/internal/model/token"
	"douyin/app/service/auth/internal/sys"
	"go.uber.org/zap"

	"douyin/app/service/auth/api/internal/svc"
	"douyin/app/service/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenByAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByAuthLogic {
	return &GetTokenByAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByAuthLogic) GetTokenByAuth(req *types.GetTokenByAuthReq) (resp *types.GetTokenByAuthRes, err error) {
	tokenGranter := token.GetTokenGranter()
	oauth2Token, erx := tokenGranter.Grant(l.ctx, token.GrantByAuth, req.Authorization, req.Obj)
	if erx != nil {
		log.Logger.Debug(consts.ErrGrantTokenByAuth, zap.Error(erx))
		return &types.GetTokenByAuthRes{
			Code: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				consts.ErrIdOprGetTokenByAuth,
				erx.Code(),
			),
			Msg: erx.Error(),
		}, nil
	}

	return &types.GetTokenByAuthRes{
		Code: 0,
		Msg:  "get token successfully",
		Data: types.GetTokenByAuthResData{Oauth2Token: oauth2Token},
	}, nil
}
