package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) UpsertVOCCriteria(ctx *gin.Context, req []*models.CreateVOCCriteriaRequest) ([]*models.VOCCriteriaFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	return impl.DbService.UpsertVOCCriteriaRecord(ctxHeader.XUserId, req)
}
