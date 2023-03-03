package adapters

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func FetchOutlets(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := fetchOutletPayload(ctx)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("fetchOutlets : invalid request param")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "request binding", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: errors.RequestDecodeError,
		}
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("fetchOutlets : req body")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateFetchOutlet(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("fetchOutlets : req validation failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	response, customErr := theProvider.FetchOutlets(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("fetchOutlets : provider failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: customErr.Details.(string), CreateTime: time.Now().Local()})

		return customErr.StatusCode, &models.APIFailureResponse{
			Status:  false,
			Error:   customErr.Err,
			Message: customErr.Message,
		}
	}

	res, _ := objecthelper.MarshalJSON(response)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, string(res)).Msg("fetchOutlets : success")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: string(res), CreateTime: time.Now().Local()})

	return http.StatusOK, response
}

func fetchOutletPayload(ctx *gin.Context) (*models.FetchDeleteOutletRequest, error) {

	userReq := &models.FetchDeleteOutletRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return userReq, nil
}
