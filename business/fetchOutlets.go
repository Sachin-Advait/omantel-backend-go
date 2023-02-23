package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchOutlets(ctx *gin.Context, req *models.FetchDeleteOutletRequest) ([]*models.OutletResponse, *errors.Error) {
	return impl.PsqlDb.FetchOutletRecord(req.OutletCodes, req.Quarter, req.Year)
}
