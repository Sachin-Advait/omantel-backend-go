package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteRegions(ctx *gin.Context, req *models.RegionRequest) *errors.Error {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if req.Quarter == 0 {
		req.Quarter = currQuarter
	}

	if req.Year == 0 {
		req.Year = currYear
	}

	return impl.PsqlDb.DeleteRegionRecord(ctxHeader.XUserId, req.RegionCodes, req.Quarter, req.Year)
}
