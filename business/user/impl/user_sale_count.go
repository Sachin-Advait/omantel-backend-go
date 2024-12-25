package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) GetUserSaleCount(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month,
	year int) (interface{}, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	// _, ok := constant.ProductGroupMap[product]
	// if !ok {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusBadRequest,
	// 		Code:       errors.InvalidReqParameters,
	// 		Message:    "INVALID PRODUCT - " + product,
	// 		Details:    "invalid product - " + product,
	// 	}
	// }

	report, err := impl.DbService.FetchUserSaleCount(userID, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchUserSaleCount: provider error - " + err.Error(),
		}
	}

	return report, nil
}
