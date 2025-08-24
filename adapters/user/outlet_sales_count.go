package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
)

func GetOutletTotalSalesCount(ctx *gin.Context, req *models.UserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := models.ValidateUserIncentiveRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetOutletTotalSalesCount : validation failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(),
			Type: "error", Msg: "invalid req params", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	response, errs := userAPIProvider.GetOutletTotalSalesCount(ctx, req.ClientId, req.Channel, isQuarterReport,
		quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).
			Msg("GetOutletTotalSalesCount : provider failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(),
			Type: "error", Msg: "provider failed",
			Text: errs.Details.(string), CreateTime: time.Now().Local()})

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Total > response[j].Total
	})

	responseBytes, _ := json.Marshal(response)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", string(responseBytes))).
		Msg("GetOutletTotalSalesCount : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(),
		Type: "info", Msg: "success", Code: http.StatusOK,
		Text: string(responseBytes), CreateTime: time.Now().Local()})

	return http.StatusOK, response
}
