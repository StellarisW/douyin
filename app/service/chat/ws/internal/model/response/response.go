package response

type Response struct {
	FromUserId int64  `json:"from_user_id"`
	MsgContent string `json:"msg_content"`
}
