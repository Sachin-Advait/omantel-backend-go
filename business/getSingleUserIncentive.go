package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// GetSingleUserIncentive - run time calculation of user payout
func (impl *IncentiveImpl) GetSingleUserIncentive(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserDetails(userID)
	if err != nil {
		return nil, err
	}

	if user.Position == constant.Promoter || user.Position == constant.CSR {
		resp, _, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	hbb := &models.IndividualIncentive{
		Payout:      mngrReport.HBBPayout,
		ExtraPayout: mngrReport.HBBExtraPayout,
	}

	hayyak := &models.IndividualIncentive{
		Payout:      mngrReport.HayyakPayout,
		ExtraPayout: mngrReport.HayyakExtraPayout,
	}

	baqati := &models.IndividualIncentive{
		Payout:      mngrReport.BaqatiPayout,
		ExtraPayout: mngrReport.BaqatiExtraPayout,
	}

	upgdare := &models.IndividualIncentive{
		Payout:      mngrReport.UpgradePayout,
		ExtraPayout: mngrReport.UpgradeExtraPayout,
	}

	resp := &models.TotalIncentive{
		HBB:     hbb,
		Hayyak:  hayyak,
		Baqati:  baqati,
		Upgrade: upgdare,
	}

	return resp, err
}
