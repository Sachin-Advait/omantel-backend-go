package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetUserSalesAndTargetChart(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	report, err := impl.CacheService.GetUserReport(userID, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return impl.userTarget(ctx, report), nil
	}

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.DbService.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	userID = vps[0]
	// }

	user, err := impl.DbService.FetchUserDetails(userID, quarter, year)
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

		_, report, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(resp)
		// logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Str("report", string(val)).Msg("GetSingleUserTargetChart : updated")

		return impl.userTarget(ctx, report), nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	// val, _ := json.Marshal(mngrReport)
	// logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Str("report", string(val)).Msg("GetSingleUserTargetChart : updated")

	return impl.userTarget(ctx, mngrReport), nil
}

func (impl *UserImpl) userTarget(ctx *gin.Context, userReport *models.UserReportDB) *models.ChartDataResponse {

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
	}
}
