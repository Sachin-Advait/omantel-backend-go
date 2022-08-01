package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type IncentiveProvider interface {
	HandleDataInputFile()
	HandleUserMasterInputFile()
	HandlePlansInputFile()
	HandleTargetInputFile()
	HandleScaleOrderInputFile()
	HandleMaxPayInputFile()
	HandleHierarchyFile()

	GetSingleUserIncentive(ctx *gin.Context, id, quarter int, isQuarterReport bool, date time.Time) (*TotalIncentive, *errors.Error)
	GetSingleUserTargetChart(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) (*ChartDataResponse, *errors.Error)
	GetSaleCountPerPlan(ctx *gin.Context, clientID, quarter int, product string, isQuarterReport bool, t time.Time) ([]*SaleCountPerPlan, *errors.Error)
	GenerateOTP(ctx *gin.Context, id int) (map[string]string, *errors.Error)
	VerifyOTP(ctx *gin.Context, clientID, otp int) (*VerifyOTPResponse, *errors.Error)
	GetUserReport(ctx *gin.Context, clientID, month, quarter int, isQuarterReport bool) (*UserRecord, *errors.Error)
	GetUsersReportingToManager(ctx *gin.Context, managerID int) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID int) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID int) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID int) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID int) ([]*ReportingStaffIDAndName, *errors.Error)
	GetUserPayout(ctx *gin.Context, id, quarter int, isQuarterly bool, t time.Time) (*TotalIncentive, *errors.Error)
	GetUserTarget(ctx *gin.Context, id, quarter int, isQuarterly bool, t time.Time) (*ChartDataResponse, *errors.Error)
	GetOutletTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetOutletSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetRegionTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerOutlet, *errors.Error)
	GetRegionSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetTopTenEmpProductSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesCountPerStaff, *errors.Error)
	GetTopTenEmpSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*ProductWiseSalesTargetPerStaff, *errors.Error)

	CreateNewRegion(ctx *gin.Context, req *NewRegionRequest) *errors.Error
	DeleteRegions(ctx *gin.Context, req *DeleteRegionRequest) *errors.Error
	CreateNewOutlet(ctx *gin.Context, req *NewOutletRequest) *errors.Error
	DeleteOutlets(ctx *gin.Context, req *DeleteOutletRequest) *errors.Error
	CreateNewProduct(ctx *gin.Context, req *NewProductRequest) *errors.Error
	DeleteProducts(ctx *gin.Context, req *DeleteProductRequest) *errors.Error
	CreateNewPlan(ctx *gin.Context, req *NewPlanRequest) *errors.Error
	DeletePlans(ctx *gin.Context, req *DeletePlanRequest) *errors.Error
	CreateNewTarget(ctx *gin.Context, req *NewTargetRequest) *errors.Error
	DeleteTargets(ctx *gin.Context, req *DeleteTargetRequest) *errors.Error
}
