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
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

// GetUserPayout - cached response
func (impl *UserImpl) GetUserPayout(ctx *gin.Context, userID string, chanel string, isQuarterReport bool, quarter,
	month, year int) (*models.TotalIncentive, *errors.Error) {

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

	// user, err := impl.DbService.FetchUserBasicDetails(userID, quarter, year)
	// if err != nil {
	// 	return nil, err
	// }

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

	if !isQuarterReport || quarter != 4 { //monthly or not last quarter
		return impl.CoalesceAjustedUserPayout(ctx, user, chanel, isQuarterReport, quarter, month, year)
	}

	q1Report, err := impl.CoalesceAjustedUserPayout(ctx, user, chanel, isQuarterReport, 1, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-1, err-%v", err.Details)
		return nil, err
	}
	q2Report, err := impl.CoalesceAjustedUserPayout(ctx, user, chanel, isQuarterReport, 2, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-2, err-%v", err.Details)
		return nil, err
	}
	q3Report, err := impl.CoalesceAjustedUserPayout(ctx, user, chanel, isQuarterReport, 3, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-3, err-%v", err.Details)
		return nil, err
	}
	q4Report, err := impl.CoalesceAjustedUserPayout(ctx, user, chanel, isQuarterReport, 4, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("coalesces-4, err-%v", err.Details)
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

// CoalesceAjustedUserPayout -
func (impl *UserImpl) CoalesceAjustedUserPayout(ctx *gin.Context, user *models.UserBasicDetailsDB, chanel string,
	isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

	// adjustedSales := &models.AdjustedReportDB{}
	// prevQuarter, prevYear := objecthelper.GetPreviousQuarterAndMonth(quarter, year)

	/* --------- fetch current payment ---------*/

	currentReport, err := impl.CacheService.GetUserReport(user.UserId, chanel, isQuarterReport, month, quarter, year)
	if err != nil || currentReport == nil || !objecthelper.IsCurrentDay(currentReport.ModifiedAt) {
		// calculation is not fresh
		if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
			//for promoter level
			_, currentReport, err = impl.updateAgentUserReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				//return error for report not found
				return nil, err
			}
		} else {
			//for manager level
			currentReport, err = impl.updateManagerReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				//return error for report not found
				return nil, err
			}
		}
	}

	// for monthly type ignore adjusted payment
	if !isQuarterReport {
		return impl.userFinalPayout(ctx, currentReport), nil
	}

	/* --------- fetch previous payment ---------*/

	prevQuarter, prevYear := objecthelper.GetPreviousQuarterAndMonth(quarter, year)

	prevReport, err := impl.CacheService.GetUserReport(user.UserId, chanel, isQuarterReport, month, prevQuarter, prevYear)
	if err != nil || currentReport == nil || !objecthelper.IsCurrentDay(prevReport.ModifiedAt) {
		// prev quarter calculation is not fresh
		if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
			//for promoter level prev quarter
			_, prevReport, err = impl.updateAgentUserReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
		} else {
			//for manager level prev quarter
			prevReport, err = impl.updateManagerReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
		}
	}
	// return currrent report without adjustment
	if err != nil || prevReport == nil {
		return impl.userFinalPayout(ctx, currentReport), nil
	}

	/* --------- fetch previous adjusted payment ---------*/

	prevAdjustedReport, err := impl.CacheService.GetUserAdjustedReport(user.UserId, chanel, isQuarterReport, month,
		prevQuarter, prevYear)
	if err != nil || prevAdjustedReport == nil || !objecthelper.IsCurrentDay(prevAdjustedReport.ModifiedAt) {
		if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
			_, prevAdjustedReport, err = impl.updateAdjustedAgentUserReport(ctx, user, chanel, isQuarterReport, prevQuarter,
				month, prevYear)
			if err != nil {
				return nil, err
			}
		} else {
			prevAdjustedReport, err = impl.updateAdjustedManagerReport(ctx, user, chanel, isQuarterReport, prevQuarter,
				month, prevYear)
			if err != nil {
				return nil, err
			}
		}
	}
	// return currrent report without adjustment
	if err != nil || prevAdjustedReport == nil {
		return impl.userFinalPayout(ctx, currentReport), nil
	}

	/*
		From current quarter totalPay subtract the difference between prevQuarterTotalPay and prevAdjustedQuarterTotalPay
		currentTotalPay = (currentTotalPay) - (prevQuarterTotalPay - prevAdjustedQuarterTotalPay)
	*/
	deductAmount := prevReport.TotalPayout - prevAdjustedReport.TotalPayout
	if deductAmount < 0 {
		//TODO: log error as prev adjusted pay cannot be greater than prev pay
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(fmt.Errorf("invalid deduction - prev:%v, adjustedPred:%v",
			prevReport.TotalPayout, prevAdjustedReport.TotalPayout)).Msg("CoalesceAjustedUserPayout : fetch empCount failed")
	} else {
		currentReport.TotalPayout = currentReport.TotalPayout - deductAmount
	}

	return impl.userFinalPayout(ctx, currentReport), nil
}

func (impl *UserImpl) userFinalPayout(_ *gin.Context, userReport *models.UserReportDB) *models.TotalIncentive {
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
		TotalPayBeforeDeduction: userReport.TotalPayout,
		TotalPay:                userReport.TotalPayout,
	}
}

// func (impl *UserImpl) calculateSaleAdjustment(_ *gin.Context, userReport *models.UserReportDB,
// 	adjustedReport *models.UserReportDB) (*models.AdjustedReportDB, *errors.Error) {

// 	// r1, _ := json.Marshal(userReport)
// 	// fmt.Println("old", string(r1))
// 	// r2, _ := json.Marshal(adjustedReport)
// 	// fmt.Println("new", string(r2))

// 	if userReport == nil || adjustedReport == nil {
// 		return nil, &errors.Error{
// 			StatusCode: http.StatusInternalServerError,
// 			Code:       errors.InternalServerErrorCode,
// 			Message:    errors.InternalServerError,
// 			Details:    "nil arguments",
// 		}
// 	}

// 	finalReport := &models.AdjustedReportDB{
// 		HayyakSalesCount:           userReport.HayyakSalesCount,
// 		HayyakUpdatedSalesCount:    adjustedReport.HayyakSalesCount,
// 		HayyakTerminatedSalesCount: 0,
// 		HayyakPayment:              userReport.HayyakPayout,
// 		HayyakUpdatedPayment:       adjustedReport.HayyakPayout,
// 		HayyakExtraPayment:         userReport.HayyakExtraPayout,
// 		HayyakUpdatedExtraPayment:  adjustedReport.HayyakExtraPayout,
// 		HayyakTotalPayment:         userReport.TotalPayout,
// 		HayyakUpdatedTotalPayment:  adjustedReport.TotalPayout,
// 		HayyakDeductAmount:         0,

// 		HBBSalesCount:           userReport.HBBSalesCount,
// 		HBBUpdatedSalesCount:    adjustedReport.HBBSalesCount,
// 		HBBTerminatedSalesCount: 0,
// 		HBBPayment:              userReport.HBBPayout,
// 		HBBUpdatedPayment:       adjustedReport.HBBPayout,
// 		HBBExtraPayment:         userReport.HBBExtraPayout,
// 		HBBUpdatedExtraPayment:  adjustedReport.HBBExtraPayout,
// 		HBBTotalPayment:         userReport.TotalPayout,
// 		HBBUpdatedTotalPayment:  adjustedReport.TotalPayout,
// 		HBBDeductAmount:         0,

// 		BaqatiSalesCount:           userReport.BaqatiSalesCount,
// 		BaqatiUpdatedSalesCount:    adjustedReport.BaqatiSalesCount,
// 		BaqatiTerminatedSalesCount: 0,
// 		BaqatiPayment:              userReport.BaqatiPayout,
// 		BaqatiUpdatedPayment:       adjustedReport.BaqatiPayout,
// 		BaqatiExtraPayment:         userReport.BaqatiExtraPayout,
// 		BaqatiUpdatedExtraPayment:  adjustedReport.BaqatiExtraPayout,
// 		BaqatiTotalPayment:         userReport.TotalPayout,
// 		BaqatiUpdatedTotalPayment:  adjustedReport.TotalPayout,
// 		BaqatiDeductAmount:         0,

// 		UpgradeSalesCount:           userReport.UpgradeSalesCount,
// 		UpgradeUpdatedSalesCount:    adjustedReport.UpgradeSalesCount,
// 		UpgradeTerminatedSalesCount: 0,
// 		UpgradePayment:              userReport.UpgradePayout,
// 		UpgradeUpdatedPayment:       adjustedReport.UpgradePayout,
// 		UpgradeExtraPayment:         userReport.UpgradeExtraPayout,
// 		UpgradeUpdatedExtraPayment:  adjustedReport.UpgradeExtraPayout,
// 		UpgradeTotalPayment:         userReport.TotalPayout,
// 		UpgradeUpdatedTotalPayment:  adjustedReport.TotalPayout,
// 		UpgradeDeductAmount:         0,
// 	}

// 	if finalReport.HayyakUpdatedSalesCount != 0 {
// 		finalReport.HayyakTerminatedSalesCount = finalReport.HayyakSalesCount - finalReport.HayyakUpdatedSalesCount
// 	}
// 	if finalReport.HayyakUpdatedTotalPayment != 0 {
// 		finalReport.HayyakDeductAmount = finalReport.HayyakTotalPayment - finalReport.HayyakUpdatedTotalPayment
// 	}

// 	if finalReport.HBBUpdatedSalesCount != 0 {
// 		finalReport.HBBTerminatedSalesCount = finalReport.HBBSalesCount - finalReport.HBBUpdatedSalesCount
// 	}
// 	if finalReport.HBBUpdatedTotalPayment != 0 {
// 		finalReport.HBBDeductAmount = finalReport.HBBTotalPayment - finalReport.HBBUpdatedTotalPayment
// 	}

// 	if finalReport.BaqatiUpdatedSalesCount != 0 {
// 		finalReport.BaqatiTerminatedSalesCount = finalReport.BaqatiSalesCount - finalReport.BaqatiUpdatedSalesCount
// 	}
// 	if finalReport.BaqatiUpdatedTotalPayment != 0 {
// 		finalReport.BaqatiDeductAmount = finalReport.BaqatiTotalPayment - finalReport.BaqatiUpdatedTotalPayment
// 	}

// 	if finalReport.UpgradeUpdatedSalesCount != 0 {
// 		finalReport.UpgradeTerminatedSalesCount = finalReport.UpgradeSalesCount - finalReport.UpgradeUpdatedSalesCount
// 	}
// 	if finalReport.UpgradeUpdatedTotalPayment != 0 {
// 		finalReport.UpgradeDeductAmount = finalReport.UpgradeTotalPayment - finalReport.UpgradeUpdatedTotalPayment
// 	}

// 	finalReport.TotalDeductedAmount = finalReport.HayyakDeductAmount + finalReport.HBBDeductAmount +
// 		finalReport.BaqatiDeductAmount + finalReport.UpgradeDeductAmount

// 	return finalReport, nil
// }

// func (impl *UserImpl) CoalesceUserPayout(ctx *gin.Context, userID string, chanel string, isQuarterReport bool,
// 	quarter, month, year int) (*models.TotalIncentive, *errors.Error) {

// 	user, err := impl.DbService.FetchUserDetails(userID, quarter, year)
// 	if err != nil {
// 		return nil, err
// 	}
// 	position, ok := constant.EmployeePositionMap[user.Position]
// 	if !ok {
// 		return nil, &errors.Error{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
// 			Details:    "incorrect user position: userId-" + user.UserId + " position-" + user.Position,
// 		}
// 	}
// 	user.Position = position

// 	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb || user.Position == constant.OutletManagerAbb {
// 		chanel = ""
// 	}

// 	report, err := impl.CacheService.GetUserReport(userID, chanel, isQuarterReport, month, quarter, year)
// 	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
// 		return impl.userPayout(ctx, report, nil), nil
// 	}

// 	if !isQuarterReport {
// 		quarter = getQuarterFromMonth(month)
// 	}

// 	/*
// 		For Agents
// 	*/

// 	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {

// 		resp, _, err := impl.updateAgentUserReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// val, _ := json.Marshal(resp)
// 		// logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Str("report", string(val)).Msg("GetUserPayout : updated")

// 		return resp, nil
// 	}

// 	/*
// 		For All Managers
// 	*/

// 	mngrReport, err := impl.updateManagerReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// val, _ := json.Marshal(mngrReport)
// 	// logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Str("report", string(val)).Msg("GetUserPayout : updated")

// 	return impl.userPayout(ctx, mngrReport, nil), nil
// }

// func (impl *UserImpl) userPayout(_ *gin.Context, userReport *models.UserReportDB,
// 	adjustedReport *models.UserReportDB) *models.TotalIncentive {

// 	finalReport := &models.TotalIncentive{
// 		HBB: &models.IndividualIncentive{
// 			Payout:      userReport.HBBPayout,
// 			ExtraPayout: userReport.HBBExtraPayout,
// 		},
// 		Hayyak: &models.IndividualIncentive{
// 			Payout:      userReport.HayyakPayout,
// 			ExtraPayout: userReport.HayyakExtraPayout,
// 		},
// 		Baqati: &models.IndividualIncentive{
// 			Payout:      userReport.BaqatiPayout,
// 			ExtraPayout: userReport.BaqatiExtraPayout,
// 		},
// 		Upgrade: &models.IndividualIncentive{
// 			Payout:      userReport.UpgradePayout,
// 			ExtraPayout: userReport.UpgradeExtraPayout,
// 		},
// 		VOC: &models.IndividualIncentive{
// 			Payout:      userReport.VOCPayout,
// 			ExtraPayout: userReport.VOCExtraPayout,
// 		},
// 		TotalPay: userReport.TotalPayout,
// 	}

// 	if adjustedReport == nil {
// 		return finalReport
// 	}

// 	if adjustedReport.HBBSalesCount != 0 {
// 		finalReport.HBB.Terminated = userReport.HBBSalesCount - adjustedReport.HBBSalesCount
// 	}
// 	if adjustedReport.BaqatiSalesCount != 0 {
// 		finalReport.Baqati.Terminated = userReport.BaqatiSalesCount - adjustedReport.BaqatiSalesCount
// 	}
// 	if adjustedReport.HayyakSalesCount != 0 {
// 		finalReport.Hayyak.Terminated = userReport.HayyakSalesCount - adjustedReport.HayyakSalesCount
// 	}
// 	if adjustedReport.UpgradeSalesCount != 0 {
// 		finalReport.Upgrade.Terminated = userReport.UpgradeSalesCount - adjustedReport.UpgradeSalesCount
// 	}

// 	hbbAdjustedAmount := adjustedReport.HBBPayout + adjustedReport.HBBExtraPayout
// 	baqatiAdjustedAmount := adjustedReport.BaqatiPayout + adjustedReport.BaqatiExtraPayout
// 	hayyakAdjustedAmount := adjustedReport.HayyakPayout + adjustedReport.HayyakExtraPayout
// 	upgradeAdjustedAmount := adjustedReport.UpgradePayout + adjustedReport.UpgradeExtraPayout

// 	if hbbAdjustedAmount != 0 {
// 		finalReport.HBB.DeductAmount = (userReport.HBBPayout + userReport.HBBExtraPayout) - hbbAdjustedAmount
// 	}
// 	if baqatiAdjustedAmount != 0 {
// 		finalReport.Baqati.DeductAmount = (userReport.BaqatiPayout + userReport.BaqatiExtraPayout) - baqatiAdjustedAmount
// 	}
// 	if hayyakAdjustedAmount != 0 {
// 		finalReport.Hayyak.DeductAmount = (userReport.HayyakPayout + userReport.HayyakExtraPayout) - hayyakAdjustedAmount
// 	}
// 	if upgradeAdjustedAmount != 0 {
// 		finalReport.Upgrade.DeductAmount = (userReport.UpgradePayout + userReport.UpgradeExtraPayout) - upgradeAdjustedAmount
// 	}

// 	netDeductedAmount := finalReport.HBB.DeductAmount + finalReport.Baqati.DeductAmount + finalReport.Hayyak.DeductAmount +
// 		finalReport.Upgrade.DeductAmount

// 	finalReport.TotalPay = userReport.TotalPayout - netDeductedAmount

// 	return finalReport
// }
