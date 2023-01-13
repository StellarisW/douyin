package token

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/service/auth/internal/auth"
	"encoding/base64"
	"strings"
)

type (
	// TokenGranter 令牌颁发器接口
	TokenGranter interface {
		Grant(ctx context.Context, grantType string, auth string, obj string) (*auth.Token, errx.Error)
	}

	// ComposeTokenGranter 集成令牌颁发器
	ComposeTokenGranter struct {
		TokenGrantDict map[string]TokenGranter
	}

	// AuthorizationTokenGranter 通过认证令牌的颁发器结构体
	AuthorizationTokenGranter struct {
		SupportGrantType string
		ClientSecret     map[string]string
		TokenService     TokenService
	}
)

func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{
		TokenGrantDict: tokenGrantDict,
	}
}

func (tokenGranter *ComposeTokenGranter) Grant(ctx context.Context, grantType string, auth string, obj string) (*auth.Token, errx.Error) {
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]

	if dispatchGranter == nil {
		return nil, errGrantTypeNotSupported
	}

	return dispatchGranter.Grant(ctx, grantType, auth, obj)
}

func (tokenGranter *AuthorizationTokenGranter) Grant(ctx context.Context,
	_ string, auth string, obj string) (*auth.Token, errx.Error) {
	// 解析 Basic 认证头
	clientId, clientSecret, ok := parseBasicAuth(auth)
	if !ok || clientSecret != tokenGranter.ClientSecret[clientId] {
		return nil, errInvalidAuthorization
	}

	// 签发token
	return tokenGranter.TokenService.GenerateToken(ctx, obj, clientId)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic ZndlYWZjdjQyM2M6MzJjNDIzYzQ=" returns ("clientId", "clientSecret", ture).
func parseBasicAuth(auth string) (clientId, clientSecret string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || (auth[:len(prefix)] != prefix) {
		return "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	seq := strings.Split(cs, ":")
	if len(seq) == 2 {
		ok = true
	}
	clientId = seq[0]
	clientSecret = seq[1]
	if !ok {
		return "", "", false
	}
	return clientId, clientSecret, true
}
