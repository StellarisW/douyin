package entity

import (
	"time"
)

const TableNameUserSubject = "user_subject"

// UserSubject mapped from table <user_subject>
type UserSubject struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`                                // 主键id
	Username   string    `gorm:"column:username;not null" json:"username"`                      // 用户名
	Password   string    `gorm:"column:password;not null" json:"password"`                      // 密码
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;not null" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime;not null" json:"update_time"` // 修改时间
}

// TableName UserSubject's table name
func (*UserSubject) TableName() string {
	return TableNameUserSubject
}
