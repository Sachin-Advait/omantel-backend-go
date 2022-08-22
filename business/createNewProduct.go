package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewProduct(ctx *gin.Context, req []*models.NewProductRequest) ([]*models.NewProductFailedResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	switch ctxHeader.XAPIVersion {
	case "1":
		return nil, impl.PsqlDb.UpsertProductRecord(req[0])
	case "2":
		return impl.PsqlDb.UpsertProductRecordViaFile(req), nil
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "api version incorrect",
		}
	}
}
