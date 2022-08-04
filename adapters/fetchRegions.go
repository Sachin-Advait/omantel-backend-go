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

func FetchRegions(ctx *gin.Context) (int, interface{}) {

	queryParam := ctx.Request.URL.Query()

	regions, ok := queryParam[constant.Region]
	if !ok || regions[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY REGIONS")).Stack().Msg("FetchRegions : invalid params")
		return http.StatusBadRequest, &models.APIResponse{
			ErrorCode: errors.InvalidReqParameterCode,
			Message:   errors.InvalidReqParameters,
			Data:      "EMPTY REGIONS",
		}
	}

	regionsReq := &models.DeleteRegionRequest{}

	if strings.ToLower(regions[0]) != constant.All {
		for _, v := range regions {
			regionsReq.RegionCodes = append(regionsReq.RegionCodes, v)
		}
	}

	response, errs := theProvider.FetchRegions(ctx, regionsReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("FetchRegions : provider failed")
		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchRegions : success")

	return http.StatusOK, response
}
