package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

// GetUserPayout - cached response
func (impl *IncentiveImpl) GetUserPayout(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	report, err := impl.RedisCache.GetUserReport(userID, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return impl.userPayout(ctx, report), nil
	}

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.PsqlDb.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	userID = vps[0]
	// }

	user, err := impl.PsqlDb.FetchUserDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	position, ok := constant.EmployeePositionMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
			Details:    "incorrect user position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {

		resp, _, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(resp)
		// logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Str("report", string(val)).Msg("GetUserPayout : updated")

		return resp, nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	// val, _ := json.Marshal(mngrReport)
	// logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Str("report", string(val)).Msg("GetUserPayout : updated")

	return impl.userPayout(ctx, mngrReport), nil
}

func (impl *IncentiveImpl) userPayout(ctx *gin.Context, userReport *models.UserReportDB) *models.TotalIncentive {

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
	}
}
