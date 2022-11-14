package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, id string, quarter int, isQuarterly bool, t time.Time) (*models.TotalIncentive, *errors.Error) {

	usrRcrd, err := impl.PsqlDb.FetchIndividualUserPayoutReport(id, int(t.Month()), quarter, isQuarterly)
	if err != nil {
		return nil, err
	}

	resp := &models.TotalIncentive{
		HBB: &models.IndividualIncentive{
			Payout:      usrRcrd.Hbb_payout,
			ExtraPayout: usrRcrd.Hbb_extra_payout,
		},
		Hayyak: &models.IndividualIncentive{
			Payout:      usrRcrd.Hayyak_payout,
			ExtraPayout: usrRcrd.Hayyak_extra_payout,
		},
		Baqati: &models.IndividualIncentive{
			Payout:      usrRcrd.Baqati_payout,
			ExtraPayout: usrRcrd.Baqati_extra_payout,
		},
		Upgrade: &models.IndividualIncentive{
			Payout:      usrRcrd.Upgrade_payout,
			ExtraPayout: usrRcrd.Upgrade_extra_payout,
		},
		TotalPay: usrRcrd.Total_payout,
	}

	return resp, nil
}
