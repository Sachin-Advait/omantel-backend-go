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

func DeleteTargets(ctx *gin.Context, req []*models.NewTargetRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if req == nil || len(req) == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("nil ids")).Msg("DeleteTargets : invalid req params")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Text:       "invalid req param",
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: "Invalid target param",
		}
	}

	targetReq := &models.DeleteTargetRequest{}

	for _, v := range req {

		if v.RegionCode == "" || v.OutletCode == "" || v.Month < 1 || v.Month > 12 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("blank ids")).Msg("DeleteTargets : invalid req params")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "blank ids",
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "Invalid target param",
			}
		}

		targetReq.TargetCodes = append(targetReq.TargetCodes, strings.ToUpper(v.RegionCode)+strings.ToUpper(v.OutletCode)+fmt.Sprintf("%d", v.Month))
	}

	errs := theProvider.DeleteTargets(ctx, targetReq)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("DeleteTargets : provider failed")
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

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("DeleteTargets : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", resp),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, resp
}
