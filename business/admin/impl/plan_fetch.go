package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchPlans(ctx *gin.Context, req *models.FetchDeletePlanRequest) ([]*models.PlanResponse, *errors.Error) {
	return impl.DbService.FetchPlanRecord(req.PlanCodes, req.Quarter, req.Year)
}
