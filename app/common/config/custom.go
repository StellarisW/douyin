package config

import (
	"douyin/app/common/log"
	"errors"
	"go.uber.org/zap"
	"sync"

	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/storage"
)

// 空解析器 (用viper解析)
type emptyParser struct {
}

func (d *emptyParser) Parse(configContent interface{}) (map[string]interface{}, error) {

	return nil, nil
}

// CustomChangeListener 自定义客户端配置监控器
type CustomChangeListener struct {
}

// DefaultCacheFactory 构造默认缓存组件工厂类
type DefaultCacheFactory struct {
}

// DefaultCache 默认缓存
type DefaultCache struct {
	defaultCache sync.Map
}

func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	//log.Logger.Info("config changed.", zap.String("namespace", changeEvent.Namespace), zap.Int64("NotificationID", changeEvent.NotificationID))
	//for k, v := range changeEvent.Changes {
	//	log.Logger.Info("", zap.String("key", k), zap.Reflect("ChangeType", v.ChangeType), zap.Reflect("OldValue", v.OldValue), zap.Reflect("NewValue", v.NewValue))
	//}
}

func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	log.Logger.Info("changed config(FullChangeEvent)", zap.String("namespace", event.Namespace), zap.Int64("NotificationID", event.NotificationID))
	for k, v := range event.Changes {
		log.Logger.Info("", zap.String("key", k), zap.Reflect("value", v))
	}
}

// Set 获取缓存
func (d *DefaultCache) Set(key string, value interface{}, expireSeconds int) (err error) {
	d.defaultCache.Store(key, value)
	return nil
}

// EntryCount 获取实体数量
func (d *DefaultCache) EntryCount() (entryCount int64) {
	count := int64(0)
	d.defaultCache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Create 创建默认缓存组件
func (d *DefaultCacheFactory) Create() agcache.CacheInterface {
	return &DefaultCache{}
}

// Get 获取缓存
func (d *DefaultCache) Get(key string) (value interface{}, err error) {
	v, ok := d.defaultCache.Load(key)
	if !ok {
		return nil, errors.New("load default cache fail")
	}
	return v.(string), nil
}

// Range 遍历缓存
func (d *DefaultCache) Range(f func(key, value interface{}) bool) {
	d.defaultCache.Range(f)
}

// Del 删除缓存
func (d *DefaultCache) Del(key string) (affected bool) {
	d.defaultCache.Delete(key)
	return true
}

// Clear 清除所有缓存
func (d *DefaultCache) Clear() {
	d.defaultCache = sync.Map{}
}
