package business

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetRegionSalesAndTarget(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID)
	if err != nil {
		return nil, err
	}

	mangerPosition, ok := constant.ManagerMap[user.Position]
	if !ok || user.Position == constant.Promoter || user.Position == constant.CSR || user.Position == constant.OutletManager || user.Position == constant.TeamLeader {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "INCORRECT USER POSITION",
			Details:    "incorrect position",
		}
	}

	var needUpdate bool

	currTime := time.Now()
	currYear, currMonth, currDay := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if isQuarterReport && year == currYear && quarter <= currQuarter {
		needUpdate = true
	} else if !isQuarterReport && year == currYear && month <= int(currMonth) {
		needUpdate = true
	}

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	// for regional manager
	if user.Position == constant.Manager {

		regionReport, err := impl.PsqlDb.FetchRegionReport(user.Region, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && regionReport.CreatedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			regionReport, err = impl.updateRegionReport(ctx, user.Region, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
		}

		eachRegionResp := &models.ProductWiseSalesTargetPerOutlet{
			Region:        regionReport.Region,
			HBBSales:      regionReport.HBBSalesCount,
			HayyakSales:   regionReport.HayyakSalesCount,
			BaqatiSales:   regionReport.BaqatiSalesCount,
			UpgradeSales:  regionReport.UpgradeSalesCount,
			HBBTarget:     regionReport.HBBTargetCount,
			HayyakTarget:  regionReport.HayyakTargetCount,
			BaqatiTarget:  regionReport.BaqatiTargetCount,
			UpgradeTarget: regionReport.UpgradeTargetCount,
		}

		eachRegionResp.TotalSales = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount
		eachRegionResp.TotalTarget = regionReport.HBBTargetCount + regionReport.HayyakTargetCount + regionReport.BaqatiTargetCount + regionReport.UpgradeTargetCount

		if eachRegionResp.HBBTarget != 0 {
			eachRegionResp.HBBAchievedPercentage = math.Round(float64(eachRegionResp.HBBSales*100/eachRegionResp.HBBTarget)*100) / 100
		}

		if eachRegionResp.HayyakTarget != 0 {
			eachRegionResp.HayyakAchievedPercentage = math.Round(float64(eachRegionResp.HayyakSales*100/eachRegionResp.HayyakTarget)*100) / 100
		}

		if eachRegionResp.BaqatiTarget != 0 {
			eachRegionResp.BaqatiAchievedPercentage = math.Round(float64(eachRegionResp.BaqatiSales*100/eachRegionResp.BaqatiTarget)*100) / 100
		}

		if eachRegionResp.UpgradeTarget != 0 {
			eachRegionResp.UpgradeAchievedPercentage = math.Round(float64(eachRegionResp.UpgradeSales*100/eachRegionResp.UpgradeTarget)*100) / 100
		}

		if eachRegionResp.TotalTarget != 0 {
			eachRegionResp.TotalAchievedPercentage = math.Round(float64(eachRegionResp.TotalSales*100/eachRegionResp.TotalTarget)*100) / 100
		}

		resp = append(resp, eachRegionResp)

		return resp, nil
	}

	// for non region manager
	allDistinctMngr, err := impl.PsqlDb.FetchAllRMReportingToGivenManager(mangerPosition, userID)
	if err != nil {
		return nil, err
	}

	if len(allDistinctMngr) == 0 {
		logger.FileLogger.CtxError(ctx).Str("position", mangerPosition).Msg("GetRegionSalesAndTarget : no OM found")
		return resp, nil
	}

	for _, eachMngr := range allDistinctMngr {

		regionReport, err := impl.PsqlDb.FetchRegionReport(eachMngr.Outlet, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && regionReport.CreatedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			regionReport, err = impl.updateRegionReport(ctx, eachMngr.Region, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
		}

		eachRegionResp := &models.ProductWiseSalesTargetPerOutlet{
			Region:        regionReport.Region,
			HBBSales:      regionReport.HBBSalesCount,
			HayyakSales:   regionReport.HayyakSalesCount,
			BaqatiSales:   regionReport.BaqatiSalesCount,
			UpgradeSales:  regionReport.UpgradeSalesCount,
			HBBTarget:     regionReport.HBBTargetCount,
			HayyakTarget:  regionReport.HayyakTargetCount,
			BaqatiTarget:  regionReport.BaqatiTargetCount,
			UpgradeTarget: regionReport.UpgradeTargetCount,
		}

		eachRegionResp.TotalSales = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount
		eachRegionResp.TotalTarget = regionReport.HBBTargetCount + regionReport.HayyakTargetCount + regionReport.BaqatiTargetCount + regionReport.UpgradeTargetCount

		if eachRegionResp.HBBTarget != 0 {
			eachRegionResp.HBBAchievedPercentage = math.Round(float64(eachRegionResp.HBBSales*100/eachRegionResp.HBBTarget)*100) / 100
		}

		if eachRegionResp.HayyakTarget != 0 {
			eachRegionResp.HayyakAchievedPercentage = math.Round(float64(eachRegionResp.HayyakSales*100/eachRegionResp.HayyakTarget)*100) / 100
		}

		if eachRegionResp.BaqatiTarget != 0 {
			eachRegionResp.BaqatiAchievedPercentage = math.Round(float64(eachRegionResp.BaqatiSales*100/eachRegionResp.BaqatiTarget)*100) / 100
		}

		if eachRegionResp.UpgradeTarget != 0 {
			eachRegionResp.UpgradeAchievedPercentage = math.Round(float64(eachRegionResp.UpgradeSales*100/eachRegionResp.UpgradeTarget)*100) / 100
		}

		if eachRegionResp.TotalTarget != 0 {
			eachRegionResp.TotalAchievedPercentage = math.Round(float64(eachRegionResp.TotalSales*100/eachRegionResp.TotalTarget)*100) / 100
		}

		resp = append(resp, eachRegionResp)
	}

	return resp, nil
}
