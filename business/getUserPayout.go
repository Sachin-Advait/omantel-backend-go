package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// GetUserPayout - cached response
func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	//TODO: check for normal employee

	payoutReport, err := impl.PsqlDb.FetchUserReport(userId, isQuarterReport, quarter, month, year)
	if err == nil && payoutReport != nil {

		currTime := time.Now()
		currYear, currMonth, currDay := currTime.Date()

		if payoutReport.UserId != "" && !payoutReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())) { //entry exists and created time is today
			return &models.TotalIncentive{
				HBB: &models.IndividualIncentive{
					Payout:      payoutReport.Hbb_payout,
					ExtraPayout: payoutReport.Hbb_extra_payout,
				},
				Baqati: &models.IndividualIncentive{
					Payout:      payoutReport.Baqati_payout,
					ExtraPayout: payoutReport.Baqati_extra_payout,
				},
				Hayyak: &models.IndividualIncentive{
					Payout:      payoutReport.Hayyak_payout,
					ExtraPayout: payoutReport.Hayyak_extra_payout,
				},
				Upgrade: &models.IndividualIncentive{
					Payout:      payoutReport.Upgrade_payout,
					ExtraPayout: payoutReport.Upgrade_extra_payout,
				},
			}, nil
		}
	}

	resp, _, err := impl.updateUserReport(ctx, userId, isQuarterReport, quarter, month, year)
	return resp, err
}
