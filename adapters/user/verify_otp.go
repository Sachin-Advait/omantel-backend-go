package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
)

func VerifyOTP(ctx *gin.Context, req *models.VerifyOTP) (int, interface{}) {

	response, errs := userAPIProvider.VerifyOTP(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("VerifyOTP : provider failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: errs.Details.(string), CreateTime: time.Now().Local()})

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	responseBytes, _ := json.Marshal(response)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, "success").Msg("VerifyOTP : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success", Code: http.StatusOK,
		Text: string(responseBytes), CreateTime: time.Now().Local()})

	return http.StatusOK, response
}
