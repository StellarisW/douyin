package consts

const (
	// LogicId

	ErrIdLogic = iota
)

const (
	// OprId

	ErrIdOprGetTokenByAuth = iota
	ErrIdOprReadToken
	ErrIdOprRefreshToken
)

// 通用错误

const (
	ErrIdParseHttpRequest = iota
	ErrIdGrantTypeNotSupported
	ErrIdRpcToken2Oauth2Token
	ErrIdInvalidAuthorization
	ErrIdRpcGenerateOauth2Token
	ErrIdCommon
)

const (
	ErrGrantTypeNotSupported = "grant type not supported"
	ErrInvalidAuthorization  = "invalid authorization"
)

// GetTokenByAuth 操作

const (
	ErrIdGrantTokenByAuth = iota + ErrIdCommon
)

const (
	ErrGrantTokenByAuth = "grant token by auth failed"
)

// ReadToken 操作

const (
	ErrIdReadToken = iota + ErrIdCommon
	ErrIdRpcReadToken
)

const (
	ErrReadToken = "read token failed"
)

// RefreshToken 操作

const (
	ErrIdInvalidRefreshToken = iota + ErrIdCommon
	ErrIdGrantTokenByRefreshToken
	ErrIdRpcRefreshOauth2Token
)

const (
	ErrInvalidRefreshToken      = "invalid refresh token"
	ErrGrantTokenByRefreshToken = "grant token by refresh token failed"
)
