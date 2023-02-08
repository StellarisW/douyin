package middleware

import (
	"douyin/app/common/douyin"
)

const (
	sysId = douyin.SysIdMiddleware

	serviceIdCORS = iota - 1
	serviceIdJWT
	serviceIdSentinel
)

// JWT 中间件

const (
	errIdGetAccessTokenByCookie = iota
	errIdRequestHttpSendAuth
	errIdGetRefreshTokenByCookie
	errIdRefreshToken
	errIdReadToken
	errIdInvalidToken
)

const (
	errGetAccessTokenByCookie  = "not logged in"
	errGetRefreshTokenByCookie = "not logged in"
	errRefreshToken            = "internal err"
	errReadToken               = "read token failed"
)

type ctxKey int

const (
	KeyUserId ctxKey = iota
	KeyScope
)

// Sentinel 中间件

const (
	errIdFlow = iota
)

const (
	errFlow = "pass flow rule failed"
)
