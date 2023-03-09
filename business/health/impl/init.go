package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/redis"
)

type HealthImpl struct {
	DbService    database.DbService
	CacheService redis.CacheService
}

func GetHealthService(dbService database.DbService, cacheService redis.CacheService) *HealthImpl {

	return &HealthImpl{
		DbService:    dbService,
		CacheService: cacheService,
	}
}
