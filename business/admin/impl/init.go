package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/fcm"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/redis"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

var (
	getQuarterFromMonth = objecthelper.GetQuarterFromMonth
)

type AdminImpl struct {
	Config       *configs.AdminConfig
	DbService    database.DbService
	CacheService redis.CacheService
	FCMService   fcm.FCMService
}

func GetAdminService(config *configs.AdminConfig, dbService database.DbService, cacheService redis.CacheService, fcmService fcm.FCMService) *AdminImpl {

	return &AdminImpl{
		Config:       config,
		DbService:    dbService,
		CacheService: cacheService,
		FCMService:   fcmService,
	}
}
