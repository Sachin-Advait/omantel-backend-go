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

func GetUsersReportingToGeneralManager(ctx *gin.Context, req *models.UserReporting) (int, interface{}) {

	id, err := strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUsersReportingToGeneralManager : incorrect client id")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "invalid userID",
		}
	}

	response, errs := theProvider.GetUsersReportingToGeneralManager(ctx, id)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GetUsersReportingToGeneralManager : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUsersReportingToGeneralManager : success")

	return http.StatusOK, response
}
