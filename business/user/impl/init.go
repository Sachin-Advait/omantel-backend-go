package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/fcm"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/redis"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/soap"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

var (
	getQuarterFromMonth = objecthelper.GetQuarterFromMonth
)

type UserImpl struct {
	Config       *configs.UserConfig
	DbService    database.DbService
	CacheService redis.CacheService
	FCMService   fcm.FCMService
	SOAPService  soap.SOAPService
}

func GetUserService(config *configs.UserConfig, dbService database.DbService, cacheService redis.CacheService, fcmService fcm.FCMService,
	soapService soap.SOAPService) *UserImpl {

	return &UserImpl{
		Config:       config,
		DbService:    dbService,
		CacheService: cacheService,
		FCMService:   fcmService,
		SOAPService:  soapService,
	}
}

func (impl *UserImpl) SessionValidation(userId, currentAccessToken string) bool {

	session, err := impl.CacheService.GetUserSessionToken(userId)
	if err != nil && err.Code != errors.UserNotFoundCode {
		return false
	}

	if err != nil && err.Code == errors.UserNotFoundCode {
		return true
	}

	if session.AccessToken == currentAccessToken {
		return true
	}

	return false
}
