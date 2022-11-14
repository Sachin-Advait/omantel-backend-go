package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserTarget(ctx *gin.Context, id string, quarter int, isQuarterly bool, t time.Time) (*models.ChartDataResponse, *errors.Error) {

	usrRcrd, err := impl.PsqlDb.FetchIndividualUserTargetReport(id, int(t.Month()), quarter, isQuarterly)
	if err != nil {
		return nil, err
	}

	resp := &models.ChartDataResponse{
		HBB: &models.TargetAchievementData{
			SalesCount:         usrRcrd.Hbb_sales_count,
			Target:             usrRcrd.Hbb_target,
			AchievedPercentage: usrRcrd.Hbb_achieved_percentage,
		},
		Hayyak: &models.TargetAchievementData{
			SalesCount:         usrRcrd.Hayyak_sales_count,
			Target:             usrRcrd.Hayyak_target,
			AchievedPercentage: usrRcrd.Hayyak_achieved_percentage,
		},
		Baqati: &models.TargetAchievementData{
			SalesCount:         usrRcrd.Baqati_sales_count,
			Target:             usrRcrd.Baqati_target,
			AchievedPercentage: usrRcrd.Baqati_achieved_percentage,
		},
		Upgrade: &models.TargetAchievementData{
			SalesCount:         usrRcrd.Upgrade_sales_count,
			Target:             usrRcrd.Upgrade_target,
			AchievedPercentage: usrRcrd.Uprade_achieved_percentage,
		},
	}

	return resp, nil
}
