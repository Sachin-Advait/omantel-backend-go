package adapters

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func FetchRegions(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	queryParam := ctx.Request.URL.Query()

	regions, ok := queryParam[constant.Region]
	if !ok || regions[0] == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New("EMPTY REGIONS")).Stack().Msg("FetchRegions : invalid params")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "invalid req param",
			CreateTime: time.Now().Local(),
		}

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
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", response)).Msg("FetchRegions : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", response),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, response
}
