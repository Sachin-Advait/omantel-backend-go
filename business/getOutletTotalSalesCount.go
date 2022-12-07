package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetOutletTotalSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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

		outletReport, err := impl.PsqlDb.FetchOutletReport(v.Outlet, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			outletReport, err = impl.updateOutletReport(ctx, v.Region, v.Outlet, isQuarterReport, quarter, month, year)
		}

		if err != nil || outletReport == nil {
			return nil, err
		}

		eachOutletResponse := &models.ProductWiseSalesCountPerOutlet{
			OutletName: v.Outlet,
			Region:     v.Region,
			Latitude:   outletReport.Latitude,
			Longitude:  outletReport.Longitude,
			HBB:        outletReport.HBBSalesCount,
			Hayyak:     outletReport.HayyakSalesCount,
			Baqati:     outletReport.BaqatiSalesCount,
			Upgrade:    outletReport.UpgradeSalesCount,
		}

		resp = append(resp, eachOutletResponse)
	}

	return resp, nil
}
