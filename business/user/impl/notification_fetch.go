package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) FetchNotification(ctx *gin.Context, category, id string) (interface{}, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	response := make([]*models.Notification, 0)

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

		response = append(response, impl.notificationResponse(dbResponse)...)

		return response, nil
	}

	for key := range constant.NotificationTypeMap {
		dbResponse, err := impl.DbService.FetchNotificationRecord(key)
		if err != nil {
			return nil, err
		}

		response = append(response, impl.notificationResponse(dbResponse)...)
	}

	return response, nil
}

func (impl *UserImpl) notificationResponse(dbResponse []*models.NotificationDB) []*models.Notification {

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
