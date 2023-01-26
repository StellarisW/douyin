package crud

const (
	// OprId

	ErrIdOprPublish = iota
	ErrIdOprFavorite
	ErrIdOprComment
)

// 通用操作

const (
	ErrIdRequestRpcReceiveSys = iota + 1
	ErrIdParseInt
	ErrIdInvalidActionType
	ErrIdGetNsqProducer
	ErrIdNsqPublish
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
