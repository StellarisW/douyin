package types

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"strings"
)

// GetViper 获取 viper
func (c *Agollo) GetViper(namespace string) (*viper.Viper, error) {
	if v, ok := c.vipers[namespace]; ok {
		// 返回已经解析的命名空间的 viper
		return v, nil
	} else {
		// 命名空间的配置文件没有解析, 则生成一个新 viper 来解析
		v := viper.New()
		// 判断命名空间的文件类型(properties, yaml)
		namespaceType := getNamespaceType(namespace)
		// 设置 viper 解析的配置文件类型
		v.SetConfigType(namespaceType)

		// 缓存配置文件
		var buffer *bytes.Buffer
		switch namespaceType {
		case "properties":
			buffer = bytes.NewBufferString(c.client.GetConfig(namespace).GetContent())

		case "yaml":
			buffer = bytes.NewBufferString(c.client.GetConfig(namespace).GetValue("content"))

		default:
			return nil, errors.New("namespace type not supported")
		}

		// 读取缓存中的配置信息
		err := v.ReadConfig(buffer)
		if err != nil {
			return nil, err
		}

		return v, nil
	}
}

// getNamespaceType 获取配置命名空间的文件格式
func getNamespaceType(namespace string) string {
	output := strings.Split(namespace, ".")
	if len(output) == 1 {
		return "properties"
	}
	return output[1]
}
