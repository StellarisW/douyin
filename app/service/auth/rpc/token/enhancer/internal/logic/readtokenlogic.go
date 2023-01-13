package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/auth/internal/sys"
	"douyin/app/service/auth/rpc/token/enhancer/internal/consts"
	"douyin/app/service/auth/rpc/token/enhancer/internal/model/jwt"
	"go.uber.org/zap"

	"douyin/app/service/auth/rpc/token/enhancer/internal/svc"
	"douyin/app/service/auth/rpc/token/enhancer/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadTokenLogic {
	return &ReadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReadTokenLogic) ReadToken(in *pb.ReadTokenReq) (*pb.ReadTokenRes, error) {
	token, erx := jwt.ParseToken(in.TokenValue)
	if erx != nil {
		log.Logger.Debug(jwt.ErrParseToken, zap.Error(erx))
		return &pb.ReadTokenRes{
			Code: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcEnhancer,
				consts.ErrIdLogic,
				jwt.ErrIdOprReadToken,
				erx.Code(),
			),
			Msg: erx.Error(), // 将解析token的错误信息原封不动的返回
		}, nil
	}

	return &pb.ReadTokenRes{
		Code: 0,
		Msg:  "read token successfully",
		Data: &pb.ReadTokenRes_Data{Payload: token},
	}, nil
}
