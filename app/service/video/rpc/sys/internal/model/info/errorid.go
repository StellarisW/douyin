package info

const (
	// OprId

	ErrIdOprFeed = iota
	ErrIdOprGetPublishList
	ErrIdOprGetFavoriteList
	ErrIdOprGetCommentList
)

// 通用错误

const (
	ErrIdMysqlGet = iota
	ErrIdRequestRpcRes
	ErrIdRedisGet
	ErrIdRedisRange
	ErrIdCommon
)
