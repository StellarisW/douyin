package entity

import (
	"time"
)

const TableNameChatMessage = "chat_message"

// ChatMessage mapped from table <chat_message>
type ChatMessage struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`                                // 主键id
	SrcUserID  int64     `gorm:"column:src_user_id;not null" json:"src_user_id"`                // 发送者id
	DstUserID  int64     `gorm:"column:dst_user_id;not null" json:"dst_user_id"`                // 接收者id
	Content    string    `gorm:"column:content" json:"content"`                                 // 消息内容
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;not null" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime;not null" json:"update_time"` // 修改时间
}

// TableName ChatMessage's table name
func (*ChatMessage) TableName() string {
	return TableNameChatMessage
}
