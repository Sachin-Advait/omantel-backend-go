package adapters

import (
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

func UpsertNotifications(ctx *gin.Context, req *models.Notification) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	err := validateUpsertNotification(req, ctxHeader)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("UpsertNotifications : validation failed")
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

	response, errs := theProvider.UpsertNotification(ctx, req)
	if errs != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Msg("UpsertNotifications : provider failed")
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

func validateUpsertNotification(req *models.Notification, ctxHeader context.ContextHeader) error {

	/*
		1. Latest Offer (Subject, Body, PDF, Excel as attachment)
		2. Notifications (Subject , Body, Attachment as PDF)
		3. Training & learning (Subject , Attachment as video or pdf)
	*/

	req.Category = strings.ToUpper(strings.TrimSpace(req.Category))
	if _, ok := constant.NotificationTypeMap[req.Category]; !ok {
		return errors.New("invalid category")
	}

	if req.Subject == "" {
		return errors.New("invalid subject")
	}

	req.FileType = strings.ToUpper(strings.TrimSpace(req.FileType))
	if _, ok := constant.SupportedFileTypeMap[req.Category]; !ok {
		return errors.New("invalid file type")
	}

	if req.FileLink == "" {
		return errors.New("invalid file link")
	}

	return nil
}
