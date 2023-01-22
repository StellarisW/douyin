package info

const (
	// OprId

	ErrIdOprFeed = iota
	ErrIdOprGetPublishList
	ErrIdOprGetFavoriteList
	ErrIdOprGetCommentList
)

// 通用操作

const (
	ErrIdRequestRpcReceiveSys = iota + 1
	ErrIdParseInt
	ErrIdCommon
)

const (
	ErrParseInt = "parse int failed"
)
