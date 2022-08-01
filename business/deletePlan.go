package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeletePlans(ctx *gin.Context, req *models.DeletePlanRequest) *errors.Error {

	return impl.PsqlDb.DeletePlanRecord(req.PlanCodes)
}
