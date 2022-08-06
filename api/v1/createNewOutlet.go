package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func createNewOutlet(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Stack().Msg("createNewOutlet : unauthorized request")
		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	userReq := &models.NewOutletRequest{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewOutlet : error binding request body")
		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: errors.ErrorCodeMap[errors.RequestDecodeErrorCode],
		})
		return
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(constant.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("createNewOutlet : req body")

	ctx.JSON(adapters.CreateNewOutlet(ctx, userReq))
}
