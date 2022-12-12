package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewProduct(ctx *gin.Context, req []*models.CreateProductRequest) ([]*models.ProductFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	for index, val := range req {
		if val.Quarter == 0 {
			val.Quarter = currQuarter
		}

		if val.Year == 0 {
			val.Year = currYear
		}

		req[index] = val
	}

	switch ctxHeader.XAPIVersion {
	case "1":
		return nil, impl.PsqlDb.UpsertProductRecord(ctxHeader.XUserId, req[0])
	case "2":
		return impl.PsqlDb.UpsertProductRecordViaFile(ctxHeader.XUserId, req), nil
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "api version incorrect",
		}
	}
}
