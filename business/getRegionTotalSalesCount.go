package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"net/http"
	"time"
)

func (impl *IncentiveImpl) GetRegionTotalSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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

	temp := make(map[string]*models.ProductWiseSalesCountPerOutlet)

	for _, v := range allDistinctOM {

		trx, err := impl.PsqlDb.FetchProductSalesPerOutletManager(v.ID, int(t.Month()), t.Year(), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		eachRegionRcrd, ok := temp[v.Region]
		if !ok {
			trx.Region = v.Region
			temp[v.Region] = trx
			continue
		}

		eachRegionRcrd.HBB += trx.HBB
		eachRegionRcrd.Hayyak += trx.Hayyak
		eachRegionRcrd.Baqati += trx.Baqati
		eachRegionRcrd.Upgrade += trx.Upgrade
		eachRegionRcrd.Total += trx.Total
	}

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	for _, m := range temp {

		resp = append(resp, m)
	}

	return resp, nil
}
