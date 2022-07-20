package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewPlan(ctx *gin.Context, req *models.NewPlanRequest) *errors.Error {

	return impl.PsqlDb.UpsertPlanRecord(req)
}
