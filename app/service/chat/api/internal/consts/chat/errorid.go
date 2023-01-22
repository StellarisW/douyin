package chat

const (
	ErrIdOprSendMessage = iota
	ErrIdOprGetMessageList
)

// 通用操作

const (
	ErrIdRequestRpcReceiveSys = iota + 1
	ErrIdParseInt
	ErrIdInvalidActionType
	ErrIdCommon
)

const (
	ErrParseInt          = "parse int failed"
	ErrInvalidActionType = "invalid action type"
)
