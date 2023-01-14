package relation

const (
	// OprId

	ErrIdOprRelation = iota
	ErrIdOprGetFollowList
	ErrIdOprGetFollowerList
	ErrIdOprGetFriendList
)

// 通用操作

const (
	ErrIdParseInt = iota + 1
	ErrIdCommon
)

const (
	ErrParseInt = "parse int failed"
)

// Relation 操作

const (
	ErrIdInvalidActionType = iota + ErrIdCommon
	ErrIdCannotFollowYourself
)

const (
	ErrInvalidActionType    = "invalid action type"
	ErrCannotFollowYourself = "cannot follow yourself"
)
