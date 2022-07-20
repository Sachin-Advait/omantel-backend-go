package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewRegion(ctx *gin.Context, req *models.NewRegionRequest) *errors.Error {

	return impl.PsqlDb.UpsertRegionRecord(req)
}
