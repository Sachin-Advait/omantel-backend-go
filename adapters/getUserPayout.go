package adapters

import (
	"fmt"
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

	var id, err = strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUserPayout : incorrect client id")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "incorrect client id: " + err.Error(),
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "invalid userID",
		}
	}

	var (
		isQuarterReport = true
		date            time.Time
		quarter         int
	)

	if req.Monthly == "" && req.Quarterly == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("monthy and quartlerly type is blank")).Stack().Msg("GetUserPayout : incorrect type")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "monthy and quartlerly type is blank",
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "monthy and quartlerly type is blank",
		}
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUserPayout : incorrect monthly time")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "incorrect monthly time: " + err.Error(),
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "monthy type is incorrect, use '06-2006' format",
			}
		}

		isQuarterReport = false
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("quarterly type is incorrect, avoid spaces")).Stack().
				Msg("GetUserPayout : incorrect quarterly type")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "quarterly type is incorrect, avoid spaces",
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "quarterly type is incorrect, avoid spaces",
			}
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUserPayout : incorrect quarterly parsing")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "incorrect quarterly parsing: " + err.Error(),
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "quarter value not found",
			}

		} else if quarter < 1 || quarter > 4 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("incorrect quarter value - should be between 1 and 4")).Stack().
				Msg("GetUserPayout : incorrect quarterly value")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "incorrect quarter value - should be between 1 and 4",
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "incorrect quarter value - should be between 1 and 4",
			}
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUserPayout : incorrect year time")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "incorrect year time: " + err.Error(),
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "year type is incorrect, use '2006' format",
			}
		}
	}

	response, errs := theProvider.GetUserPayout(ctx, id, quarter, isQuarterReport, date)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GetUserPayout : provider failed")
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

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUserPayout : success")
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
