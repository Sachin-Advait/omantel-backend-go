package adapters

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CreateNewOutlet(ctx *gin.Context, req *models.NewOutletRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if req.OutletCode == "" || req.OutletName == "" || req.RegionCode == "" || req.Lat == 0 || req.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Err(errors.New("invalid req param")).Stack().Msg("CreateNewOutlet : invalid req params")
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
			Message: "Blank value found",
		}
	}

	errs := theProvider.CreateNewOutlet(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewOutlet : provider failed")
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

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewOutlet : success")
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
