package request

type Request struct {
	UserId     int64  `json:"user_id"`
	ToUserId   int64  `json:"to_user_id"`
	MsgContent string `json:"msg_content"`
}
