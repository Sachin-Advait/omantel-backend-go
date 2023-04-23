package impl

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

// GetUserPayout - cached response
func (impl *UserImpl) GetUserPayout(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

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

	if !isQuarterReport || quarter != 4 {
		return impl.CoalesceUserPayout(ctx, userID, isQuarterReport, quarter, month, year)
	}

	q1Report, err := impl.CoalesceUserPayout(ctx, userID, isQuarterReport, 1, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-1, err-%v", err)
		return nil, err
	}

	q2Report, err := impl.CoalesceUserPayout(ctx, userID, isQuarterReport, 2, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-2, err-%v", err)
		return nil, err
	}

	q3Report, err := impl.CoalesceUserPayout(ctx, userID, isQuarterReport, 3, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-3, err-%v", err)
		return nil, err
	}

	q4Report, err := impl.CoalesceUserPayout(ctx, userID, isQuarterReport, 4, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-4, err-%v", err)
		return nil, err
	}

	vocCriteria, vocErr := impl.getVOCCriteria(year, q4Report.VOC.ExtraPayout)
	if vocErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    vocErr.Error(),
			Details:    fmt.Sprintf("vocCriteria: year-%v, err-%v", year, vocErr),
		}
	}

	var achievedCount int

	if q1Report.TotalPay != 0 {
		achievedCount++
	}
	if q2Report.TotalPay != 0 {
		achievedCount++
	}
	if q3Report.TotalPay != 0 {
		achievedCount++
	}
	if q4Report.TotalPay != 0 {
		achievedCount++
	}

	if achievedCount < 3 {
		q4Report.VOC.Payout = 0
	} else {
		totalPay := q1Report.TotalPay + q2Report.TotalPay + q3Report.TotalPay + q4Report.TotalPay
		q4Report.VOC.Payout = (vocCriteria * totalPay) / 100
		q4Report.VOC.Payout = math.Round(q4Report.VOC.Payout*1000) / 1000
		q4Report.TotalPay = q4Report.TotalPay + q4Report.VOC.Payout
	}

	return q4Report, nil
}

func (impl *UserImpl) CoalesceUserPayout(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	report, err := impl.CacheService.GetUserReport(userID, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return impl.userPayout(ctx, report), nil
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

func (impl *UserImpl) userPayout(ctx *gin.Context, userReport *models.UserReportDB) *models.TotalIncentive {

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
		VOC: &models.IndividualIncentive{
			Payout:      userReport.VOCPayout,
			ExtraPayout: userReport.VOCExtraPayout,
		},
		TotalPay: userReport.TotalPayout,
	}
}
