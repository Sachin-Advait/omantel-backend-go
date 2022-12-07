package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	userReport, err := impl.PsqlDb.FetchUserReport(userId, isQuarterReport, quarter, month, year)
	if err == nil && userReport != nil {

		currTime := time.Now()
		currYear, currMonth, currDay := currTime.Date()

		if userReport.UserId != "" && !userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())) { //entry exists and created time is today
			return &models.ChartDataResponse{
				HBB: &models.TargetAchievementData{
					SalesCount:         userReport.Hbb_sales_count,
					Target:             userReport.Hbb_target,
					AchievedPercentage: userReport.Hbb_achieved_percentage,
				},
				Hayyak: &models.TargetAchievementData{
					SalesCount:         userReport.Hayyak_sales_count,
					Target:             userReport.Hayyak_target,
					AchievedPercentage: userReport.Hayyak_achieved_percentage,
				},
				Baqati: &models.TargetAchievementData{
					SalesCount:         userReport.Baqati_sales_count,
					Target:             userReport.Baqati_target,
					AchievedPercentage: userReport.Baqati_achieved_percentage,
				},
				Upgrade: &models.TargetAchievementData{
					SalesCount:         userReport.Upgrade_sales_count,
					Target:             userReport.Upgrade_target,
					AchievedPercentage: userReport.Uprade_achieved_percentage,
				},
			}, nil
		}
	}

	_, userReport, err = impl.updateUserReport(ctx, userId, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	return &models.ChartDataResponse{
		HBB: &models.TargetAchievementData{
			SalesCount:         userReport.Hbb_sales_count,
			Target:             userReport.Hbb_target,
			AchievedPercentage: userReport.Hbb_achieved_percentage,
		},
		Hayyak: &models.TargetAchievementData{
			SalesCount:         userReport.Hayyak_sales_count,
			Target:             userReport.Hayyak_target,
			AchievedPercentage: userReport.Hayyak_achieved_percentage,
		},
		Baqati: &models.TargetAchievementData{
			SalesCount:         userReport.Baqati_sales_count,
			Target:             userReport.Baqati_target,
			AchievedPercentage: userReport.Baqati_achieved_percentage,
		},
		Upgrade: &models.TargetAchievementData{
			SalesCount:         userReport.Upgrade_sales_count,
			Target:             userReport.Upgrade_target,
			AchievedPercentage: userReport.Uprade_achieved_percentage,
		},
	}, nil
}
