package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchPlans(ctx *gin.Context, req *models.FetchDeletePlanRequest) ([]*models.PlanResponse, *errors.Error) {
	return impl.PsqlDb.FetchPlanRecord(req.PlanCodes, req.Quarter, req.Year)
}
