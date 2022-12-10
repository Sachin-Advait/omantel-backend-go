package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// GetUserPayout - cached response
func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserDetails(userID)
	if err != nil {
		return nil, err
	}

	if _, ok := constant.EmployeePositionMap[user.Position]; !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    userID + "INVALID POSITION FOUND",
			Details:    "invalid position for user payout, userid-" + userID + ", position-" + user.Position,
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

		if err != nil || (err == nil && userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			resp, _, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			return resp, nil
		}

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

	if err != nil || (err == nil && mngrReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

		mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

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
