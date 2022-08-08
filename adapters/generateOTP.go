package adapters

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GenerateOTP(ctx *gin.Context, req *models.GenerateOTP) (int, interface{}) {

	id, err := strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GenerateOTP : incorrect client id")
		logger.LogChan <- &logger.Logging{
			UserID:     "",
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "invalid userID",
		}
	}

	response, errs := theProvider.GenerateOTP(ctx, id)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GenerateOTP : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     req.ClientId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GenerateOTP : success")
	logger.LogChan <- &logger.Logging{
		UserID:     req.ClientId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", response),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
