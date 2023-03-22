package impl

import (
	"fmt"
	"net/http"
	"time"

	"firebase.google.com/go/messaging"
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

	var (
		topic string
		notif *messaging.Notification
		data  map[string]string
	)

	switch req.Category {
	case "OFFER":
		topic, notif, data = impl.getOfferFCMNotif(req.FileType, req.FileLink)
	case "TRAINING":
		topic, notif, data = impl.getTrainingFCMNotif(req.FileType, req.FileLink)
	case "NOTIFICATION":
		topic, notif, data = impl.getNotificationFCMNotif(req.FileType, req.FileLink)
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Err:        errors.InternalServerError,
			Message:    "INVALID NOTIFICATION CATEGORY",
			Details:    "invalid category-" + req.Category,
		}
	}

	fcmResp, fcmErr := impl.FCMService.SendViaTopicMessage(topic, notif, data)
	if fcmErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Err:        errors.InternalServerError,
			Message:    "FAILED FCM MESSAGE",
			Details:    "fcmService: " + fcmErr.Error(),
		}
	}

	fmt.Println("fcmResponse", fcmResp)

	return &models.APISuccessResponse{
		Status: true,
		Data: map[string]string{
			"id": resp,
		},
	}, nil
}

func (impl *AdminImpl) getOfferFCMNotif(fileType, fileLink string) (topic string, notif *messaging.Notification, data map[string]string) {
	topic = impl.Config.FCMOfferTopic
	notif = &messaging.Notification{
		Title:    "A NEW OFFER",
		Body:     "A new limited offer available. Do checkout",
		ImageURL: "",
	}
	data = map[string]string{
		"start": "from today 12:00AM",
		"file":  fileType,
		"link":  fileLink,
	}

	return
}

func (impl *AdminImpl) getNotificationFCMNotif(fileType, fileLink string) (topic string, notif *messaging.Notification, data map[string]string) {
	topic = impl.Config.FCMNotificationTopic
	notif = &messaging.Notification{
		Title:    "A NEW NOTIFICATION",
		Body:     "A new notification found. Do checkout",
		ImageURL: "",
	}
	data = map[string]string{
		"file": fileType,
		"link": fileLink,
	}

	return
}

func (impl *AdminImpl) getTrainingFCMNotif(fileType, fileLink string) (topic string, notif *messaging.Notification, data map[string]string) {
	topic = impl.Config.FCMTrainingTopic
	notif = &messaging.Notification{
		Title:    "NEW TRAINING VIDEO",
		Body:     "A new training video uploaded. Buckle up",
		ImageURL: "",
	}
	data = map[string]string{
		"file": fileType,
		"link": fileLink,
	}

	return
}
