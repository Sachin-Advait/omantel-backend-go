package cache

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

var printIncentiveCacheOnce sync.Once

type IncentiveCache interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	HGet(ctx context.Context, key string, field string) *redis.StringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Close() error
}

type RedisIncentiveCacheService struct {
	IncentiveRedis IncentiveCache
	ContextTimeOut time.Duration
	KeyTTL         time.Duration
}

var RedisIncentiveCacheServiceImpl = &RedisIncentiveCacheService{}

func NewRedisIncentiveCacheService(redisConfig *models.RedisConfigStruct) *RedisIncentiveCacheService {

	printIncentiveCacheOnce.Do(func() {
		logger.FileLogger.Info().Str("env", flags.Env()).
			Str("redisMode", redisConfig.Mode).
			Interface("address", redisConfig.Address).
			Str("password", redisConfig.Password).
			Int("db", redisConfig.Db).
			Int("poolSize", redisConfig.PoolSize).
			Int("minIdleConn", redisConfig.MinIdleConn).
			Dur("maxConnAgeInHour", redisConfig.MaxConnAgeInHour).
			Dur("poolTimeoutInMilliSecond", redisConfig.PoolTimeoutInMilliSecond).
			Dur("readTimeoutInMilliSecond", redisConfig.ReadTimeoutInMilliSecond).
			Dur("writeTimeoutInMilliSecond", redisConfig.WriteTimeoutInMilliSecond).
			Dur("contextTimeoutInMilliSecond", redisConfig.ContextTimeoutInMilliSecond).
			Msg("Incentive Cache Configuration")
	})
	mode := redisConfig.Mode
	var redisIncentiveCache IncentiveCache

	switch mode {
	case constant.Cluster:
		redisIncentiveCache = redis.NewClusterClient(&redis.ClusterOptions{
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
		redisIncentiveCache = redis.NewClient(&redis.Options{
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

	if redisIncentiveCache == nil {
		logger.FileLogger.Fatal().Str("env", flags.Env()).Stack().Msg("IncentiveCacheService - incorrect redis mode " + mode)
	}
	ctx, _ := context.WithTimeout(context.Background(), redisConfig.ContextTimeoutInMilliSecond)

	_, err := redisIncentiveCache.Ping(ctx).Result()
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Str("env", flags.Env()).Stack().Msg("IncentiveCacheService - ping error")
	}

	return &RedisIncentiveCacheService{
		IncentiveRedis: redisIncentiveCache,
		ContextTimeOut: redisConfig.ContextTimeoutInMilliSecond,
		KeyTTL:         redisConfig.KeyTTLInHour,
	}
}
