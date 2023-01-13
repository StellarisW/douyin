package common

import (
	"douyin/app/common/config/internal/consts"
	"github.com/spf13/viper"
)

func (g *Group) GetViper(namespace string) (*viper.Viper, error) {
	return g.Agollo.GetViper(namespace)
}

func (g *Group) GetDomain() (string, error) {
	if g.Agollo == nil {
		return "", consts.ErrEmptyConfigClient
	}

	v, err := g.Agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return "", consts.ErrGetViper
	}

	return v.GetString("Domain"), nil
}
