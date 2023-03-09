package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/redis"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/cache"
)

type CacheServiceImpl struct {
	RedisClient cache.IncentiveCache
}

func GetCacheServiceInstance(redisClient cache.IncentiveCache) redis.CacheService {
	return &CacheServiceImpl{
		RedisClient: redisClient,
	}
}
