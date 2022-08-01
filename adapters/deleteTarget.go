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

func DeleteTargets(ctx *gin.Context, req []*models.NewTargetRequest) (int, interface{}) {

	if req == nil || len(req) == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("nil ids")).Stack().Msg("DeleteTargets : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid target param",
		}
	}

	targetReq := &models.DeleteTargetRequest{}

	for _, v := range req {

		if v.RegionCode == "" || v.OutletCode == "" || v.Month < 1 || v.Month > 12 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("blank ids")).Stack().Msg("DeleteTargets : invalid req params")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "Invalid target param",
			}
		}

		targetReq.TargetCodes = append(targetReq.TargetCodes, strings.ToUpper(v.RegionCode)+strings.ToUpper(v.OutletCode)+fmt.Sprintf("%d", v.Month))
	}

	errs := theProvider.DeleteTargets(ctx, targetReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("DeleteTargets : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("DeleteTargets : success")

	return http.StatusOK, resp
}
