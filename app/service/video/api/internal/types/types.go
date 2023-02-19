// Code generated by goctl. DO NOT EDIT.
package types

type Profile struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FollowCount    uint64 `json:"follow_count"`
	FollowerCount  uint64 `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited uint64 `json:"total_favorited"`
	WorkCount      uint64 `json:"work_count"`
	FavoriteCount  uint64 `json:"favorite_count"`
}

type Video struct {
	Id            int64    `json:"id"`
	Author        *Profile `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount uint64   `json:"favorite_count"`
	CommentCount  uint64   `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}

type Comment struct {
	Id         int64    `json:"id"`
	User       *Profile `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}

type FeedReq struct {
	LastestTime string `form:"latest_time,optional"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token       string `form:"token,optional"`       // 可选参数，登录用户设置
}

type FeedRes struct {
	StatusCode uint32      `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`           // 返回状态描述
	NextTime   int64       `json:"next_time,omitempty"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList  interface{} `json:"video_list,omitempty"` // 视频列表
}

type PublishReq struct {
	Token string `form:"token"` // 用户鉴权token
	Title string `form:"title"` // 视频标题
}

type PublishRes struct {
	StatusCode uint32 `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type GetPublishListReq struct {
	Token  string `form:"token"`
	UserId string `form:"user_id"`
}

type GetPublishListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
}

type FavoriteReq struct {
	Token      string `form:"token"`       // 用户鉴权token
	VideoId    string `form:"video_id"`    // 视频id
	ActionType string `form:"action_type"` // 1-点赞，2-取消点赞
}

type FavoriteRes struct {
	StatusCode uint32 `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type GetFavoriteListReq struct {
	UserId string `form:"user_id"` // 用户id
	Token  string `form:"token"`   // 用户鉴权token
}

type GetFavoriteListRes struct {
	StatusCode uint32      `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`           // 返回状态描述
	VideoList  interface{} `json:"video_list,omitempty"` // 用户点赞视频列表
}

type CommentReq struct {
	Token       string `form:"token"`                 // 用户鉴权token
	VideoId     string `form:"video_id"`              // 视频id
	ActionType  string `form:"action_type"`           // 1-发布评论，2-删除评论
	CommentText string `form:"comment_text,optional"` // 用户填写的评论内容，在action_type=1的时候使用
	CommentId   string `form:"comment_id,optional"`   // 要删除的评论id，在action_type=2的时候使用
}

type CommentRes struct {
	StatusCode uint32      `json:"status_code"`       // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`        // 返回状态描述
	Comment    interface{} `json:"comment,omitempty"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

type GetCommentListReq struct {
	Token   string `form:"token"`    // 用户鉴权token
	VideoId string `form:"video_id"` // 视频id
}

type GetCommentListRes struct {
	StatusCode  uint32      `json:"status_code"`            // 状态码，0-成功，其他值-失败
	StatusMsg   string      `json:"status_msg"`             // 返回状态描述
	CommentList interface{} `json:"comment_list,omitempty"` // 评论列表
}
