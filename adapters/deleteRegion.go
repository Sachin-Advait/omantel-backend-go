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

func DeleteRegions(ctx *gin.Context, req *models.DeleteRegionRequest) (int, interface{}) {

	if req == nil || len(req.RegionCodes) == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("nil ids")).Stack().Msg("DeleteRegions : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid target param",
		}
	}

	for i, v := range req.RegionCodes {

		if v == "" {
			logger.FileLogger.CtxError(ctx).Err(errors.New("blank ids")).Stack().Msg("DeleteRegions : invalid req params")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "Invalid target param",
			}
		}

		req.RegionCodes[i] = v
	}

	errs := theProvider.DeleteRegions(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("DeleteRegions : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("DeleteRegions : success")

	return http.StatusOK, resp
}
