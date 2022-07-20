package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewTarget(ctx *gin.Context, req *models.NewTargetRequest) *errors.Error {

	return impl.PsqlDb.UpsertTargetRecord(req)
}
