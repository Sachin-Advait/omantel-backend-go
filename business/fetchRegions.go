package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchRegions(ctx *gin.Context, req *models.FetchDeleteRegionRequest) ([]*models.RegionResponse, *errors.Error) {
	return impl.PsqlDb.FetchRegionRecord(req.RegionCodes, req.Quarter, req.Year)
}
