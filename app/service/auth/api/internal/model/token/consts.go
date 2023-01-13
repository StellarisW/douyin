package token

import (
	"douyin/app/common/errx"
	"douyin/app/service/auth/api/internal/consts"
)

var (
	// token_granter 错误

	errGrantTypeNotSupported = errx.New(consts.ErrIdGrantTypeNotSupported, consts.ErrGrantTypeNotSupported)
	errInvalidAuthorization  = errx.New(consts.ErrIdInvalidAuthorization, consts.ErrInvalidAuthorization)
	errInvalidRefreshToken   = errx.New(consts.ErrIdInvalidRefreshToken, consts.ErrInvalidRefreshToken)
)

var (
	// token_service 错误

	errRpcGenerateOauth2Token = errx.New(consts.ErrIdRpcGenerateOauth2Token, errx.RequestRpcRes)

	errRpcReadToken = errx.New(consts.ErrIdRpcReadToken, errx.RequestRpcRes)

	errRpcRefreshOauth2Token = errx.New(consts.ErrIdRpcRefreshOauth2Token, errx.RequestRpcRes)
)
