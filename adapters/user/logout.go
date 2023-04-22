package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func UserLogout(ctx *gin.Context) (int, interface{}) {

	err := userAPIProvider.UserLogout(ctx)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(err.Details.(string))).Msg("UserLogout : provider failed")
		return err.StatusCode, &models.APIFailureResponse{
			Message: err.Message,
			Error:   err.Details.(string),
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, "success").Msg("UserLogout : success")

	return http.StatusOK, &models.APISuccessResponse{Status: true, Data: "success"}
}
