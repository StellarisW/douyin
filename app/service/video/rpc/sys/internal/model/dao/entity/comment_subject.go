package entity

import (
	"time"
)

const TableNameCommentSubject = "comment_subject"

// CommentSubject mapped from table <comment_subject>
type CommentSubject struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"`                                // 主键id
	UserID      int64     `gorm:"column:user_id;not null" json:"user_id"`                        // 评论者id
	VideoID     int64     `gorm:"column:video_id;not null" json:"video_id"`                      // 评论的视频id
	CommentText string    `gorm:"column:comment_text" json:"comment_text"`                       // 评论内容
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime;not null" json:"create_time"` // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime;not null" json:"update_time"` // 修改时间
}

// TableName CommentSubject's table name
func (*CommentSubject) TableName() string {
	return TableNameCommentSubject
}
