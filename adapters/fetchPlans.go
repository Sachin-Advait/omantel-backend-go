package adapters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func FetchPlans(ctx *gin.Context) (int, interface{}) {

	queryParam := ctx.Request.URL.Query()

	plans, ok := queryParam[constant.Plan]
	if !ok || plans[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY PLANS")).Stack().Msg("FetchPlans : invalid params")
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
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchPlans : success")

	return http.StatusOK, response
}
