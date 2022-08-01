package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteOutlets(ctx *gin.Context, req *models.DeleteOutletRequest) *errors.Error {

	return impl.PsqlDb.DeleteOutletRecord(req.OutletCodes)
}
