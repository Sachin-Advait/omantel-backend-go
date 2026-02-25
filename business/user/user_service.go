package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type UserProvider interface {
	FetchUserRegions(ctx *gin.Context) ([]string, *errors.Error)

	GenerateOTP(ctx *gin.Context, id string) (map[string]string, *errors.Error)
	GenerateOTPDev(ctx *gin.Context, userID string) (map[string]string, *errors.Error)
	VerifyOTP(ctx *gin.Context, req *models.VerifyOTP) (*models.VerifyOTPResponse, *errors.Error)
	GenerateToken(ctx *gin.Context) (string, *errors.Error)
	RenewAccessToken(ctx *gin.Context) (string, *errors.Error)
	CaptchaVerification(ctx *gin.Context, req *models.CaptchRequest) (interface{}, *errors.Error)
	UserLogout(ctx *gin.Context) *errors.Error

	GetUserFromRegionOrOutlet(ctx *gin.Context, outlet, region string, quarter,
		year int) ([]*models.UserMapResponse, *errors.Error)

	FetchNotification(ctx *gin.Context, category, id string) (interface{}, *errors.Error)

	GetUserPayout(ctx *gin.Context, userId string, chanel string, isQuarterReport bool, quarter, month,
		year int) (*models.TotalIncentive, *errors.Error)
	GetUserSalesAndTargetChart(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) (*models.ChartDataResponse, *errors.Error)
	GetUserSalesAndTargetChartPerPlan(ctx *gin.Context, userID string, chanel string,
		planType string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error)

	GetSaleCountPerPlan(ctx *gin.Context, userId, product string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.SaleCountPerPlan, *errors.Error)
	GetUserReport(ctx *gin.Context, userId string, chanel string, isQuarterReport bool, quarter, month,
		year int) (*models.UserReportDB, *errors.Error)

	GetUsersReportingToManager(ctx *gin.Context, managerID string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)

	GetOutletTotalSalesCount(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error)
	GetOutletSalesAndTarget(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetRegionTotalSalesCount(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) (map[string]*models.ProductWiseSalesCountPerOutlet, *errors.Error)
	GetRegionSalesAndTarget(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) (map[string]*models.ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetTopTenEmpProductSalesCount(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ProductWiseSalesCountPerStaff, *errors.Error)
	GetTopTenEmpSalesAndTarget(ctx *gin.Context, userId string, chanel string, isQuarterReport bool,
		quarter, month, year int) ([]*models.ProductWiseSalesTargetPerStaff, *errors.Error)

	GetUserSalesCount(ctx *gin.Context, userID, product string, chanel string, isQuarterReport bool,
		quarter, month, year int) (interface{}, *errors.Error)
	GetUserPlanCategory(ctx *gin.Context, userID, product string, chanel string, isQuarterReport bool,
		quarter, month, year int) (interface{}, *errors.Error)

	GetUserMNPSalesCount(ctx *gin.Context, userID string, chanel string, isQuarterReport bool,
		quarter, month, year int) (interface{}, *errors.Error)
	GetUserMNPPlanCount(ctx *gin.Context, userID string, chanel string, isQuarterReport bool,
		quarter, month, year int) (interface{}, *errors.Error)
}
