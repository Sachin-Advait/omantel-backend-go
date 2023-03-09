package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

type IncentiveCache interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

var client IncentiveCache

// InitCache is used to initialise the cache
func InitCache(redisConfig *models.RedisConfigStruct) {

	mode := redisConfig.Mode

	switch mode {
	case constant.Cluster:
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        []string{redisConfig.Address},
			Password:     redisConfig.Password,
			PoolSize:     redisConfig.PoolSize,
			MinIdleConns: redisConfig.MinIdleConn,
			MaxConnAge:   redisConfig.MaxConnAgeInHour,
			PoolTimeout:  redisConfig.PoolTimeoutInMilliSecond,
			ReadTimeout:  redisConfig.ReadTimeoutInMilliSecond,
			WriteTimeout: redisConfig.WriteTimeoutInMilliSecond,
		})
	case constant.NonCluster:
		client = redis.NewClient(&redis.Options{
			Addr:         redisConfig.Address,
			Password:     redisConfig.Password,
			PoolSize:     redisConfig.PoolSize,
			MinIdleConns: redisConfig.MinIdleConn,
			MaxConnAge:   redisConfig.MaxConnAgeInHour,
			PoolTimeout:  redisConfig.PoolTimeoutInMilliSecond,
			ReadTimeout:  redisConfig.ReadTimeoutInMilliSecond,
			WriteTimeout: redisConfig.WriteTimeoutInMilliSecond,
			DB:           redisConfig.Db,
		})
	}

	if client == nil {
		logger.FileLogger.Fatal().Str("env", flags.Env()).Stack().Msg("IncentiveCacheService - incorrect redis mode " + mode)
	}
	ctx, _ := context.WithTimeout(context.Background(), redisConfig.ContextTimeoutInMilliSecond)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Str("env", flags.Env()).Stack().Msg("IncentiveCacheService - ping error")
	}
}

// Get is used to get the cache client instance
func Get() IncentiveCache {
	return client
}

// Close is used to close the cache client
func Close() error {
	return client.Close()
}
