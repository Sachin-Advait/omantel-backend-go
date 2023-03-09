package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchOutlets(ctx *gin.Context, req *models.FetchDeleteOutletRequest) ([]*models.OutletResponse, *errors.Error) {
	return impl.DbService.FetchOutletRecord(req.OutletCodes, req.Quarter, req.Year)
}
