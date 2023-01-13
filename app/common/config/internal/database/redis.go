package database

import (
	"context"
	"douyin/app/common/config/internal/consts"
	"fmt"
	"github.com/go-redis/redis/v9"
	"net"
	"time"
)

// GetRedisClusterClient 获取redis集群客户端实例
func (g *Group) GetRedisClusterClient() (rdb *redis.ClusterClient, err error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	redisOptions, err := g.NewRedisClusterOptions()
	if err != nil {
		return nil, err
	}

	rdb = redis.NewClusterClient(redisOptions)
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

// NewRedisClusterOptions 返回 *redis.ClusterOptions
func (g *Group) NewRedisClusterOptions() (*redis.ClusterOptions, error) {
	v, err := g.agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return nil, err
	}

	return &redis.ClusterOptions{
		Addrs:    v.GetStringSlice("Database.Redis.Addrs"),
		Username: "",
		Password: v.GetString("Database.Redis.Password"),

		NewClient: nil,

		ReadOnly:       false,
		RouteByLatency: false,
		RouteRandomly:  false,

		ClusterSlots: nil,

		// 可自定义连接函数
		Dialer: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial(network, addr)
		},

		// 钩子函数
		// 仅当客户端执行命令时需要从连接池获取连接时, 如果连接池需要新建连接时则会调用此钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error {
			fmt.Printf("conn=%v\n", conn)
			return nil
		},

		MaxRedirects: 3,

		// 命令执行失败时的重试策略
		MaxRetries:      3,                      // 命令执行失败时, 最多重试多少次, 默认为0即不错红石
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限, 默认8毫秒, -1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上线, 默认512毫秒, -1表示取消间隔

		// 超时
		DialTimeout:  5 * time.Second, // 连接建立超时时间, 默认5秒
		ReadTimeout:  3 * time.Second, // 读超时, 默认3秒, -1表示取消读超时
		WriteTimeout: 3 * time.Second, // 写超时, 默认等于读超时
		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时, 客户端等待可用连接的最大等待时长, 默认为读超时+1秒

		PoolFIFO: false,
		PoolSize: 16, // 连接池最大socket连接数, 默认为4倍CPU数, 4*runtime.NumCPU

		MinIdleConns:    8, // 在启动阶段创建指定数量的Idle连接, 并长期维持idle状态的连接数不少于指定数量
		MaxIdleConns:    30,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 5 * time.Minute,

		ContextTimeoutEnabled: true,

		TLSConfig: nil,
	}, nil
}
