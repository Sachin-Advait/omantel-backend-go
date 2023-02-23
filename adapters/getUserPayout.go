package adapters

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GetUserPayout(ctx *gin.Context, req *models.SingleUserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := validateUserPayoutRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetUserPayout : validation failed")
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

	response, errs := theProvider.GetUserPayout(ctx, req.ClientId, isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetUserPayout : provider failed")
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

	validateUserPayoutResponse(ctx, response)

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUserPayout : success")
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

func validateUserPayoutRequest(req *models.SingleUserIncentiveRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if req.ClientId == "" || ctxHeader.XUserId != req.ClientId {
	// 	err = errors.New("incorrect user token")
	// 	return
	// }

	var date time.Time

	isQuarterReport = true

	if req.Monthly == "" && req.Quarterly == "" {
		err = errors.New("monthy and quartlerly type is blank")
		return
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			err = errors.New("monthly parsing: " + err.Error())
			return
		}

		isQuarterReport = false
		month = int(date.Month())
		year = date.Year()
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			err = errors.New("quarterly type is incorrect, avoid spaces")
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			err = errors.New("quarter conversion: " + err.Error())
			return

		}

		if quarter < 1 || quarter > 4 {
			err = errors.New("incorrect quarter value - should be between 1 and 4")
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			err = errors.New("quarter year parsing: " + err.Error())
			return
		}

		year = date.Year()
	}

	return
}

func validateUserPayoutResponse(ctx *gin.Context, res *models.TotalIncentive) {

	var achievedCount int

	if res.HBB == nil {
		res.HBB = &models.IndividualIncentive{}
	}
	if res.Hayyak == nil {
		res.Hayyak = &models.IndividualIncentive{}
	}
	if res.Baqati == nil {
		res.Baqati = &models.IndividualIncentive{}
	}
	if res.Upgrade == nil {
		res.Upgrade = &models.IndividualIncentive{}
	}

	if res.HBB.Payout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("HBB", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Hayyak.Payout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("HAYYAK", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Baqati.Payout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("BAQATI", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Upgrade.Payout != 0 {
		achievedCount++
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("UPGRADE", "not achieved").Msg("GetUserPayout : response validation")
	}

	if achievedCount < 3 {
		res.TotalPay = 0
		res.HBB.Payout = 0
		res.HBB.ExtraPayout = 0
		res.Hayyak.Payout = 0
		res.Hayyak.ExtraPayout = 0
		res.Baqati.Payout = 0
		res.Baqati.ExtraPayout = 0
		res.Upgrade.Payout = 0
		res.Upgrade.ExtraPayout = 0

		return
	}

	res.TotalPay = math.Round(res.TotalPay*1000) / 1000
	res.HBB.Payout = math.Round(res.HBB.Payout*1000) / 1000
	res.HBB.ExtraPayout = math.Round(res.HBB.ExtraPayout*1000) / 1000
	res.Hayyak.Payout = math.Round(res.Hayyak.Payout*1000) / 1000
	res.Hayyak.ExtraPayout = math.Round(res.Hayyak.ExtraPayout*1000) / 1000
	res.Baqati.Payout = math.Round(res.Baqati.Payout*1000) / 1000
	res.Baqati.ExtraPayout = math.Round(res.Baqati.ExtraPayout*1000) / 1000
	res.Upgrade.Payout = math.Round(res.Upgrade.Payout*1000) / 1000
	res.Upgrade.ExtraPayout = math.Round(res.Upgrade.ExtraPayout*1000) / 1000
}
