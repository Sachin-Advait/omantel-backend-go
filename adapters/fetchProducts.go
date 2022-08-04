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

func FetchProducts(ctx *gin.Context) (int, interface{}) {

	queryParam := ctx.Request.URL.Query()

	products, ok := queryParam[constant.Product]
	if !ok || products[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY PRODUCTS")).Stack().Msg("FetchProducts : invalid params")
		return http.StatusBadRequest, &models.APIResponse{
			ErrorCode: errors.InvalidReqParameterCode,
			Message:   errors.InvalidReqParameters,
			Data:      "EMPTY PRODUCTS",
		}
	}

	productsReq := &models.DeleteProductRequest{}

	if strings.ToLower(products[0]) != constant.All {
		for _, v := range products {
			productsReq.ProductCodes = append(productsReq.ProductCodes, v)
		}
	}

	response, errs := theProvider.FetchProducts(ctx, productsReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("FetchProducts : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchProducts : success")

	return http.StatusOK, response
}
