package config

import (
	"douyin/app/common/config/internal/consts"
	"douyin/app/common/config/internal/types"
	"douyin/app/common/log"
	"github.com/apolloconfig/agollo/v4"
	"go.uber.org/zap"
	"os"

	"github.com/apolloconfig/agollo/v4/constant"
	agolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
	"github.com/spf13/viper"
)

var agolloClient *types.Agollo // 声明一个 Agollo 对象, 用于在当前包范围内的函数操作调用

// InitClient 初始化 Agollo 客户端
func InitClient() (err error) {
	connConfig := types.AgolloConnConfig{
		AppID:       os.Getenv("APOLLO_APP_ID"),
		ClusterName: os.Getenv("APOLLO_CLUSTER_NAME"),
		IP:          os.Getenv("APOLLO_IP"),
		Secret:      os.Getenv("APOLLO_SECRET"),
	}
	log.Logger.Debug("", zap.Reflect("connConfig", connConfig))
	vipers := make(map[string]*viper.Viper)
	appConfig := &agolloConfig.AppConfig{
		AppID:         connConfig.AppID,
		Cluster:       connConfig.ClusterName,
		IP:            connConfig.IP,
		NamespaceName: consts.Namespaces,
		Secret:        connConfig.Secret,
		MustStart:     true,
	}
	// 客户端不解析 content, viper 来解析
	extension.AddFormatParser(constant.Properties, &emptyParser{})
	extension.AddFormatParser(constant.YAML, &emptyParser{})

	// 设置 apollo 的日志器
	agollo.SetLogger(log.Logger.Sugar())

	// 设置 apollo 的缓存组件
	agollo.SetCache(&DefaultCacheFactory{})

	client, err := agollo.StartWithConfig(func() (*agolloConfig.AppConfig, error) {
		return appConfig, nil
	})
	if err != nil {
		return err
	}

	// 设置 配置监听功能
	//client.AddChangeListener(&CustomChangeListener{})

	agolloClient = types.NewAgolloClient(client, vipers)

	return nil
}
