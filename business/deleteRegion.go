package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteRegions(ctx *gin.Context, req *models.DeleteRegionRequest) *errors.Error {

	return impl.PsqlDb.DeleteRegionRecord(req.RegionCodes)
}
