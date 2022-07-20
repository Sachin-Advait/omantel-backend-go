package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"time"
)

func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, id, quarter int, isQuarterly bool, t time.Time) (*psqldb.TotalIncentive, *errors.Error) {

	usrRcrd, err := impl.PsqlDb.FetchIndividualUserPayoutReport(id, int(t.Month()), quarter, isQuarterly)
	if err != nil {
		return nil, err
	}

	resp := &psqldb.TotalIncentive{
		HBB: &psqldb.IndividualIncentive{
			Payout:      usrRcrd.Hbb_payout,
			ExtraPayout: usrRcrd.Hbb_extra_payout,
		},
		Hayyak: &psqldb.IndividualIncentive{
			Payout:      usrRcrd.Hayyak_payout,
			ExtraPayout: usrRcrd.Hayyak_extra_payout,
		},
		Baqati: &psqldb.IndividualIncentive{
			Payout:      usrRcrd.Baqati_payout,
			ExtraPayout: usrRcrd.Baqati_extra_payout,
		},
		Upgrade: &psqldb.IndividualIncentive{
			Payout:      usrRcrd.Upgrade_payout,
			ExtraPayout: usrRcrd.Upgrade_extra_payout,
		},
		TotalPay: usrRcrd.Total_payout,
	}

	return resp, nil
}
