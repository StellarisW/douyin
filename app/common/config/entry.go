package config

import (
	"douyin/app/common/config/internal/common"
	"douyin/app/common/config/internal/database"
	"douyin/app/common/config/internal/middleware"
)

func Database() *database.Group {
	return database.GetGroup()
}
func Common() *common.Group {
	return common.GetGroup()
}
func Middleware() *middleware.Group {
	return middleware.GetGroup()
}
