package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/admin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func GetUserLoginCount(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("getUserLoginCount : unauthorized request")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "unauthorized request", Code: http.StatusBadRequest,
			Text: "auth token not found", CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	if !ctxHeader.IsAdmin {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("getUserLoginCount : not an admin")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "not admin", Code: http.StatusForbidden,
			Text: "not an admin", CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusForbidden, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	userReq := &models.UserLoginCountRequest{
		Region: ctx.Query("region"),
	}

	var err error
	from, fromErr := time.Parse("02-01-2006", ctx.Query("from"))
	to, toErr := time.Parse("02-01-2006", ctx.Query("to"))
	if fromErr != nil || toErr != nil {
		if fromErr != nil {
			err = fromErr
		} else {
			err = toErr
		}
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("getUserLoginCount : validation failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "validation failed", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		})
		return
	}
	userReq.From = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local)
	userReq.To = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, time.Local)

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("getUserLoginCount : req body")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	ctx.JSON(admin.GetUserLoginCount(ctx, userReq))
}
