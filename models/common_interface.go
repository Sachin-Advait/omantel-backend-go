package models

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"time"
)

type IncentiveProvider interface {
	HandleDataInputFile()
	HandleUserMasterInputFile()
	HandlePlansInputFile()
	HandleTargetInputFile()
	HandleScaleOrderInputFile()
	HandleMaxPayInputFile()
	HandleHierarchyFile()

	GetSingleUserIncentive(ctx *gin.Context, id, quarter int, isQuarterReport bool, date time.Time) (*psqldb.TotalIncentive, *errors.Error)
	GetSingleUserTargetChart(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) (*psqldb.ChartDataResponse, *errors.Error)
	GetSaleCountPerPlan(ctx *gin.Context, clientID, quarter int, product string, isQuarterReport bool, t time.Time) ([]*psqldb.SaleCountPerPlan, *errors.Error)
	GenerateOTP(ctx *gin.Context, id int) (map[string]string, *errors.Error)
	VerifyOTP(ctx *gin.Context, clientID, otp int) (*VerifyOTPResponse, *errors.Error)
	GetUserReport(ctx *gin.Context, clientID, month, quarter int, isQuarterReport bool) (*psqldb.UserRecord, *errors.Error)
	GetUsersReportingToManager(ctx *gin.Context, managerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error)
	GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error)
	GetUserPayout(ctx *gin.Context, id, quarter int, isQuarterly bool, t time.Time) (*psqldb.TotalIncentive, *errors.Error)
	GetUserTarget(ctx *gin.Context, id, quarter int, isQuarterly bool, t time.Time) (*psqldb.ChartDataResponse, *errors.Error)
	GetOutletTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesCountPerOutlet, *errors.Error)
	GetOutletSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetRegionTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesCountPerOutlet, *errors.Error)
	GetRegionSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesTargetPerOutlet, *errors.Error)
	GetTopTenEmpProductSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesCountPerStaff, *errors.Error)
	GetTopTenEmpSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesTargetPerStaff, *errors.Error)
}
