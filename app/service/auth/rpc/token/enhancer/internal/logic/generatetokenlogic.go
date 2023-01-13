package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/auth/internal/sys"
	"douyin/app/service/auth/rpc/token/enhancer/internal/consts"
	"douyin/app/service/auth/rpc/token/enhancer/internal/model/jwt"
	"douyin/app/service/auth/rpc/token/store/tokenstore"
	"go.uber.org/zap"
	"time"

	"douyin/app/service/auth/rpc/token/enhancer/internal/svc"
	"douyin/app/service/auth/rpc/token/enhancer/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenRes, error) {
	existTokenRes, _ := l.svcCtx.TokenStoreRpcClient.GetToken(l.ctx, &tokenstore.GetTokenReq{
		Username: in.Username,
	})

	if existTokenRes.Code == 0 {
		// 获取到已经存在的令牌
		if !time.Unix(existTokenRes.Data.Token.ExpiresAt, 0).Before(time.Now()) {
			return &pb.GenerateTokenRes{
				Code: 0,
				Msg:  "get exist token",
				Data: &pb.GenerateTokenRes_Data{Token: &pb.Token{
					TokenValue: existTokenRes.Data.Token.TokenValue,
					ExpiresAt:  existTokenRes.Data.Token.ExpiresAt,
				}},
			}, nil
		}
	}

	// 生成新令牌
	token, erx := jwt.GenerateToken(in.Username,
		in.ClientId,
		in.Username,
	)
	if erx != nil {
		log.Logger.Error(jwt.ErrGenerateToken, zap.Error(erx))
		return &pb.GenerateTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcEnhancer,
				consts.ErrIdLogic,
				jwt.ErrIdOprGenerateToken,
				erx.Code(),
			),
			Msg: errx.Internal,
		}, nil
	}

	storeTokenRes, _ := l.svcCtx.TokenStoreRpcClient.StoreToken(l.ctx,
		&tokenstore.StoreTokenReq{
			Username: in.Username,
			Token: &tokenstore.Token{
				TokenValue: token.TokenValue,
				ExpiresAt:  token.ExpiresAt,
			},
		})
	if storeTokenRes.Code != 0 {
		return &pb.GenerateTokenRes{
			Code: storeTokenRes.Code,
			Msg:  storeTokenRes.Msg,
		}, nil
	}

	return &pb.GenerateTokenRes{
		Code: 0,
		Msg:  "generate token successfully",
		Data: &pb.GenerateTokenRes_Data{Token: &pb.Token{
			TokenValue: token.TokenValue,
			ExpiresAt:  token.ExpiresAt,
		}},
	}, nil
}
