package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserDetails(userId)
	if err != nil {
		return nil, err
	}

	var needUpdate bool

	currTime := time.Now()
	currYear, currMonth, currDay := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if isQuarterReport && year == currYear && quarter <= currQuarter {
		needUpdate = true
	} else if !isQuarterReport && year == currYear && month <= int(currMonth) {
		needUpdate = true
	}

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err := impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			_, userReport, err = impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
		}

		return &models.ChartDataResponse{
			HBB: &models.TargetAchievementData{
				SalesCount:         userReport.HBBSalesCount,
				Target:             userReport.HBBTargetCount,
				AchievedPercentage: userReport.HBBAchievedPercentage,
			},
			Hayyak: &models.TargetAchievementData{
				SalesCount:         userReport.HayyakSalesCount,
				Target:             userReport.HayyakTargetCount,
				AchievedPercentage: userReport.HayyakAchievedPercentage,
			},
			Baqati: &models.TargetAchievementData{
				SalesCount:         userReport.BaqatiSalesCount,
				Target:             userReport.BaqatiTargetCount,
				AchievedPercentage: userReport.BaqatiAchievedPercentage,
			},
			Upgrade: &models.TargetAchievementData{
				SalesCount:         userReport.UpgradeSalesCount,
				Target:             userReport.UpgradeTargetCount,
				AchievedPercentage: userReport.UpradeAchievedPercentage,
			},
		}, nil
	}

	userReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if needUpdate && (err != nil || (err == nil && userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
		userReport, err = impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	}

	if err != nil {
		return nil, err
	}

	return &models.ChartDataResponse{
		HBB: &models.TargetAchievementData{
			SalesCount:         userReport.HBBSalesCount,
			Target:             userReport.HBBTargetCount,
			AchievedPercentage: userReport.HBBAchievedPercentage,
		},
		Hayyak: &models.TargetAchievementData{
			SalesCount:         userReport.HayyakSalesCount,
			Target:             userReport.HayyakTargetCount,
			AchievedPercentage: userReport.HayyakAchievedPercentage,
		},
		Baqati: &models.TargetAchievementData{
			SalesCount:         userReport.BaqatiSalesCount,
			Target:             userReport.BaqatiTargetCount,
			AchievedPercentage: userReport.BaqatiAchievedPercentage,
		},
		Upgrade: &models.TargetAchievementData{
			SalesCount:         userReport.UpgradeSalesCount,
			Target:             userReport.UpgradeTargetCount,
			AchievedPercentage: userReport.UpradeAchievedPercentage,
		},
	}, nil
}
