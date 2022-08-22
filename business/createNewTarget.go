package business

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewTarget(ctx *gin.Context, req []*models.NewTargetRequest) ([]*models.NewTargetFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	for index, val := range req {
		req[index].TargetCode = val.RegionCode + val.OutletCode + fmt.Sprintf("%d", val.Month)
	}

	switch ctxHeader.XAPIVersion {
	case "1":
		return nil, impl.PsqlDb.UpsertTargetRecord(req[0])
	case "2":
		return impl.PsqlDb.UpsertTargetRecordViaFile(req), nil
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "api version incorrect",
		}
	}
}
