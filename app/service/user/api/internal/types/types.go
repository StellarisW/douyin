// Code generated by goctl. DO NOT EDIT.
package types

type Profile struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type RegisterReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterRes struct {
	StatusCode uint32 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
}

type LoginReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginRes struct {
	StatusCode uint32 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
}

type GetProfileReq struct {
	UserId string `form:"user_id"`
	Token  string `form:"token"`
}

type GetProfileRes struct {
	StatusCode uint32  `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	User       Profile `json:"user"`
}

type RelationReq struct {
	Token      string `form:"token"`
	ToUserId   string `form:"to_user_id"`
	ActionType string `form:"action_type"`
}

type RelationRes struct {
	StatusCode uint32 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetFollowListReq struct {
	UserId string `form:"user_id"`
	Token  string `form:"token"`
}

type GetFollowListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	UserList   interface{} `json:"user_list"`
}

type GetFollowerListReq struct {
	UserId string `form:"user_id"`
	Token  string `form:"token"`
}

type GetFollowerListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	UserList   interface{} `json:"user_list"`
}

type GetFriendListReq struct {
	UserId string `form:"user_id"`
	Token  string `form:"token"`
}

type GetFriendListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	UserList   interface{} `json:"user_list"`
}
