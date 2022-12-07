package business

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetRegionSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

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

		regionReport, err := impl.PsqlDb.FetchRegionReport(v.Region, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			regionReport, err = impl.updateRegionReport(ctx, v.Region, isQuarterReport, quarter, month, year)
		}

		if err != nil || regionReport == nil {
			return nil, err
		}

		eachRegionResp := &models.ProductWiseSalesTargetPerOutlet{
			OutletName:    v.Outlet,
			Region:        v.Region,
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
		eachRegionResp.TotalTarget = eachRegionResp.HBBTarget + eachRegionResp.HayyakTarget + eachRegionResp.BaqatiTarget + eachRegionResp.UpgradeTarget

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
