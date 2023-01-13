package middleware

import (
	"douyin/app/common/config/internal/consts"
	"douyin/app/common/middleware"
)

func (g *Group) NewCORSMiddleware() (*middleware.CORSMiddleware, error) {
	mode, err := g.getCORSMode()
	if err != nil {
		return nil, err
	}

	list, err := g.getCORSList()
	if err != nil {
		return nil, err
	}

	return middleware.NewCORSMiddleware(mode, list), nil
}

func (g *Group) getCORSMode() (string, error) {
	if g.agollo == nil {
		return "", consts.ErrEmptyConfigClient
	}

	v, err := g.agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return "", consts.ErrGetViper
	}

	mode := v.GetString("Middleware.CORS.Mode")
	return mode, nil
}

func (g *Group) getCORSList() (map[string]*middleware.CORSHeader, error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	v, err := g.agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return nil, consts.ErrGetViper
	}

	list := make(map[string]*middleware.CORSHeader, 0)
	rawList := v.GetStringMap("Middleware.CORS.List")
	for k := range rawList {
		list[k] = &middleware.CORSHeader{}
		err = v.UnmarshalKey("Middleware.CORS.List."+k, list[k])
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
