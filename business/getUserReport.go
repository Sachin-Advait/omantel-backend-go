package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

func (impl *IncentiveImpl) GetUserReport(ctx *gin.Context, clientID, month, quarter int, isQuarterReport bool) (*psqldb.UserRecord, *errors.Error) {

	return impl.PsqlDb.FetchIndividualUserOverAllReport(clientID, month, quarter, isQuarterReport)
}
