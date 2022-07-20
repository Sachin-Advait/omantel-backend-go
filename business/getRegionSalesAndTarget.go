package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"math"
	"net/http"
	"time"
)

func (impl *IncentiveImpl) GetRegionSalesAndTarget(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(clientID)
	if err != nil {
		return nil, err
	}

	manger, ok := commons.ManagerMap[position]
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

	temp := make(map[string]*psqldb.ProductWiseSalesTargetPerOutlet)

	for _, v := range allDistinctOM {

		trx, err := impl.PsqlDb.FetchProductSalesPerOutletManager(v.ID, int(t.Month()), t.Year(), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		eachRegionRcrd, ok := temp[v.Region]
		if !ok {
			temp[v.Region] = &psqldb.ProductWiseSalesTargetPerOutlet{
				Region:        v.Region,
				HBBSales:      trx.HBB,
				HayyakSales:   trx.Hayyak,
				BaqatiSales:   trx.Baqati,
				UpgradeSales:  trx.Upgrade,
				HBBTarget:     impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.HBBCategory, v.Region, v.Outlet),
				HayyakTarget:  impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.HayyakCategory, v.Region, v.Outlet),
				BaqatiTarget:  impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.BaqatiCategory, v.Region, v.Outlet),
				UpgradeTarget: impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.UpgradeCategory, v.Region, v.Outlet),
			}

			continue
		}

		eachRegionRcrd.HBBSales += trx.HBB
		eachRegionRcrd.HayyakSales += trx.Hayyak
		eachRegionRcrd.BaqatiSales += trx.Baqati
		eachRegionRcrd.UpgradeSales += trx.Upgrade
		eachRegionRcrd.HBBTarget += impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.HBBCategory, v.Region, v.Outlet)
		eachRegionRcrd.HayyakTarget += impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.HayyakCategory, v.Region, v.Outlet)
		eachRegionRcrd.BaqatiTarget += impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.BaqatiCategory, v.Region, v.Outlet)
		eachRegionRcrd.UpgradeTarget += impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.UpgradeCategory, v.Region, v.Outlet)

	}

	resp := make([]*psqldb.ProductWiseSalesTargetPerOutlet, 0)

	for _, v := range temp {

		v.TotalSales = v.HBBSales + v.HayyakSales + v.BaqatiSales + v.UpgradeSales
		v.TotalTarget = v.HBBTarget + v.HayyakTarget + v.BaqatiTarget + v.UpgradeTarget

		v.HBBAchievedPercentage = math.Round(float64(v.HBBSales*100/v.HBBTarget)*100) / 100
		v.HayyakAchievedPercentage = math.Round(float64(v.HayyakSales*100/v.HayyakTarget)*100) / 100
		v.BaqatiAchievedPercentage = math.Round(float64(v.BaqatiSales*100/v.BaqatiTarget)*100) / 100
		v.UpgradeAchievedPercentage = math.Round(float64(v.UpgradeSales*100/v.UpgradeTarget)*100) / 100
		v.TotalAchievedPercentage = math.Round(float64(v.TotalSales*100/v.TotalTarget)*100) / 100

		resp = append(resp, v)
	}

	return resp, nil
}
