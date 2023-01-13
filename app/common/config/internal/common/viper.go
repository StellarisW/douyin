package common

import "github.com/mitchellh/mapstructure"

// UnmarshalServiceConfig 将微服务的配置信息反序列化
func (g *Group) UnmarshalServiceConfig(namespace, serviceType, serviceName string, dst interface{}) (err error) {
	if serviceType == "api" || serviceType == "mq" {
		return g.UnmarshalKey(namespace, serviceType, dst)
	}

	return g.UnmarshalKey(namespace, serviceType+"."+serviceName, dst)
}

// UnmarshalKey 将对应 key 的配置信息反序列化到指定变量
func (g *Group) UnmarshalKey(namespace, key string, dst interface{}) (err error) {
	v, err := g.agollo.GetViper(namespace)
	if err != nil {
		return
	}

	err = v.UnmarshalKey(key, dst, func(decoderConfig *mapstructure.DecoderConfig) {
		// 允许解析匿名字段
		decoderConfig.Squash = true
	})

	return
}
