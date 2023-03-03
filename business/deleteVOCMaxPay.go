package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteVOCMaxPay(ctx *gin.Context, req []*models.FetchDeleteVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	position := make([]string, 0)
	year := make([]int, 0)

	for _, v := range req {
		position = append(position, v.Position)
		year = append(year, v.Year)
	}

	return impl.PsqlDb.DeleteVOCMaxPayRecord(ctxHeader.XUserId, position, year)
}
