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

func GetUsersReportingToOutletManager(ctx *gin.Context, req *models.UserReporting) (int, interface{}) {

	id, err := strconv.Atoi(req.ClientId)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("GetUsersReportingToOutletManager : incorrect client id")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "invalid userID",
		}
	}

	response, errs := theProvider.GetUsersReportingToOutletManager(ctx, id)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("GetUsersReportingToOutletManager : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("GetUsersReportingToOutletManager : success")

	return http.StatusOK, response
}
