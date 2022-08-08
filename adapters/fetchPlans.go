package adapters

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func FetchPlans(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	queryParam := ctx.Request.URL.Query()

	plans, ok := queryParam[constant.Plan]
	if !ok || plans[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY PLANS")).Stack().Msg("FetchPlans : invalid params")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "invalid req param",
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIResponse{
			ErrorCode: errors.InvalidReqParameterCode,
			Message:   errors.InvalidReqParameters,
			Data:      "EMPTY PLANS",
		}
	}

	plansReq := &models.DeletePlanRequest{}

	if strings.ToLower(plans[0]) != constant.All {
		for _, v := range plans {
			plansReq.PlanCodes = append(plansReq.PlanCodes, v)
		}
	}

	response, errs := theProvider.FetchPlans(ctx, plansReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("FetchPlans : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchPlans : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", response),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
