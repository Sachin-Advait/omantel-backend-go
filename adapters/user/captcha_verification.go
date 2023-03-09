package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CaptchaVerification(ctx *gin.Context, req *models.CaptchRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	resp, errs := userAPIProvider.CaptchaVerification(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("CaptchaVerification : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Code:       errs.StatusCode,
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error:   errs.Code,
			Message: errs.Message,
		}
	}

	respBytes, _ := json.Marshal(resp)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, string(respBytes)).Msg("CaptchaVerification : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Code:       http.StatusOK,
		Text:       string(respBytes),
		CreateTime: time.Now().Local(),
	}

	response := models.APISuccessResponse{
		Status: true,
		Data:   resp,
	}

	return http.StatusOK, response
}
