package token

import (
	apollo "douyin/app/common/config"
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	GrantByAuth         = "authorization"
	GrantByRefreshToken = "refresh_token"

	EnhancerRpcClientConf = "Api.TokenEnhancerRpcClientConf.Target"
)

var (
	tokenService TokenService
	tokenGranter TokenGranter
)

func InitTokenService() (err error) {
	v, err := apollo.Common().GetViper("auth.yaml")
	if err != nil {
		return fmt.Errorf("get viper failed, %v", err)
	}

	tokenService = NewRpcTokenService(zrpc.RpcClientConf{Target: v.GetString(EnhancerRpcClientConf)})

	return nil
}

func InitTokenGranter() (err error) {
	authorizationTokenGranter := &AuthorizationTokenGranter{
		SupportGrantType: GrantByAuth,
		ClientSecret:     ClientSecrets,
		TokenService:     tokenService,
	}

	tokenGranter = NewComposeTokenGranter(map[string]TokenGranter{
		GrantByAuth: authorizationTokenGranter,
	})
	return nil
}

func GetTokenService() TokenService {
	return tokenService
}

func GetTokenGranter() TokenGranter {
	return tokenGranter
}
