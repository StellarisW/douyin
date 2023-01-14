// Code generated by goctl. DO NOT EDIT.
package types

type Profile struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
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
	LastestTime string `form:"lastest_time,optional"`
	Token       string `form:"token,optional"`
}

type FeedRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	NextTime   int64       `json:"next_time,optional"`
	VideoList  interface{} `json:"video_list,optional"`
}

type PublishReq struct {
	Token string `form:"token"`
	Title string `form:"title"`
}

type PublishRes struct {
	StatusCode uint32 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetPublishListReq struct {
	Token  string `form:"token"`
	UserId string `form:"user_id"`
}

type GetPublishListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,optional"`
}

type FavoriteReq struct {
	Token      string `form:"token"`
	VideoId    string `form:"video_id"`
	ActionType string `form:"action_type"`
}

type FavoriteRes struct {
	StatusCode uint32 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetFavoriteListReq struct {
	UserId string `form:"user_id"`
	Token  string `form:"token"`
}

type GetFavoriteListRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,optional"`
}

type CommentReq struct {
	Token       string `form:"token"`
	VideoId     string `form:"video_id"`
	ActionType  string `form:"action_type"`
	CommentText string `form:"comment_text,optional"`
	CommentId   string `form:"comment_id,optional"`
}

type CommentRes struct {
	StatusCode uint32      `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	Comment    interface{} `json:"comment,optional"`
}

type GetCommentListReq struct {
	Token   string `form:"token"`
	VideoId string `form:"video_id"`
}

type GetCommentListRes struct {
	StatusCode  uint32      `json:"status_code"`
	StatusMsg   string      `json:"status_msg"`
	CommentList interface{} `json:"comment_list,optional"`
}
