package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) FetchNotification(ctx *gin.Context, category string) ([]*models.Notification, *errors.Error) {

	dbResponse, err := impl.DbService.FetchNotificationRecord(category)
	if err != nil {
		return nil, err
	}

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

	return response, nil
}
