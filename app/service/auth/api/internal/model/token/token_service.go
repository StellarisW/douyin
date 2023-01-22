package token

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/service/auth/internal/auth"
	"douyin/app/service/auth/rpc/token/enhancer/tokenenhancer"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	// TokenService 令牌服务接口
	TokenService interface {
		GenerateToken(ctx context.Context, subject string, audience string) (*auth.Token, errx.Error)
		ReadToken(ctx context.Context, tokenValue string) (string, errx.Error)
	}

	// DefaultTokenService 默认令牌服务结构体
	DefaultTokenService struct {
	}

	// RpcTokenService rpc 令牌服务结构体
	RpcTokenService struct {
		TokenEnhancerClient tokenenhancer.TokenEnhancer
	}
)

func NewRpcTokenService(enhancerConf zrpc.RpcClientConf) TokenService {
	return &RpcTokenService{
		TokenEnhancerClient: tokenenhancer.NewTokenEnhancer(zrpc.MustNewClient(enhancerConf)),
	}
}

// GenerateToken 生成令牌
func (tokenService *RpcTokenService) GenerateToken(ctx context.Context, subject string, audience string) (*auth.Token, errx.Error) {
	rpcRes, _ := tokenService.TokenEnhancerClient.GenerateToken(ctx,
		&tokenenhancer.GenerateTokenReq{
			Username: subject,
			ClientId: audience,
		})
	if rpcRes == nil || rpcRes.Code != 0 {
		return nil, errRpcGenerateOauth2Token
	}

	return &auth.Token{
		TokenValue: rpcRes.Data.Token.TokenValue,
		ExpiresAt:  rpcRes.Data.Token.ExpiresAt,
	}, nil
}

// ReadToken 读取令牌
func (tokenService *RpcTokenService) ReadToken(ctx context.Context, tokenValue string) (string, errx.Error) {
	rpcRes, _ := tokenService.TokenEnhancerClient.ReadToken(ctx,
		&tokenenhancer.ReadTokenReq{
			TokenValue: tokenValue,
		})
	if rpcRes == nil || rpcRes.Code != 0 {
		return "", errx.New(rpcRes.Code, rpcRes.Msg)
	}

	return rpcRes.Data.Payload, nil
}
