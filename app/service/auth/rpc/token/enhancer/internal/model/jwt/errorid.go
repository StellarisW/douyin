package jwt

import "douyin/app/service/auth"

const (
	// OprId

	ErrIdOprGenerateToken = iota
	ErrIdOprReadToken
	ErrIdOprRefreshToken
)

/*
------------------------------
*/
// 通用错误

const (
	ErrIdTokenBuild = iota
	ErrIdTokenSerialize
	ErrIdInvalidSignature
	ErrIdInvalidKey
	ErrIdTokenNotValidYet
	ErrIdTokenExpired
	ErrIdCommon
)

const (
	ErrTokenBuild       = "build token failed"
	ErrTokenSerialize   = "serialize token failed"
	ErrGenerateToken    = "generate token failed"
	ErrInvalidSignature = auth.ErrInvalidSignature
	ErrInvalidKey       = auth.ErrInvalidKey
	ErrTokenNotValidYet = auth.ErrTokenNotValidYet
	ErrTokenExpired     = auth.ErrTokenExpired
)

// GenerateToken 操作

// ReadToken 操作

const (
	ErrIdParseToken = iota + ErrIdCommon
)

const (
	ErrParseToken = "parse token failed"
)

// RefreshToken 操作

const (
	ErrIdParseRefreshToken = iota + ErrIdCommon
)

const (
	ErrParseRefreshToken = "parse refresh_token failed"
)
