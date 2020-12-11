package config

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"time"
)
type RedisConfig struct {
	// host
	Host string

	// port, default 6379
	Port int

	// auth
	Auth string

	// pool size, default 100
	PoolSize int

	// max retries, default 3
	MaxRetries int

	// timeout, default 10 seconds
	IdleTimeout time.Duration
}

// Options get redis options
func (conf *RedisConfig) Options() *redis.Options {
	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	return &redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}


// Redis return redis config
func Redis() (options *RedisConfig) {
	if err := container.Invoke(func(o *RedisConfig) {
		options = o
	}); err != nil {
		options = &RedisConfig{
			Host:        viper.GetString(redisHostKey),
			Port:        viper.GetInt(redisPortKey),
			Auth:        viper.GetString(redisPassKey),
			PoolSize:    viper.GetInt(redisPoolSizeKey),
			MaxRetries:  viper.GetInt(redisMaxRetriesKey),
			IdleTimeout: time.Duration(viper.GetInt(redisIdleTimeoutKey)) * time.Second,
		}

		_ = container.Provide(func() *RedisConfig {
			return options
		})
	}
	return
}
