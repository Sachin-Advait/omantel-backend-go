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

func CreateNewProduct(ctx *gin.Context, req *models.NewProductRequest) (int, interface{}) {

	if req.ProductCode == "" || req.ProductName == "" || req.Weightage == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid product param")).Stack().Msg("CreateNewProduct : invalid req params")
		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid product params",
		}
	}

	errs := theProvider.CreateNewProduct(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewProduct : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewProduct : success")

	return http.StatusOK, resp
}
