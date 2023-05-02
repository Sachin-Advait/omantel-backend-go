package user

import (
	"encoding/json"
	"fmt"
	"math"
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

func GetRegionSalesAndTarget(ctx *gin.Context, req *models.UserIncentiveRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	isQuarterReport, quarter, month, year, err := models.ValidateUserIncentiveRequest(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("GetRegionSalesAndTarget : validation failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	response, errs := userAPIProvider.GetRegionSalesAndTarget(ctx, req.ClientId, isQuarterReport, quarter, month, year)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("GetRegionSalesAndTarget : provider failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: errs.Details.(string), CreateTime: time.Now().Local()})

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	for _, val := range response {
		val.BaqatiAchievedPercentage = math.Round(val.BaqatiAchievedPercentage*1000) / 1000
		val.HBBAchievedPercentage = math.Round(val.HBBAchievedPercentage*1000) / 1000
		val.HayyakAchievedPercentage = math.Round(val.HayyakAchievedPercentage*1000) / 1000
		val.UpgradeAchievedPercentage = math.Round(val.UpgradeAchievedPercentage*1000) / 1000
		val.TotalAchievedPercentage = math.Round(val.TotalAchievedPercentage*1000) / 1000
		resp = append(resp, val)
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].TotalSales > resp[j].TotalSales
	})

	responseBytes, _ := json.Marshal(response)
	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("GetRegionSalesAndTarget : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success", Code: http.StatusOK,
		Text: string(responseBytes), CreateTime: time.Now().Local()})

	return http.StatusOK, resp
}
