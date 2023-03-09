package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func DeleteVOCUser(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := deleteVOCUserPayload(ctx)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("deleteVOCUser : invalid request param")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "request binding", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: errors.RequestDecodeError,
		}
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("deleteVOCUser : req body")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateDeleteVOCUser(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("deleteVOCUser : req validation failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	customErr := adminAPIProvider.DeleteVOCUer(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("deleteVOCUser : provider failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: customErr.Details.(string), CreateTime: time.Now().Local()})

		return customErr.StatusCode, &models.APIFailureResponse{
			Status:  false,
			Error:   customErr.Err,
			Message: customErr.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, errors.SuccessInfoMessage).Msg("deleteVOCUser : success")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: errors.SuccessInfoMessage, CreateTime: time.Now().Local()})

	return http.StatusOK, models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}
}

func deleteVOCUserPayload(ctx *gin.Context) (*models.FetchDeleteVOCUserRequest, error) {

	userReq := &models.FetchDeleteVOCUserRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return userReq, nil
}
