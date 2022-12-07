package business

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

//GetOutletSalesAndTarget: userID should be a manager
func (impl *IncentiveImpl) GetOutletSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(userId)
	if err != nil {
		return nil, err
	}

	manger, ok := constant.ManagerMap[position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	allDistinctOM, err := impl.PsqlDb.FetchAllOutletManagersReportingToGivenManagerRole(manger, userId)
	if err != nil {
		return nil, err
	}

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	for _, v := range allDistinctOM {

		outletReport, err := impl.PsqlDb.FetchOutletReport(v.Outlet, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			outletReport, err = impl.updateOutletReport(ctx, v.Region, v.Outlet, isQuarterReport, quarter, month, year)
		}

		if err != nil || outletReport == nil {
			return nil, err
		}

		eachOutletResp := &models.ProductWiseSalesTargetPerOutlet{
			OutletName:    v.Outlet,
			Region:        v.Region,
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
		eachOutletResp.TotalTarget = eachOutletResp.HBBTarget + eachOutletResp.HayyakTarget + eachOutletResp.BaqatiTarget + eachOutletResp.UpgradeTarget

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
