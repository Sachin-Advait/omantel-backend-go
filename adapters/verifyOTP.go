package adapters

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func VerifyOTP(ctx *gin.Context, req *models.VerifyOTP) (int, interface{}) {

	response, errs := theProvider.VerifyOTP(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("VerifyOTP : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     req.ClientId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Code:       errs.StatusCode,
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, "success").Msg("VerifyOTP : success")
	logger.LogChan <- &logger.Logging{
		UserID:     req.ClientId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Code:       http.StatusOK,
		Text:       "success response",
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
