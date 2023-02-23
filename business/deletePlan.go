package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeletePlans(ctx *gin.Context, req *models.FetchDeletePlanRequest) *errors.Error {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	return impl.PsqlDb.DeletePlanRecord(ctxHeader.XUserId, req.PlanCodes, req.Quarter, req.Year)
}
