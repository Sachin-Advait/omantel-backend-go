package user

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GetRegionSalesAndTarget(ctx *gin.Context, req *models.UserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := models.ValidateUserIncentiveRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetRegionSalesAndTarget : validation failed")
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

	response, errs := userAPIProvider.GetRegionSalesAndTarget(ctx, req.ClientId, isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetRegionSalesAndTarget : provider failed")
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

	for index, val := range response {
		response[index].BaqatiAchievedPercentage = math.Round(val.BaqatiAchievedPercentage*1000) / 1000
		response[index].HBBAchievedPercentage = math.Round(val.HBBAchievedPercentage*1000) / 1000
		response[index].HayyakAchievedPercentage = math.Round(val.HayyakAchievedPercentage*1000) / 1000
		response[index].UpgradeAchievedPercentage = math.Round(val.UpgradeAchievedPercentage*1000) / 1000
		response[index].TotalAchievedPercentage = math.Round(val.TotalAchievedPercentage*1000) / 1000
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetRegionSalesAndTarget : success")
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
