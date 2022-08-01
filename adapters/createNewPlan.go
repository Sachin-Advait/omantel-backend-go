package adapters

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CreateNewPlan(ctx *gin.Context, req *models.NewPlanRequest) (int, interface{}) {

	if req.PlanCode == "" || req.PlanName == "" || req.ProductCode == "" || req.PlanPrice == 0 || req.SIPPlanRate == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid plan param")).Stack().Msg("CreateNewPlan : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid plan params",
		}
	}

	errs := theProvider.CreateNewPlan(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewPlan : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewPlan : success")

	return http.StatusOK, resp
}
