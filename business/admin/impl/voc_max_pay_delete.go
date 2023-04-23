package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) DeleteVOCMaxPay(ctx *gin.Context, req []*models.FetchDeleteVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error) {

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

	position := make([]string, 0)
	year := make([]int, 0)

	for _, v := range req {
		position = append(position, v.Position)
		year = append(year, v.Year)
	}

	return impl.DbService.DeleteVOCMaxPayRecord(ctxHeader.XUserId, position, year)
}
