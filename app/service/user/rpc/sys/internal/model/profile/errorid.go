package profile

const (
	// OprId

	ErrIdOprGetProfile = iota
)

const (
	// 通用错误

	ErrIdMysqlGet = iota
	ErrIdRequestRpcRes
	ErrIdRequestRpcReceive
	ErrIdRedisGet
	ErrIdCommon
)

// GetProfile 操作

const (
	ErrIdUserNotFound = iota + ErrIdCommon
)

const (
	ErrUserNotFound = "user not found"
)
