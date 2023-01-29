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

func GetUserTarget(ctx *gin.Context, req *models.SingleUserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := validateUserTargetRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetUserTarget : validation failed")
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

	response, errs := theProvider.GetUserTarget(ctx, req.ClientId, isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetUserTarget : provider failed")
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

	validateUserTargetResponse(response)

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUserTarget : success")
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

func validateUserTargetRequest(req *models.SingleUserIncentiveRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

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

func validateUserTargetResponse(res *models.ChartDataResponse) {

	res.HBB.Target = math.Round(res.HBB.Target*1000) / 1000
	res.HBB.AchievedPercentage = math.Round(res.HBB.AchievedPercentage*1000) / 1000
	res.Hayyak.Target = math.Round(res.Hayyak.Target*1000) / 1000
	res.Hayyak.AchievedPercentage = math.Round(res.Hayyak.AchievedPercentage*1000) / 1000
	res.Baqati.Target = math.Round(res.Baqati.Target*1000) / 1000
	res.Baqati.AchievedPercentage = math.Round(res.Baqati.AchievedPercentage*1000) / 1000
	res.Upgrade.Target = math.Round(res.Upgrade.Target*1000) / 1000
	res.Upgrade.AchievedPercentage = math.Round(res.Upgrade.AchievedPercentage*1000) / 1000
}
