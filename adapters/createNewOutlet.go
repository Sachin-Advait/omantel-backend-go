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

func CreateNewOutlet(ctx *gin.Context, req *models.NewOutletRequest) (int, interface{}) {

	if req.OutletCode == "" || req.OutletName == "" || req.RegionCode == "" || req.Lat == 0 || req.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid req param")).Stack().Msg("CreateNewOutlet : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Blank value found",
		}
	}

	errs := theProvider.CreateNewOutlet(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewOutlet : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewOutlet : success")

	return http.StatusOK, resp
}
