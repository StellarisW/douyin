package types

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/spf13/viper"
)

// Agollo 配置总管理对象
type Agollo struct {
	client agollo.Client
	vipers map[string]*viper.Viper
}

// AgolloConnConfig 连接配置
type AgolloConnConfig struct {
	ClusterName string // apollo 集群类型
	IP          string // apollo 集群 IP
	AppID       string // apollo app ID
	Secret      string // apollo app secret
}

var AgolloClient *Agollo // 声明一个 Agollo 对象, 用于在当前包范围内的函数操作调用

func NewAgolloClient(client agollo.Client, vipers map[string]*viper.Viper) *Agollo {
	return &Agollo{
		client: client,
		vipers: vipers,
	}
}
