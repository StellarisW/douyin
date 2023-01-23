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
	ErrIdRedisGet
	ErrIdRedisIncr
	ErrIdRedisRem
	ErrIdRedisDecr
	ErrIdRedisPipeExec
	ErrIdInvalidActionType
	ErrIdMysqlDelete
	ErrIdCommon
)

const (
	ErrInvalidActionType = "invalid action type"
)

// Publish 操作

const (
	ErrIdInvalidContentType = iota + ErrIdCommon
)

const (
	ErrInvalidContentType = "invalid content type"
)

// Favorite 操作

const (
	ErrIdAlreadyLike = iota + ErrIdCommon
	ErrIdAlreadyDisLike
)

const (
	ErrAlreadyLike    = "already liked"
	ErrAlreadyDisLike = "already disliked"
)
