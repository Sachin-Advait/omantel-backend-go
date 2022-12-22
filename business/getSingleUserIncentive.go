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

// GetSingleUserIncentive - run time calculation of user payout
func (impl *IncentiveImpl) GetSingleUserIncentive(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

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

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err := impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserIncentive : precalculated")

			return &models.TotalIncentive{
				HBB: &models.IndividualIncentive{
					Payout:      userReport.HBBPayout,
					ExtraPayout: userReport.HBBExtraPayout,
				},
				Hayyak: &models.IndividualIncentive{
					Payout:      userReport.HayyakPayout,
					ExtraPayout: userReport.HayyakExtraPayout,
				},
				Baqati: &models.IndividualIncentive{
					Payout:      userReport.BaqatiPayout,
					ExtraPayout: userReport.BaqatiExtraPayout,
				},
				Upgrade: &models.IndividualIncentive{
					Payout:      userReport.UpgradePayout,
					ExtraPayout: userReport.UpgradeExtraPayout,
				},
			}, nil
		}

		if err != nil || (err == nil && userReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			resp, _, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserIncentive : updated")

			return resp, nil
		}

		logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetSingleUserIncentive : trace")

		return &models.TotalIncentive{
			HBB: &models.IndividualIncentive{
				Payout:      userReport.HBBPayout,
				ExtraPayout: userReport.HBBExtraPayout,
			},
			Hayyak: &models.IndividualIncentive{
				Payout:      userReport.HayyakPayout,
				ExtraPayout: userReport.HayyakExtraPayout,
			},
			Baqati: &models.IndividualIncentive{
				Payout:      userReport.BaqatiPayout,
				ExtraPayout: userReport.BaqatiExtraPayout,
			},
			Upgrade: &models.IndividualIncentive{
				Payout:      userReport.UpgradePayout,
				ExtraPayout: userReport.UpgradeExtraPayout,
			},
		}, nil
	}

	////

	mngrReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if !needUpdate {
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserIncentive : precalculated")

		return &models.TotalIncentive{
			HBB: &models.IndividualIncentive{
				Payout:      mngrReport.HBBPayout,
				ExtraPayout: mngrReport.HBBExtraPayout,
			},
			Hayyak: &models.IndividualIncentive{
				Payout:      mngrReport.HayyakPayout,
				ExtraPayout: mngrReport.HayyakExtraPayout,
			},
			Baqati: &models.IndividualIncentive{
				Payout:      mngrReport.BaqatiPayout,
				ExtraPayout: mngrReport.BaqatiExtraPayout,
			},
			Upgrade: &models.IndividualIncentive{
				Payout:      mngrReport.UpgradePayout,
				ExtraPayout: mngrReport.UpgradeExtraPayout,
			},
		}, nil
	}

	if err != nil || (err == nil && mngrReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

		mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserIncentive : updated")

		return &models.TotalIncentive{
			HBB: &models.IndividualIncentive{
				Payout:      mngrReport.HBBPayout,
				ExtraPayout: mngrReport.HBBExtraPayout,
			},
			Hayyak: &models.IndividualIncentive{
				Payout:      mngrReport.HayyakPayout,
				ExtraPayout: mngrReport.HayyakExtraPayout,
			},
			Baqati: &models.IndividualIncentive{
				Payout:      mngrReport.BaqatiPayout,
				ExtraPayout: mngrReport.BaqatiExtraPayout,
			},
			Upgrade: &models.IndividualIncentive{
				Payout:      mngrReport.UpgradePayout,
				ExtraPayout: mngrReport.UpgradeExtraPayout,
			},
		}, nil
	}

	logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetSingleUserIncentive : trace")

	return &models.TotalIncentive{
		HBB: &models.IndividualIncentive{
			Payout:      mngrReport.HBBPayout,
			ExtraPayout: mngrReport.HBBExtraPayout,
		},
		Hayyak: &models.IndividualIncentive{
			Payout:      mngrReport.HayyakPayout,
			ExtraPayout: mngrReport.HayyakExtraPayout,
		},
		Baqati: &models.IndividualIncentive{
			Payout:      mngrReport.BaqatiPayout,
			ExtraPayout: mngrReport.BaqatiExtraPayout,
		},
		Upgrade: &models.IndividualIncentive{
			Payout:      mngrReport.UpgradePayout,
			ExtraPayout: mngrReport.UpgradeExtraPayout,
		},
	}, nil
}
