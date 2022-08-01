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

func DeleteProducts(ctx *gin.Context, req *models.DeleteProductRequest) (int, interface{}) {

	if req == nil || len(req.ProductCodes) == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("nil ids")).Stack().Msg("DeleteProducts : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid target param",
		}
	}

	for i, v := range req.ProductCodes {

		if v == "" {
			logger.FileLogger.CtxError(ctx).Err(errors.New("blank ids")).Stack().Msg("DeleteProducts : invalid req params")
			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "Invalid target param",
			}
		}

		req.ProductCodes[i] = v
	}

	errs := theProvider.DeleteProducts(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("DeleteProducts : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("DeleteProducts : success")

	return http.StatusOK, resp
}
