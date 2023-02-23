package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteOutlets(ctx *gin.Context, req *models.FetchDeleteOutletRequest) *errors.Error {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	return impl.PsqlDb.DeleteOutletRecord(ctxHeader.XUserId, req.OutletCodes, req.Quarter, req.Year)
}
