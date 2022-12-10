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
)

func FetchTargets(ctx *gin.Context, req *models.GetTargetRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	//TODO: check

	response, errs := theProvider.FetchTargets(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("CreateNewTarget : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
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

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("CreateNewTarget : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", response),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
