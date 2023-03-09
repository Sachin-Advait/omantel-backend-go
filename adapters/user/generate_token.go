package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GenerateToken(ctx *gin.Context) (int, interface{}) {

	token, err := userAPIProvider.GenerateToken(ctx)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(err.Details.(string))).Msg("GenerateToken : provider failed")
		return err.StatusCode, &models.APIFailureResponse{
			Message: err.Message,
			Error:   err.Details.(string),
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, "success").Msg("GenerateToken : success")

	return http.StatusOK, &models.APISuccessResponse{Status: true, Data: token}
}
