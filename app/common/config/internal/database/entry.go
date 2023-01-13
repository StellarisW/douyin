package database

import (
	"douyin/app/common/config/internal/types"
)

type Group struct {
	agollo *types.Agollo
}

var insDatabase = &Group{}

func GetGroup() *Group {
	if insDatabase.agollo == nil {
		insDatabase = &Group{agollo: types.AgolloClient}
	}

	return insDatabase
}
