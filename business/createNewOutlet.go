package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewOutlet(ctx *gin.Context, req []*models.NewOutletRequest) ([]*models.NewOutletFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	switch ctxHeader.XAPIVersion {
	case "1":
		return nil, impl.PsqlDb.UpsertOutletRecord(req[0])
	case "2":
		return impl.PsqlDb.UpsertOutletRecordViaFile(req), nil
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "api version incorrect",
		}
	}
}
