package middleware

import (
	"douyin/app/common/config/internal/common"
	"douyin/app/common/config/internal/consts"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"strings"
)

func (g *Group) NewSentinelMiddleware(namespace string) (*middleware.SentinelMiddleware, error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	conf, err := g.newSentinelConfig(namespace)
	if err != nil {
		return nil, err
	}

	rules, err := g.newFlowRules(namespace)
	if err != nil {
		return nil, err
	}

	return middleware.NewSentinelMiddleware(conf, rules), nil
}

func (g *Group) newSentinelConfig(namespace string) (*config.Entity, error) {
	//v, err := g.agollo.GetViper(consts.MainNamespace)
	//if err != nil {
	//	return nil, consts.ErrGetViper
	//}

	logger, err := log.GetSentinelLogger()
	if err != nil {
		return nil, err
	}

	return &config.Entity{
		Version: "v1",
		Sentinel: config.SentinelConfig{
			App: struct {
				Name string
				Type int32
			}{
				Name: "douyin-" + strings.Split(namespace, ".")[0] + "-api",
				Type: 0,
			},
			//Exporter: config.ExporterConfig{
			//	Metric: config.MetricExporterConfig{
			//		HttpAddr: "",
			//		HttpPath: "",
			//	},
			//},
			Log: config.LogConfig{
				Logger: logger,
				Dir:    "",
				UsePid: false,
				Metric: config.MetricLogConfig{
					SingleFileMaxSize: 1024 * 1024 * 50,
					MaxFileCount:      8,
					FlushIntervalSec:  1,
				},
			},
			Stat: config.StatConfig{
				GlobalStatisticSampleCountTotal: 20,
				GlobalStatisticIntervalMsTotal:  10000,
				MetricStatisticSampleCount:      2,
				MetricStatisticIntervalMs:       1000,
				System: config.SystemStatConfig{
					CollectIntervalMs:       1000,
					CollectLoadIntervalMs:   1000,
					CollectCpuIntervalMs:    1000,
					CollectMemoryIntervalMs: 150,
				},
			},
			UseCacheTime: true,
		},
	}, nil
}

func (g *Group) newFlowRules(namespace string) ([]*flow.Rule, error) {
	//v, err := g.agollo.GetViper(consts.MainNamespace)
	//if err != nil {
	//	return nil, consts.ErrGetViper
	//}

	var rules []*flow.Rule

	err := common.GetGroup().UnmarshalKey(namespace, "Api.Sentinel.Flow", &rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
