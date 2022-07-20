package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"net/http"
	"time"
)

func (impl *IncentiveImpl) GetOutletTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	for _, v := range allDistinctOM {

		lat, lng, err := impl.PsqlDb.FetchOutletDetails(v.Outlet, v.Region)
		if err != nil {
			return nil, err
		}

		trx, err := impl.PsqlDb.FetchProductSalesPerOutletManager(v.ID, int(t.Month()), t.Year(), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		trx.Latitude = lat
		trx.Longitude = lng
		trx.Region = v.Region
		trx.OutletName = v.Outlet

		resp = append(resp, trx)
	}

	return resp, nil
}
