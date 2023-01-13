package config

import (
	"douyin/app/common/config/internal/common"
	"douyin/app/common/config/internal/database"
)

var insDatabase = database.Group{Agollo: agolloClient}
var insCommon = common.Group{Agollo: agolloClient}

func Database() *database.Group {
	return &insDatabase
}

func Common() *common.Group {
	return &insCommon
}
