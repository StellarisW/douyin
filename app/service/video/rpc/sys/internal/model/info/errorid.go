package info

const (
	// OprId

	ErrIdOprFeed = iota
	ErrIdOprGetPublishList
	ErrIdOprGetPublishCount
	ErrIdOprGetFavoriteList
	ErrIdOprGetFavoriteCount
	ErrIdOprGetCommentList
	ErrIdOprGetTotalFavorited
)

// 通用错误

const (
	ErrIdMysqlGet = iota
	ErrIdRequestRpcRes
	ErrIdRequestRpcReceive
	ErrIdRedisGet
	ErrIdRedisPipeExec
	ErrIdRedisRange
	ErrIdCommon
)
