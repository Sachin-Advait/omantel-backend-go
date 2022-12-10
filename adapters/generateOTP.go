package adapters

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GenerateOTP(ctx *gin.Context, req *models.GenerateOTP) (int, interface{}) {

	response, errs := theProvider.GenerateOTP(ctx, req.ClientId)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GenerateOTP : provider failed")
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

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GenerateOTP : success")
	logger.LogChan <- &logger.Logging{
		UserID:     req.ClientId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Code:       http.StatusOK,
		Text:       fmt.Sprintf("%+v", response),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
