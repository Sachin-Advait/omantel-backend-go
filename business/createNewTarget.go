package business

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *IncentiveImpl) CreateNewTarget(ctx *gin.Context, req []*models.CreateTargetRequest) ([]*models.TargetFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	for index, val := range req {
		val.Quarter = impl.getQuarterFromMonth(val.Month)
		if val.Year == 0 {
			val.Year = currYear
		}

		val.TargetCode = strings.ToUpper(objecthelper.RandStringRunes(10))

		req[index] = val
	}

	switch ctxHeader.XAPIVersion {
	case "1":
		return nil, impl.PsqlDb.UpsertTargetRecord(ctxHeader.XUserId, req[0])
	case "2":
		return impl.PsqlDb.UpsertTargetRecordViaFile(ctxHeader.XUserId, req), nil
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "api version incorrect",
		}
	}
}
