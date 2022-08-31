package adapters

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CaptchaVerification(ctx *gin.Context, req *models.CaptchRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if req.Token == "" || ctxHeader.XUserId == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid req param")).Stack().Msg("CaptchaVerification : invalid req params")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "invalid req param",
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Blank value found",
		}
	}

	resp, errs := theProvider.CaptchaVerification(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CaptchaVerification : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error:   errs.Code,
			Message: errs.Message,
		}
	}

	response := models.APISuccessResponse{
		Status: true,
		Data:   resp,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("CaptchaVerification : logging failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "logging failed",
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error:   errors.InternalServerErrorCode,
			Message: errors.InternalServerError,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, string(respBytes)).Msg("CaptchaVerification : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       string(respBytes),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
