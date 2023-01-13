package consts

const (
	// LogicId

	ErrIdLogic = iota
)

const (
	// OprId

	ErrIdOprStoreToken = iota
	ErrIdOprGetToken
	ErrIdOprRemoveToken
)

// 通用错误

const (
	ErrIdEmptyParam = iota
	ErrIdCommon     // 通用错误末尾id, 用于其他操作的第一个错误确定开头的id
)

const (
	ErrEmptyParam = "empty param"
)

// StoreToken 操作

const (
	ErrIdJsonMarshalOauth2Token = iota + ErrIdCommon
	ErrIdRedisSetOauth2Token
)

// GetToken 操作

const (
	ErrIdTokenNotFound = iota + ErrIdCommon
	ErrIdRedisGetOauth2Token
	ErrIdJsonUnmarshalOauth2Token
)

const (
	ErrTokenNotFound = "token not found"
)

// RemoveToken 操作

const (
	ErrIdRedisDelOauth2Token = iota + ErrIdCommon
)
