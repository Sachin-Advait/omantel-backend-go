package adapters

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GetUserReport(ctx *gin.Context, req *models.UserReportRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, err := validateUserReportRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetUserReport : validation failed")
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

	response, errs := theProvider.GetUserReport(ctx, req.ClientId, isQuarterReport, req.Quarter, req.Month, time.Now().Year())
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetUserReport : provider failed")
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

	validateUserReportResponse(ctx, response)

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUserReport : success")
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

func validateUserReportRequest(req *models.UserReportRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if req.ClientId == "" || ctxHeader.XUserId != req.ClientId {
	// 	err = errors.New("incorrect user token")
	// 	return
	// }

	if !((req.Month >= 1 && req.Month <= 12) || (req.Quarter >= 1 && req.Quarter <= 4)) {
		err = errors.New("monthly or quarterly type is incorrect")
		return
	}

	if req.Quarter >= 1 && req.Quarter <= 4 {
		isQuarterReport = true
	}

	return isQuarterReport, nil
}

func validateUserReportResponse(ctx *gin.Context, res *models.UserReportDB) {

	var achievedCount int

	if res.HBBPayout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("HBB", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.HayyakPayout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("HAYYAK", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.BaqatiPayout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("BAQATI", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.UpgradePayout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("UPGRADE", "not achieved").Msg("GetUserReport : response validation")
	}

	if achievedCount < 3 {
		res.HBBPayout = 0
		res.HBBExtraPayout = 0
		res.HayyakPayout = 0
		res.HayyakExtraPayout = 0
		res.BaqatiPayout = 0
		res.BaqatiExtraPayout = 0
		res.UpgradePayout = 0
		res.UpgradeExtraPayout = 0
		res.TotalPayout = 0

		return
	}

	res.HBBTargetCount = math.Round(res.HBBTargetCount*1000) / 1000
	res.HBBPayout = math.Round(res.HBBPayout*1000) / 1000
	res.HBBExtraPayout = math.Round(res.HBBExtraPayout*1000) / 1000
	res.HBBAchievedPercentage = math.Round(res.HBBAchievedPercentage*1000) / 1000

	res.HayyakTargetCount = math.Round(res.HayyakTargetCount*1000) / 1000
	res.HayyakPayout = math.Round(res.HayyakPayout*1000) / 1000
	res.HayyakExtraPayout = math.Round(res.HayyakExtraPayout*1000) / 1000
	res.HayyakAchievedPercentage = math.Round(res.HayyakAchievedPercentage*1000) / 1000

	res.BaqatiTargetCount = math.Round(res.BaqatiTargetCount*1000) / 1000
	res.BaqatiPayout = math.Round(res.BaqatiPayout*1000) / 1000
	res.BaqatiExtraPayout = math.Round(res.BaqatiExtraPayout*1000) / 1000
	res.BaqatiAchievedPercentage = math.Round(res.BaqatiAchievedPercentage*1000) / 1000

	res.UpgradeTargetCount = math.Round(res.UpgradeTargetCount*1000) / 1000
	res.UpgradePayout = math.Round(res.UpgradePayout*1000) / 1000
	res.UpgradeExtraPayout = math.Round(res.UpgradeExtraPayout*1000) / 1000
	res.UpradeAchievedPercentage = math.Round(res.UpradeAchievedPercentage*1000) / 1000

	res.TotalPayout = math.Round(res.TotalPayout*1000) / 1000
}
