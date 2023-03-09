package user

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
)

func GetOutletTotalSalesCount(ctx *gin.Context, req *models.UserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := models.ValidateUserIncentiveRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetOutletTotalSalesCount : validation failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	response, errs := userAPIProvider.GetOutletTotalSalesCount(ctx, req.ClientId, isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetOutletTotalSalesCount : provider failed")
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

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetOutletTotalSalesCount : success")
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
