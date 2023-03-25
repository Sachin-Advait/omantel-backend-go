package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchNotification(ctx *gin.Context, category, id string) (interface{}, *errors.Error) {

	response := make(map[string][]*models.Notification)

	if id != "" {
		dbResponse, err := impl.DbService.FetchSingleNotificationRecord(id)
		if err != nil {
			return nil, err
		}

		return dbResponse, nil
	}

	if category != "" {
		dbResponse, err := impl.DbService.FetchNotificationRecord(category)
		if err != nil {
			return nil, err
		}

		response[category] = append(response[category], impl.notificationResponse(dbResponse)...)

		return response, nil
	}

	for key := range constant.NotificationTypeMap {
		dbResponse, err := impl.DbService.FetchNotificationRecord(key)
		if err != nil {
			return nil, err
		}

		response[key] = append(response[key], impl.notificationResponse(dbResponse)...)
	}

	return response, nil
}

func (impl *AdminImpl) notificationResponse(dbResponse []*models.NotificationDB) []*models.Notification {

	response := make([]*models.Notification, 0)

	for _, val := range dbResponse {
		response = append(response, &models.Notification{
			Id:         val.Id,
			Category:   val.Category,
			Subject:    val.Subject,
			Body:       val.Body,
			FileType:   val.FileType,
			FileLink:   val.FileLink,
			CreatedAt:  val.CreatedAt,
			ModifiedAt: val.ModifiedAt,
		})
	}

	return response
}
