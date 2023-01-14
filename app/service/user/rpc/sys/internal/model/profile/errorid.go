package profile

const (
	// OprId

	ErrIdOprGetProfile = iota
)

const (
	// 通用错误

	ErrIdMysqlGet = iota
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
