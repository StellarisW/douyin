package chat

const (
	ErrIdOprSendMessage = iota
	ErrIdOprGetMessageList
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
