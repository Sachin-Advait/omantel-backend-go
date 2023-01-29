package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

// GetUserPayout - cached response
func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)
	if ctxHeader.IsAdmin {
		vps, err := impl.PsqlDb.FetchVP(quarter, year)
		if err != nil {
			return nil, err
		}
		userID = vps[0]
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

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserPayout : precalculated")

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
				TotalPay: userReport.TotalPayout,
			}, nil
		}

		if err != nil || (err == nil && userReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			resp, _, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserPayout : updated")

			return resp, nil
		}

		logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserPayout : trace")

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
			TotalPay: userReport.TotalPayout,
		}, nil
	}

	////

	mngrReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if !needUpdate {
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserPayout : precalculated")

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
			TotalPay: mngrReport.TotalPayout,
		}, nil
	}

	if err != nil || (err == nil && mngrReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

		mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserPayout : updated")

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
			TotalPay: mngrReport.TotalPayout,
		}, nil
	}

	logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserPayout : trace")

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
		TotalPay: mngrReport.TotalPayout,
	}, nil
}
