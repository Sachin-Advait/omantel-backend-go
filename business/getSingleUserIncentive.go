package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// GetSingleUserIncentive - run time calculation of user payout
func (impl *IncentiveImpl) GetSingleUserIncentive(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	resp, _, err := impl.updateUserReport(ctx, userId, isQuarterReport, quarter, month, year)

	return resp, err
}
