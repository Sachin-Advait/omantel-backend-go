package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) DeleteRegions(ctx *gin.Context, req *models.FetchDeleteRegionRequest) *errors.Error {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	return impl.DbService.DeleteRegionRecord(ctxHeader.XUserId, req.RegionCodes, req.Quarter, req.Year)
}
