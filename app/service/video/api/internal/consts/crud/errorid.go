package crud

const (
	// OprId

	ErrIdOprPublish = iota
	ErrIdOprFavorite
	ErrIdOprComment
)

// 通用操作

const (
	ErrIdParseInt = iota + 1
	ErrIdInvalidActionType
	ErrIdCommon
)

const (
	ErrParseInt          = "parse int failed"
	ErrInvalidActionType = "invalid action type"
)

// Publish 操作

const (
	ErrIdOpenFile = iota + ErrIdCommon
	ErrIdReadBytes
	ErrIdInvalidVideoType
)

const (
	ErrOpenFile         = "open file failed"
	ErrInvalidVideoType = "invalid video type"
)
