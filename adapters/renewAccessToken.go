package adapters

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func RenewAccessToken(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	response, errs := theProvider.RenewAccessToken(ctx)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("RenewAccessToken : provider failed")
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
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, "success").Msg("RenewAccessToken : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Code:       http.StatusOK,
		Text:       "success response",
		CreateTime: time.Now().Local(),
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   response,
	}

	return http.StatusOK, resp
}
