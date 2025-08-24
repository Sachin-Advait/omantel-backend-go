package redis

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type CacheService interface {
	InsertUseSessionToken(userId string, userSession *models.UserSessionToken, keyTTL time.Duration) *errors.Error
	GetUserSessionToken(userId string) (*models.UserSessionToken, *errors.Error)
	DelUserSessionToken(userId string) *errors.Error
	IsSessionTokenValid(userId, sessionToken string) *errors.Error

	InsertUserReport(chanel string, isQuarterReport bool, month, quarter, year int,
		userReport *models.UserReportDB, keyTTL time.Duration) *errors.Error
	GetUserReport(userId string, chanel string, isQuarterReport bool, month, quarter,
		year int) (*models.UserReportDB, *errors.Error)
	DelUserReport(userId string, chanel string, isQuarterReport bool, month, quarter,
		year int) *errors.Error
	InsertOutletReport(chanel string, isQuarterReport bool, month, quarter, year int,
		outletReport *models.OutletReportDB, keyTTL time.Duration) *errors.Error
	GetOutletReport(outlet, region string, chanel string, isQuarterReport bool, month,
		quarter, year int) (*models.OutletReportDB, *errors.Error)
	DelOutletReport(outlet, region string, chanel string, isQuarterReport bool, month,
		quarter, year int) *errors.Error
	InsertRegionReport(chanel string, isQuarterReport bool, month, quarter, year int,
		regionReport *models.RegionReportDB, keyTTL time.Duration) *errors.Error
	GetRegionReport(region string, chanel string, isQuarterReport bool, month, quarter,
		year int) (*models.RegionReportDB, *errors.Error)
	DelRegionReport(region string, chanel string, isQuarterReport bool, month, quarter,
		year int) *errors.Error

	InsertUserSaleCountPerPlanReport(userId string, product string, chanel string, isQuarterReport bool,
		month, quarter, year int, userReport interface{}, keyTTL time.Duration) *errors.Error
	GetUserSaleCountPerPlanReport(userId string, product string, chanel string, isQuarterReport bool,
		month, quarter, year int) (interface{}, *errors.Error)

	InsertUserPlanCategoryReport(userId string, product string, chanel string, isQuarterReport bool,
		month, quarter, year int, userReport interface{}, keyTTL time.Duration) *errors.Error
	GetUserPlanCategoryReport(userId string, product string, chanel string, isQuarterReport bool,
		month, quarter, year int) (interface{}, *errors.Error)

	InsertUserMNPSaleCountPerPlanReport(userId string, chanel string, isQuarterReport bool, month,
		quarter, year int, userMNPReport interface{}, keyTTL time.Duration) *errors.Error
	GetUserMNPSaleCountPerPlanReport(userId string, chanel string, isQuarterReport bool, month, quarter,
		year int) (interface{}, *errors.Error)

	InsertUserMNPPlanCountPerPlanReport(userId string, chanel string, isQuarterReport bool, month,
		quarter, year int, userMNPReport interface{}, keyTTL time.Duration) *errors.Error
	GetUserMNPPlanCountPerPlanReport(userId string, chanel string, isQuarterReport bool, month, quarter,
		year int) (interface{}, *errors.Error)
}
