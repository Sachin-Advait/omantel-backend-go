package business

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetOutletSalesAndTarget(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(clientID)
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

	allDistinctOM, err := impl.PsqlDb.FetchAllOutletManagersReportingToGivenManagerRole(manger, clientID)
	if err != nil {
		return nil, err
	}

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	for _, v := range allDistinctOM {

		trx, err := impl.PsqlDb.FetchProductSalesPerOutletManager(v.ID, int(t.Month()), t.Year(), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		eachOutletResp := &models.ProductWiseSalesTargetPerOutlet{
			OutletName:    v.Outlet,
			Region:        v.Region,
			HBBSales:      trx.HBB,
			HayyakSales:   trx.Hayyak,
			BaqatiSales:   trx.Baqati,
			UpgradeSales:  trx.Upgrade,
			TotalSales:    trx.Total,
			HBBTarget:     impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, constant.HBBCategory, v.Region, v.Outlet),
			HayyakTarget:  impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, constant.HayyakCategory, v.Region, v.Outlet),
			BaqatiTarget:  impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, constant.BaqatiCategory, v.Region, v.Outlet),
			UpgradeTarget: impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, constant.UpgradeCategory, v.Region, v.Outlet),
		}

		eachOutletResp.TotalTarget = eachOutletResp.HBBTarget + eachOutletResp.HayyakTarget + eachOutletResp.BaqatiTarget + eachOutletResp.UpgradeTarget

		eachOutletResp.HBBAchievedPercentage = math.Round(float64(eachOutletResp.HBBSales*100/eachOutletResp.HBBTarget)*100) / 100
		eachOutletResp.HayyakAchievedPercentage = math.Round(float64(eachOutletResp.HayyakSales*100/eachOutletResp.HayyakTarget)*100) / 100
		eachOutletResp.BaqatiAchievedPercentage = math.Round(float64(eachOutletResp.BaqatiSales*100/eachOutletResp.BaqatiTarget)*100) / 100
		eachOutletResp.UpgradeAchievedPercentage = math.Round(float64(eachOutletResp.UpgradeSales*100/eachOutletResp.UpgradeTarget)*100) / 100
		eachOutletResp.TotalAchievedPercentage = math.Round(float64(eachOutletResp.TotalSales*100/eachOutletResp.TotalTarget)*100) / 100

		resp = append(resp, eachOutletResp)
	}

	return resp, nil
}
