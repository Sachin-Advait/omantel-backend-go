package impl

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) UpsertNotification(ctx *gin.Context, req *models.Notification) (*models.APISuccessResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	dbReq := &models.NotificationDB{
		Id:         req.Id,
		Category:   req.Category,
		Subject:    req.Subject,
		Body:       req.Body,
		FileType:   req.FileType,
		FileLink:   req.FileLink,
		ModifiedAt: time.Now(),
		IsUpdate:   true,
	}

	if req.Id == "" {
		dbReq.IsUpdate = false
		dbReq.Id = uuid.NewString()
		dbReq.CreatedAt = dbReq.ModifiedAt
	}

	resp, err := impl.DbService.UpsertNotificationRecord(ctxHeader.XUserId, dbReq)
	if err != nil {
		return nil, err
	}

	return &models.APISuccessResponse{
		Status: true,
		Data: map[string]string{
			"id": resp,
		},
	}, nil
}
