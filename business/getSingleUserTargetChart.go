package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetSingleUserTargetChart(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	user, err := impl.PsqlDb.FetchUserDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	_, ok := constant.EmployeePositionMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
			Details:    "incorrect user position: userId-" + user.UserId + " position-" + user.Position,
		}
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

	//TODO: remove at last
	needUpdate = true

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err := impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserTargetChart : precalculated")

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

		if err != nil || (err == nil && userReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			_, resp, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserTargetChart : updated")

			return &models.ChartDataResponse{
				HBB: &models.TargetAchievementData{
					SalesCount:         resp.HBBSalesCount,
					Target:             resp.HBBTargetCount,
					AchievedPercentage: resp.HBBAchievedPercentage,
				},
				Hayyak: &models.TargetAchievementData{
					SalesCount:         resp.HayyakSalesCount,
					Target:             resp.HayyakTargetCount,
					AchievedPercentage: resp.HayyakAchievedPercentage,
				},
				Baqati: &models.TargetAchievementData{
					SalesCount:         resp.BaqatiSalesCount,
					Target:             resp.BaqatiTargetCount,
					AchievedPercentage: resp.BaqatiAchievedPercentage,
				},
				Upgrade: &models.TargetAchievementData{
					SalesCount:         resp.UpgradeSalesCount,
					Target:             resp.UpgradeTargetCount,
					AchievedPercentage: resp.UpradeAchievedPercentage,
				},
			}, nil
		}

		logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserTargetChart : trace")

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

	////////

	mngrReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if !needUpdate {
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserTargetChart : precalculated")

		return &models.ChartDataResponse{
			HBB: &models.TargetAchievementData{
				SalesCount:         mngrReport.HBBSalesCount,
				Target:             mngrReport.HBBTargetCount,
				AchievedPercentage: mngrReport.HBBAchievedPercentage,
			},
			Hayyak: &models.TargetAchievementData{
				SalesCount:         mngrReport.HayyakSalesCount,
				Target:             mngrReport.HayyakTargetCount,
				AchievedPercentage: mngrReport.HayyakAchievedPercentage,
			},
			Baqati: &models.TargetAchievementData{
				SalesCount:         mngrReport.BaqatiSalesCount,
				Target:             mngrReport.BaqatiTargetCount,
				AchievedPercentage: mngrReport.BaqatiAchievedPercentage,
			},
			Upgrade: &models.TargetAchievementData{
				SalesCount:         mngrReport.UpgradeSalesCount,
				Target:             mngrReport.UpgradeTargetCount,
				AchievedPercentage: mngrReport.UpradeAchievedPercentage,
			},
		}, nil
	}

	if err != nil || (err == nil && mngrReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

		resp, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserTargetChart : updated")

		return &models.ChartDataResponse{
			HBB: &models.TargetAchievementData{
				SalesCount:         resp.HBBSalesCount,
				Target:             resp.HBBTargetCount,
				AchievedPercentage: resp.HBBAchievedPercentage,
			},
			Hayyak: &models.TargetAchievementData{
				SalesCount:         resp.HayyakSalesCount,
				Target:             resp.HayyakTargetCount,
				AchievedPercentage: resp.HayyakAchievedPercentage,
			},
			Baqati: &models.TargetAchievementData{
				SalesCount:         resp.BaqatiSalesCount,
				Target:             resp.BaqatiTargetCount,
				AchievedPercentage: resp.BaqatiAchievedPercentage,
			},
			Upgrade: &models.TargetAchievementData{
				SalesCount:         resp.UpgradeSalesCount,
				Target:             resp.UpgradeTargetCount,
				AchievedPercentage: resp.UpradeAchievedPercentage,
			},
		}, nil
	}

	logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserTargetChart : trace")

	return &models.ChartDataResponse{
		HBB: &models.TargetAchievementData{
			SalesCount:         mngrReport.HBBSalesCount,
			Target:             mngrReport.HBBTargetCount,
			AchievedPercentage: mngrReport.HBBAchievedPercentage,
		},
		Hayyak: &models.TargetAchievementData{
			SalesCount:         mngrReport.HayyakSalesCount,
			Target:             mngrReport.HayyakTargetCount,
			AchievedPercentage: mngrReport.HayyakAchievedPercentage,
		},
		Baqati: &models.TargetAchievementData{
			SalesCount:         mngrReport.BaqatiSalesCount,
			Target:             mngrReport.BaqatiTargetCount,
			AchievedPercentage: mngrReport.BaqatiAchievedPercentage,
		},
		Upgrade: &models.TargetAchievementData{
			SalesCount:         mngrReport.UpgradeSalesCount,
			Target:             mngrReport.UpgradeTargetCount,
			AchievedPercentage: mngrReport.UpradeAchievedPercentage,
		},
	}, nil
}
