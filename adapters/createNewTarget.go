package adapters

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	"net/http"
)

func CreateNewTarget(ctx *gin.Context, req *models.NewTargetRequest) (int, interface{}) {

	if req.RegionName == "" || req.OutletName == "" || req.Month < 1 || req.Month > 12 || req.Baqati == 0 || req.HBB == 0 || req.Hayyak == 0 || req.Upgrade == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid target param")).Stack().Msg("CreateNewTarget : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid target param",
		}
	}

	errs := theProvider.CreateNewTarget(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewTarget : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewTarget : success")

	return http.StatusOK, resp
}
