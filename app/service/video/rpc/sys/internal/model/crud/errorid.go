package crud

const (
	// OprId

	ErrIdOprPublish = iota
	ErrIdOprFavorite
	ErrIdOprComment
)

// 通用错误

const (
	ErrIdMysqlInsert = iota
	ErrIdMinioPut
	ErrIdRedisAdd
	ErrIdRedisIncr
	ErrIdRedisRem
	ErrIdRedisDecr
	ErrIdInvalidActionType
	ErrIdMysqlDelete
	ErrIdCommon
)

const (
	ErrInvalidActionType = "invalid action type"
)

// Publish 操作

const (
	ErrIdInvalidContentType = iota + ErrIdOprComment
)

const (
	ErrInvalidContentType = "invalid content type"
)
