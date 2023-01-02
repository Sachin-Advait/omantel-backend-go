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

//GetOutletSalesAndTarget: userID should be a manager
func (impl *IncentiveImpl) GetOutletSalesAndTarget(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	_, ok := constant.ManagerMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT MANAGER POSITION - " + user.Position,
			Details:    "incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
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

	//TODO: remove at last
	needUpdate = true

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	// for oulletManager and teamLeader
	if user.Position == constant.OutletManager { //|| user.Position == constant.TeamLeader

		outletReport, err := impl.PsqlDb.FetchOutletReport(user.Outlet, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && outletReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				outletResp, err := impl.updateOutletReport(ctx, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Interface("OM outletReport", outletResp).Msg("GetOutletSalesAndTarget : trace")
				//TODO: test this out for struct assignment
				outletReport = outletResp
			}
		}

		eachOutletResp := &models.ProductWiseSalesTargetPerOutlet{
			OutletName:    outletReport.Outlet,
			Region:        outletReport.Region,
			HBBSales:      outletReport.HBBSalesCount,
			HayyakSales:   outletReport.HayyakSalesCount,
			BaqatiSales:   outletReport.BaqatiSalesCount,
			UpgradeSales:  outletReport.UpgradeSalesCount,
			HBBTarget:     outletReport.HBBTargetCount,
			HayyakTarget:  outletReport.HayyakTargetCount,
			BaqatiTarget:  outletReport.BaqatiTargetCount,
			UpgradeTarget: outletReport.UpgradeTargetCount,
		}

		eachOutletResp.TotalSales = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount
		eachOutletResp.TotalTarget = outletReport.HBBTargetCount + outletReport.HayyakTargetCount + outletReport.BaqatiTargetCount + outletReport.UpgradeTargetCount

		if eachOutletResp.HBBTarget != 0 {
			eachOutletResp.HBBAchievedPercentage = math.Round(float64(eachOutletResp.HBBSales*100/eachOutletResp.HBBTarget)*100) / 100
		}

		if eachOutletResp.HayyakTarget != 0 {
			eachOutletResp.HayyakAchievedPercentage = math.Round(float64(eachOutletResp.HayyakSales*100/eachOutletResp.HayyakTarget)*100) / 100
		}

		if eachOutletResp.BaqatiTarget != 0 {
			eachOutletResp.BaqatiAchievedPercentage = math.Round(float64(eachOutletResp.BaqatiSales*100/eachOutletResp.BaqatiTarget)*100) / 100
		}

		if eachOutletResp.UpgradeTarget != 0 {
			eachOutletResp.UpgradeAchievedPercentage = math.Round(float64(eachOutletResp.UpgradeSales*100/eachOutletResp.UpgradeTarget)*100) / 100
		}

		if eachOutletResp.TotalTarget != 0 {
			eachOutletResp.TotalAchievedPercentage = math.Round(float64(eachOutletResp.TotalSales*100/eachOutletResp.TotalTarget)*100) / 100
		}

		resp = append(resp, eachOutletResp)

		return resp, nil
	}

	// for non outlet and tl manager
	allDistinctOM, err := impl.PsqlDb.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("user", userID).Str("position", user.Position).Msg("GetOutletSalesAndTarget : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {

		outletReport, err := impl.PsqlDb.FetchOutletReport(eachOM.Outlet, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && outletReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				outletResp, err := impl.updateOutletReport(ctx, eachOM.Region, eachOM.Outlet, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr outletReport", outletResp).Msg("GetOutletSalesAndTarget : trace")
				//TODO: test this out for struct assignment
				outletReport = outletResp
			}
		}

		eachOutletResp := &models.ProductWiseSalesTargetPerOutlet{
			OutletName:    outletReport.Outlet,
			Region:        outletReport.Region,
			HBBSales:      outletReport.HBBSalesCount,
			HayyakSales:   outletReport.HayyakSalesCount,
			BaqatiSales:   outletReport.BaqatiSalesCount,
			UpgradeSales:  outletReport.UpgradeSalesCount,
			HBBTarget:     outletReport.HBBTargetCount,
			HayyakTarget:  outletReport.HayyakTargetCount,
			BaqatiTarget:  outletReport.BaqatiTargetCount,
			UpgradeTarget: outletReport.UpgradeTargetCount,
		}

		eachOutletResp.TotalSales = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount
		eachOutletResp.TotalTarget = outletReport.HBBTargetCount + outletReport.HayyakTargetCount + outletReport.BaqatiTargetCount + outletReport.UpgradeTargetCount

		if eachOutletResp.HBBTarget != 0 {
			eachOutletResp.HBBAchievedPercentage = math.Round(float64(eachOutletResp.HBBSales*100/eachOutletResp.HBBTarget)*100) / 100
		}

		if eachOutletResp.HayyakTarget != 0 {
			eachOutletResp.HayyakAchievedPercentage = math.Round(float64(eachOutletResp.HayyakSales*100/eachOutletResp.HayyakTarget)*100) / 100
		}

		if eachOutletResp.BaqatiTarget != 0 {
			eachOutletResp.BaqatiAchievedPercentage = math.Round(float64(eachOutletResp.BaqatiSales*100/eachOutletResp.BaqatiTarget)*100) / 100
		}

		if eachOutletResp.UpgradeTarget != 0 {
			eachOutletResp.UpgradeAchievedPercentage = math.Round(float64(eachOutletResp.UpgradeSales*100/eachOutletResp.UpgradeTarget)*100) / 100
		}

		if eachOutletResp.TotalTarget != 0 {
			eachOutletResp.TotalAchievedPercentage = math.Round(float64(eachOutletResp.TotalSales*100/eachOutletResp.TotalTarget)*100) / 100
		}

		resp = append(resp, eachOutletResp)
	}

	return resp, nil
}
