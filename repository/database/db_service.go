package database

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type DbService interface {
	UpsertRegionRecord(userId string, dbReq []*models.CreateRegionRequest) ([]*models.RegionFailedResponse, *errors.Error)
	UpsertOutletRecord(userId string, dbReq []*models.CreateOutletRequest) ([]*models.OutletFailedResponse, *errors.Error)
	UpsertProductRecord(userId string, dbReq []*models.CreateProductRequest) ([]*models.ProductFailedResponse, *errors.Error)
	UpsertPlanRecord(userId string, dbReq []*models.CreatePlanRequest) ([]*models.PlanFailedResponse, *errors.Error)
	UpsertTargetRecord(userId string, dbReq []*models.CreateTargetRequest) ([]*models.TargetFailedResponse, *errors.Error)
	UpsertVOCCriteriaRecord(userId string, dbReq []*models.CreateVOCCriteriaRequest) ([]*models.VOCCriteriaFailedResponse, *errors.Error)
	UpsertVOCMaxPayRecord(userId string, dbReq []*models.CreateVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error)
	UpsertVOCUserRecord(userId string, dbReq []*models.CreateVOCUserRequest) ([]*models.VOCCUserFailedResponse, *errors.Error)
	FetchRegionRecord(regionCodes []string, quarter, year int) ([]*models.RegionResponse, *errors.Error)
	FetchOutletRecord(outletCodes []string, quarter, year int) ([]*models.OutletResponse, *errors.Error)
	FetchProductRecord(productCodes []string, quarter, year int) ([]*models.ProductResponse, *errors.Error)
	FetchPlanRecord(planCodes []string, quarter, year int) ([]*models.PlanResponse, *errors.Error)
	FetchTargetRecord(targetCodes []string, outletCodes []string, month, quarter, year int) (map[string][]*models.TargetResponse, *errors.Error)
	FetchVOCCriteriaRecord(year int) (*models.VOCCriteriaResponse, *errors.Error)
	FetchVOCMaxPayRecord(year int) ([]*models.VOCMaxPayResponse, *errors.Error)
	FetchVOCUserRecord(userIds []string, year int) ([]*models.VOCUserResponse, *errors.Error)
	DeleteRegionRecord(userId string, regionCodes []string, quarter, year int) *errors.Error
	DeleteOutletRecord(userId string, outletCode []string, quarter, year int) *errors.Error
	DeleteProductRecord(userId string, productCode []string, quarter, year int) *errors.Error
	DeletePlanRecord(userId string, planCode []string, quarter, year int) *errors.Error
	DeleteTargetRecord(userId string, targetCodes []string, outletCodes []string, month, quarter, year int) *errors.Error
	DeleteVOCCriteriaRecord(userId string, year int) *errors.Error
	DeleteVOCMaxPayRecord(userId string, positions []string, year []int) ([]*models.VOCMaxPayFailedResponse, *errors.Error)
	DeleteVOCUserRecord(userId string, agentIds []string, year int) *errors.Error

	FetchPlans(year int) error
	FetchSaleTarget(year int) error
	FetchValueOfPoint(year int) error
	FetchMaxPayout(year int) error
	FetchMaxExtraPayout(year int) error
	FetchMaxThreshold(year int) error
	FetchVOCCriteria(year int) error
	FetchVOCMaxPay(year int) error

	UpsertNotificationRecord(userId string, dbReq *models.NotificationDB) (string, *errors.Error)
	FetchBasicNotificationRecord(category string) ([]*models.NotificationDB, *errors.Error)
	FetchNotificationRecord(category string) ([]*models.NotificationDB, *errors.Error)

	FetchOutletDetailsViaName(regionName, outletName string, quarter, year int) (*models.OutletResponse, *errors.Error)
	FetchAllOutletNames(regionName string, quarter, year int) ([]string, *errors.Error)
	FetchOutletSalesCount(region, outlet string, isQuarterReport bool, quarter, month, year int) (hbbSales, hayyakSales, baqatiSales, upgradeSales int, err *errors.Error)
	FetchActiveUserSales(userID string, region string, outlet string, isQuarterReport bool, quarter, month, year int) (*models.SaleTransaction, *errors.Error)
	FetchAllAgentsUserOfOM(outletManagerID string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error)
	FetchAllOMReportingToGivenManager(managerPosition string, managerID string, quarter, year int) (map[string]*models.UserBasicDetailsDB, *errors.Error)
	FetchAllMngrReportingToManagerRole(managerID, managerPosition string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error)
	FetchAllEmpReportingToGivenManagerRole(managerPosition string, managerID string, quarter, year int) ([]*models.ReportingStaffIDAndName, *errors.Error)
	FetchVP(quarter, year int) ([]string, *errors.Error)

	FetchUserBasicDetails(userID string, quarter, year int, args ...bool) (*models.UserBasicDetailsDB, *errors.Error)
	FetchUserIdFromOutletAndRegion(region, outlet string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error)
	FetchUserDetails(userID string, quarter, year int) (*models.UserBasicDetailsDB, *errors.Error)
	FetchOTPValidity(userID string) (*models.OtpModel, *errors.Error)
	InsertStaffOTP(update bool, otpModel *models.OtpModel) *errors.Error
	FetchTotalNonManagerialEmployees(region, outlet string, quarter, year int) (int, *errors.Error)
}
