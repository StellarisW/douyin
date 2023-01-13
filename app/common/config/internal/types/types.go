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

func NewAgolloClient(client agollo.Client, vipers map[string]*viper.Viper) *Agollo {
	return &Agollo{
		client: client,
		vipers: vipers,
	}
}
