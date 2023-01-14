package database

import (
	"douyin/app/common/config/internal/common"
	"douyin/app/common/config/internal/consts"
	"fmt"
	"github.com/yitter/idgenerator-go/idgen"
)

// NewIdGenerator 创建新的ID生成器
func (g *Group) NewIdGenerator(namespace string) (idGenerator *idgen.DefaultIdGenerator, err error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	options, err := g.newIdGeneratorOptions(namespace)
	if err != nil {
		return nil, fmt.Errorf("get options failed, %v", err)
	}

	idGenerator = idgen.NewDefaultIdGenerator(options)

	return idGenerator, nil
}

func (g *Group) newIdGeneratorOptions(namespace string) (options *idgen.IdGeneratorOptions, err error) {
	options = &idgen.IdGeneratorOptions{}
	err = common.GetGroup().UnmarshalKey(namespace, "IdGeneratorOptions", options)
	if err != nil {
		return nil, err
	}

	return options, nil
}
