package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
)

func GetSaleCountPerPlan(ctx *gin.Context, req *models.UserSalesRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := models.ValidateUserSalesRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetSaleCountPerPlan : validation failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(),
			Type: "error", Msg: "invalid req params", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	product := strings.ToUpper(req.Product)

	response, errs := userAPIProvider.GetSaleCountPerPlan(ctx, req.ClientId, product, req.Channel,
		isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).
			Msg("GetSaleCountPerPlan : provider failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(),
			Type: "error", Msg: "provider failed",
			Text: errs.Details.(string), CreateTime: time.Now().Local()})

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	responseBytes, _ := json.Marshal(response)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", response)).
		Msg("GetSaleCountPerPlan : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: req.ClientId, URL: ctx.Request.URL.String(),
		Type: "info", Msg: "success", Code: http.StatusOK,
		Text: string(responseBytes), CreateTime: time.Now().Local()})

	return http.StatusOK, response
}
