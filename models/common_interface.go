package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type IncentiveProvider interface {
	LoadPlansFromDB(printLogs bool) error
	LoadTargetFromDB(printLogs bool) error
	LoadScaleOrderFromDB(printLogs bool) error
	LoadMaxPayFromDB(printLogs bool) error

	GetSingleUserIncentive(ctx *gin.Context, id string, quarter int, isQuarterReport bool, date time.Time) (*TotalIncentive, *errors.Error)
	GetSingleUserTargetChart(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) (*ChartDataResponse, *errors.Error)
	GetSaleCountPerPlan(ctx *gin.Context, clientID string, quarter int, product string, isQuarterReport bool, t time.Time) ([]*SaleCountPerPlan, *errors.Error)
	GenerateOTP(ctx *gin.Context, id string) (map[string]string, *errors.Error)
	VerifyOTP(ctx *gin.Context, clientID string, otp int) (*VerifyOTPResponse, *errors.Error)
	GetUserReport(ctx *gin.Context, clientID string, month, quarter int, isQuarterReport bool) (*UserRecord, *errors.Error)
	GetUsersReportingToManager(ctx *gin.Context, managerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUserPayout(ctx *gin.Context, id string, quarter int, isQuarterly bool, t time.Time) (*TotalIncentive, *errors.Error)
	GetUserTarget(ctx *gin.Context, id string, quarter int, isQuarterly bool, t time.Time) (*ChartDataResponse, *errors.Error)
	GetOutletTotalSalesCount(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetOutletSalesAndTarget(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetRegionTotalSalesCount(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetRegionSalesAndTarget(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetTopTenEmpProductSalesCount(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerStaff, *errors.Error)
	GetTopTenEmpSalesAndTarget(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerStaff, *errors.Error)
	RenewAccessToken(ctx *gin.Context) (string, *errors.Error)

	CaptchaVerification(ctx *gin.Context, req *CaptchRequest) (interface{}, *errors.Error)

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
