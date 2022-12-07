package models

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type IncentiveProvider interface {
	LoadPlansFromDB(printLogs bool) error
	LoadTargetFromDB(printLogs bool) error
	LoadScaleOrderFromDB(printLogs bool) error
	LoadMaxPayFromDB(printLogs bool) error

	GenerateOTP(ctx *gin.Context, id string) (map[string]string, *errors.Error)
	VerifyOTP(ctx *gin.Context, clientID string, otp int) (*VerifyOTPResponse, *errors.Error)
	CaptchaVerification(ctx *gin.Context, req *CaptchRequest) (interface{}, *errors.Error)

	GetSingleUserIncentive(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*TotalIncentive, *errors.Error)
	GetSingleUserTargetChart(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*ChartDataResponse, *errors.Error)
	GetSaleCountPerPlan(ctx *gin.Context, userId, product string, isQuarterReport bool, quarter, month, year int) ([]*SaleCountPerPlan, *errors.Error)
	GetUserReport(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*UserRecord, *errors.Error)
	GetUsersReportingToManager(ctx *gin.Context, managerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUserPayout(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*TotalIncentive, *errors.Error)
	GetUserTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*ChartDataResponse, *errors.Error)
	GetOutletTotalSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetOutletSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetRegionTotalSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetRegionSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetTopTenEmpProductSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesCountPerStaff, *errors.Error)
	GetTopTenEmpSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*ProductWiseSalesTargetPerStaff, *errors.Error)
	RenewAccessToken(ctx *gin.Context) (string, *errors.Error)

	FetchRegions(ctx *gin.Context, req *DeleteRegionRequest) ([]*GetRegionResponse, *errors.Error)
	CreateNewRegion(ctx *gin.Context, req []*NewRegionRequest) ([]*NewRegionFailedResponse, *errors.Error)
	DeleteRegions(ctx *gin.Context, req *DeleteRegionRequest) *errors.Error

	FetchOutlets(ctx *gin.Context, req *DeleteOutletRequest) ([]*GetOutletResponse, *errors.Error)
	CreateNewOutlet(ctx *gin.Context, req []*NewOutletRequest) ([]*NewOutletFailedResponse, *errors.Error)
	DeleteOutlets(ctx *gin.Context, req *DeleteOutletRequest) *errors.Error

	FetchProducts(ctx *gin.Context, req *DeleteProductRequest) ([]*GetProductResponse, *errors.Error)
	CreateNewProduct(ctx *gin.Context, req []*NewProductRequest) ([]*NewProductFailedResponse, *errors.Error)
	DeleteProducts(ctx *gin.Context, req *DeleteProductRequest) *errors.Error

	FetchPlans(ctx *gin.Context, req *DeletePlanRequest) ([]*GetPlanResponse, *errors.Error)
	CreateNewPlan(ctx *gin.Context, req []*NewPlanRequest) ([]*NewPlanFailedResponse, *errors.Error)
	DeletePlans(ctx *gin.Context, req *DeletePlanRequest) *errors.Error

	FetchTargets(ctx *gin.Context, req *GetTargetRequest) (map[string][]*GetTargetResponse, *errors.Error)
	CreateNewTarget(ctx *gin.Context, req []*NewTargetRequest) ([]*NewTargetFailedResponse, *errors.Error)
	DeleteTargets(ctx *gin.Context, req *DeleteTargetRequest) *errors.Error
}
