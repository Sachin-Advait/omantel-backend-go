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
)

func GetUserReport(ctx *gin.Context, req *models.UserReportRequest) (int, interface{}) {

	id, err := strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUserReport : incorrect client id")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "invalid userID",
		}
	}

	if !((req.Month >= 1 && req.Month <= 12) || (req.Quarter >= 1 && req.Quarter <= 4)) {
		logger.FileLogger.CtxError(ctx).Err(errors.New("monthly and quarterly type is incorrect")).Stack().Msg("GetUserReport : incorrect client id")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "monthly and quarterly type is incorrect",
		}
	}

	var isQuarterReport bool

	if req.Quarter >= 1 && req.Quarter <= 4 {
		isQuarterReport = true
	}

	response, errs := theProvider.GetUserReport(ctx, id, req.Month, req.Quarter, isQuarterReport)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GetUserReport : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUserReport : success")

	return http.StatusOK, response
}
