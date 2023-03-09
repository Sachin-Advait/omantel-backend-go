package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/admin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func FetchVOCCriteria(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("fetchVOCCriteria : unauthorized request")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "unauthorized request", Code: http.StatusBadRequest,
			Text: "auth token not found", CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.TokenNotFound,
			Message: "AUTHENTICATE YOUR REQUEST",
		})
		return
	}

	if !ctxHeader.IsAdmin {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("fetchVOCCriteria : not an admin")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "not admin", Code: http.StatusForbidden,
			Text: "not an admin", CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusForbidden, models.APIFailureResponse{
			Error:   errors.UnauthorizedUser,
			Message: "USER IS NOT ADMIN",
		})
		return
	}

	ctx.JSON(admin.FetchVOCCriteria(ctx))
}
