package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetUserSalesAndTargetChart(ctx *gin.Context, userID string, chanel string,
	isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

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

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb || user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	report, err := impl.CacheService.GetUserReport(userID, chanel, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return impl.userTarget(ctx, report), nil
	}

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
		_, report, err := impl.updateAgentUserReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return impl.userTarget(ctx, report), nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}
	return impl.userTarget(ctx, mngrReport), nil
}

func (impl *UserImpl) GetUserSalesAndTargetChartPerPlan(ctx *gin.Context, userID string, chanel string,
	planType string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

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

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb || user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	report, err := impl.CacheService.GetUserReportPerCategory(userID, chanel, planType, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return impl.userTarget(ctx, report), nil
	}

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
		_, report, err := impl.updateAgentUserReportPerPlanType(ctx, user, chanel, planType, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return impl.userTarget(ctx, report), nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReportPerPlanType(ctx, user, chanel, planType, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}
	return impl.userTarget(ctx, mngrReport), nil
}

func (impl *UserImpl) userTarget(_ *gin.Context, userReport *models.UserReportDB) *models.ChartDataResponse {

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
