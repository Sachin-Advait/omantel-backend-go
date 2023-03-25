package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func FetchNotifications(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	str := ctx.Request.URL.Query().Get("category")
	id := ctx.Request.URL.Query().Get("id")

	category := strings.ToUpper(strings.TrimSpace(str))
	_, ok := constant.NotificationTypeMap[category]
	if category != "" && !ok {
		err := errors.New("invalid category query param")
		logger.FileLogger.CtxError(ctx).Err(err).Msg("FetchNotifications : validation failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid req params",
			Code:       http.StatusBadRequest,
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	response, errs := adminAPIProvider.FetchNotification(ctx, category, id)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("FetchNotifications : provider failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "provider failed",
			Code:       errs.StatusCode,
			Text:       errs.Details.(string),
			CreateTime: time.Now().Local(),
		}

		return errs.StatusCode, &models.APIFailureResponse{
			Error: errs.Message,
		}
	}

	return http.StatusOK, response
}
