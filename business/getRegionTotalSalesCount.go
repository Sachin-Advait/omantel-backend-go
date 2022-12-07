package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetRegionTotalSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	for _, v := range allDistinctOM {

		regionReport, err := impl.PsqlDb.FetchRegionReport(v.Region, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			regionReport, err = impl.updateRegionReport(ctx, v.Region, isQuarterReport, quarter, month, year)
		}

		if err != nil || regionReport == nil {
			return nil, err
		}

		eachOutletResponse := &models.ProductWiseSalesCountPerOutlet{
			OutletName: v.Outlet,
			Region:     v.Region,
			HBB:        regionReport.HBBSalesCount,
			Hayyak:     regionReport.HayyakSalesCount,
			Baqati:     regionReport.BaqatiSalesCount,
			Upgrade:    regionReport.UpgradeSalesCount,
		}

		resp = append(resp, eachOutletResponse)
	}

	return resp, nil
}
