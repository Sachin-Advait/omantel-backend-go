package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
)

func GenerateOTP(ctx *gin.Context, req *models.GenerateOTP) (int, interface{}) {

	var (
		response map[string]string
		errs     *errors.Error
	)

	env := flags.Env()

	if env == "prod" {
		response, errs = userAPIProvider.GenerateOTP(ctx, req.ClientId)
		if errs != nil {
			logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GenerateOTP : provider failed")
			dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(), Type: "error",
				Msg:  "provider failed",
				Text: errs.Details.(string), CreateTime: time.Now().Local()})

			return errs.StatusCode, &models.APIFailureResponse{
				Error: errs.Message,
			}
		}
	} else {
		response, errs = userAPIProvider.GenerateOTPDev(ctx, req.ClientId)
		if errs != nil {
			logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GenerateOTP : provider failed")
			dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(), Type: "error",
				Msg:  "provider failed",
				Text: errs.Details.(string), CreateTime: time.Now().Local()})

			return errs.StatusCode, &models.APIFailureResponse{
				Error: errs.Message,
			}
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).
		Msg("GenerateOTP : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(),
		Type: "info", Msg: "success", Code: http.StatusOK,
		Text: "", CreateTime: time.Now().Local()})

	return http.StatusOK, response
}
