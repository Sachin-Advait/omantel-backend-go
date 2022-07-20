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

func CreateNewRegion(ctx *gin.Context, req *models.NewRegionRequest) (int, interface{}) {

	if req.RegionCode == "" || req.RegionName == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid req param")).Stack().Msg("CreateNewRegion : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "RegionName or RegionCode is blank",
		}
	}

	errs := theProvider.CreateNewRegion(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewRegion : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewRegion : success")

	return http.StatusOK, resp
}
