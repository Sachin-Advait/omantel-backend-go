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

func FetchOutlets(ctx *gin.Context) (int, interface{}) {

	queryParam := ctx.Request.URL.Query()

	outlets, ok := queryParam[constant.Outlet]
	if !ok || outlets[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY OUTLETS")).Stack().Msg("FetchOutlets : invalid params")
		return http.StatusBadRequest, &models.APIResponse{
			ErrorCode: errors.InvalidReqParameterCode,
			Message:   errors.InvalidReqParameters,
			Data:      "EMPTY OUTLETS",
		}
	}

	outletsReq := &models.DeleteOutletRequest{}

	if strings.ToLower(outlets[0]) != constant.All {
		for _, v := range outlets {
			outletsReq.OutletCodes = append(outletsReq.OutletCodes, v)
		}
	}

	response, errs := theProvider.FetchOutlets(ctx, outletsReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("FetchOutlets : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchOutlets : success")

	return http.StatusOK, response
}
