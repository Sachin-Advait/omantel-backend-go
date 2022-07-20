package adapters

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetRegionSalesAndTarget(ctx *gin.Context, req *models.SingleUserIncentiveRequest) (int, interface{}) {

	id, err := strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetRegionSalesAndTarget : incorrect client id")
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
		logger.FileLogger.CtxError(ctx).Err(errors.New("monthy and quartlerly type is blank")).Stack().Msg("GetRegionSalesAndTarget : incorrect type")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "monthy and quartlerly type is blank",
		}
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetRegionSalesAndTarget : incorrect monthly time")
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
				Msg("GetRegionSalesAndTarget : incorrect quarterly type")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "quarterly type is incorrect, avoid spaces",
			}
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetRegionSalesAndTarget : incorrect quarterly parsing")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "quarter value not found",
			}

		} else if quarter < 1 || quarter > 4 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("incorrect quarter value - should be between 1 and 4")).Stack().
				Msg("GetRegionSalesAndTarget : incorrect quarterly value")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "incorrect quarter value - should be between 1 and 4",
			}
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetRegionSalesAndTarget : incorrect year time")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "year type is incorrect, use '2006' format",
			}
		}
	}

	response, errs := theProvider.GetRegionSalesAndTarget(ctx, id, quarter, isQuarterReport, date)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GetRegionSalesAndTarget : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetRegionSalesAndTarget : success")

	return http.StatusOK, response
}
