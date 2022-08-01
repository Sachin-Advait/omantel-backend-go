package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteTargets(ctx *gin.Context, req *models.DeleteTargetRequest) *errors.Error {

	return impl.PsqlDb.DeleteTargetRecord(req.TargetCodes)
}
