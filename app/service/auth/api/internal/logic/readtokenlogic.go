package logic

import (
	"context"
	"douyin/app/service/auth/api/internal/model/token"

	"douyin/app/service/auth/api/internal/svc"
	"douyin/app/service/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadTokenLogic {
	return &ReadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadTokenLogic) ReadToken(req *types.ReadTokenReq) (resp *types.ReadTokenRes, err error) {
	tokenService := token.GetTokenService()
	payload, erx := tokenService.ReadToken(l.ctx,
		req.TokenValue)
	if erx != nil {
		return &types.ReadTokenRes{
			Code: erx.Code(),
			Msg:  erx.Error(),
		}, nil
	}

	return &types.ReadTokenRes{
		Code: 0,
		Msg:  "read token successfully",
		Data: types.ReadTokenResData{Payload: payload},
	}, nil
}
