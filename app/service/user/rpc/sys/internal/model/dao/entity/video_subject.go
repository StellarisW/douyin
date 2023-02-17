package entity

import (
	"time"
)

const TableNameVideoSubject = "video_subject"

// VideoSubject mapped from table <video_subject>
type VideoSubject struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`             // 主键id
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`     // 作者id
	PlayURL    string    `gorm:"column:play_url;not null" json:"play_url"`   // 视频播放地址
	CoverURL   string    `gorm:"column:cover_url;not null" json:"cover_url"` // 视频封面地址
	Title      string    `gorm:"column:title;not null" json:"title"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;not null" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime;not null" json:"update_time"` // 修改时间
}

// TableName VideoSubject's table name
func (*VideoSubject) TableName() string {
	return TableNameVideoSubject
}
