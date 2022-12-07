package business

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetSaleCountPerPlan(ctx *gin.Context, userId, product string, isQuarterReport bool, quarter, month, year int) ([]*models.SaleCountPerPlan, *errors.Error) {

	_, ok := constant.ProductGroupMap[product]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "Invalid product",
			Details:    "Invalid product",
		}
	}

	resp, err := impl.PsqlDb.FetchUserSalesCountPerPlan(userId, product, month, year, quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Quantity > resp[j].Quantity
	})

	return resp, nil
}
