package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserReport(ctx *gin.Context, clientID, month, quarter int, isQuarterReport bool) (*models.UserRecord, *errors.Error) {

	return impl.PsqlDb.FetchIndividualUserOverAllReport(clientID, month, quarter, isQuarterReport)
}
