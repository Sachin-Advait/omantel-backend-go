package impl

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) updateAdjustedAgentUserReport(ctx *gin.Context, user *models.UserBasicDetailsDB, chanel string,
	isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
				Msg("AAdjustedReport : nil info")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAdjustedAgentUserReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("AAdjustedReport : invalid info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAdjustedAgentUserReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).
		Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("AAdjustedReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).
			Msg("AAdjustedReport : invalid position")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details: "updateAdjustedAgentUserReport: incorrect agent position - userId:" +
				user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Msg("AAdjustedReport : invalid outlet info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.
				Region + " outlet:" + user.Outlet,
			Details: "updateAdjustedAgentUserReport: agent blank region or outlet - userId: " +
				user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAdjustedAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
			Msg("AAdjustedReport : fetch outlet info failed")
		return nil, nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Msg("AAdjustedReport : outlet co-ordinates invalid")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" +
				user.Region + " outlet:" + user.Outlet,
			Details: "updateAdjustedAgentUserReport: " + user.UserId + ": AGENT invalid outlet co-ordintaes - region:" +
				user.Region + " outlet:" + user.Outlet,
		}
	}

	totalEmp, err := impl.DbService.FetchTotalNonManagerialEmployees(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAdjustedAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AAdjustedReport : fetch empCount failed")
		return nil, nil, err
	}

	allUserSales, err := impl.DbService.FetchAdjustedActiveUserSales(user.UserId, chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAdjustedAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AAdjustedReport : fetch sales failed")
		return nil, nil, err
	}
	allUserSales.Region = user.Region
	allUserSales.Outlet = user.Outlet

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("agentCount", totalEmp).
		Int("saleGroupCount", len(allUserSales.SaleRecords)).Msg("AAdjustedReport : total agents and salesGroup")

	resp := &models.TotalIncentive{}
	report := &models.UserReportDB{
		UserId:          user.UserId,
		UserName:        user.Name,
		Position:        constant.EmployeePositionMap[user.Position],
		ManagerId:       user.ManagerId,
		ManagerPosition: user.ManagerPosition,
		Region:          user.Region,
		Outlet:          user.Outlet,
		Latitude:        outletResponse.Lat,
		Longitude:       outletResponse.Lang,
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}

	if len(allUserSales.SaleRecords) == 0 {
		ind := &models.IndividualIncentive{Error: "NO SALES FOUND FOR OUTLET - " + user.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		resp.VOC = ind
		return resp, report, nil
	}

	if totalEmp == 0 {
		ind := &models.IndividualIncentive{Error: "NO AGENT FOUND FOR OUTLET - " + user.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		resp.VOC = ind
		return resp, report, nil
	}

	totalPay := 0.0

	grpPlans := make(map[string]map[string]int)
	upgradedPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)
	totalActiveUpgrade := 0

	for grp, rcds := range allUserSales.SaleRecords {

		for _, v := range rcds {

			if v.Category != constant.UpgradeCategory {
				if grpPlans[grp] == nil {
					grpPlans[grp] = make(map[string]int)
				}
				grpPlans[grp][v.CurrentProduct] += v.Quantity
				totalActive[grp] += v.Quantity
			} else {
				if upgradedPlans[v.CurrentProduct] == nil {
					upgradedPlans[v.CurrentProduct] = make(map[string]int)
				}

				upgradedPlans[v.CurrentProduct][v.PreviousProduct] += v.Quantity
				totalActiveUpgrade += v.Quantity
			}
		}
	}

	for grp, qty := range totalActive {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("category", grp).
			Int("count", qty).Msg("AAdjustedReport : sale count")
	}
	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
		Int("count", totalActiveUpgrade).Msg("AAdjustedReport : sale count")

	for grp, plans := range grpPlans {
		targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport,
			grp, user.Region, user.Outlet)
		if internalErr != nil {
			internalErr = fmt.Errorf("updateAdjustedAgentUserReport: %v", internalErr.Error())
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
				Err(internalErr).Msg("AAdjustedReport : fetch target failed")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Message:    fmt.Sprintf("INVALID TARGET: reg: %v, out: %v, grp: %v", user.Region, user.Outlet, grp),
				Details:    internalErr.Error(),
			}
		}

		if targetPerOutlet <= 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).
				Str("category", grp).Msg("AAdjustedReport : zero target found")
			ind := &models.IndividualIncentive{Error: "NO TARGET FOUND FOR OUTLET - " + user.Outlet}
			resp.HBB = ind
			resp.Baqati = ind
			resp.Hayyak = ind
			resp.Upgrade = ind
			resp.VOC = ind
			return resp, report, nil
		}

		switch grp {

		case constant.HBBCategory:

			ind := &models.IndividualIncentive{}

			/*
					For each group like Baqati, Hayyak, HBB
					target := targetPerOutlet / totalEmployeeInThatOutlet
					achievedPercent := saleCountDone * 100 / target
					valueOfPoint := from achievedPercent

					for each plan sale do the following:
						1. quantityOfEachPlan * valueOfPoint * SIPRateOfThatPlan
						2. add all these values to get a totalPayout

				   if achievedPercent reaches a given value then extraPay is also given
			*/

			targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)
			// if targetPerEmp <= 0 {
			// 	ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE HBB TARGET-%v", user, targetPerEmp)
			// 	resp.HBB = ind
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).
			// 		Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HBB sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : HBB fetch vop failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : HBB fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : HBB fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : HBB fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).
					Msg("AAdjustedReport : HBB fetch threshold failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}

			report.HBBTargetCount = targetPerEmp
			report.HBBSalesCount = totalActive[grp]
			report.HBBAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			if ind.ExtraPayout > maxExtraPayAllowed {
				ind.ExtraPayout = maxExtraPayAllowed
			}

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).
				Float64("targetPerEmp", targetPerEmp).Int("achievedCount", totalActive[grp]).
				Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).
				Float64("threshold", threshold).Msg("AAdjustedReport : HBB sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)
			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).
					Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
					Float64("cumulativePay", ind.Payout).Msg("AAdjustedReport : HBB per plan info")
			}

			if ind.Error != "" {
				break
			}

			sort.Slice(saleCount, func(i, j int) bool {
				return saleCount[i].Quantity > saleCount[j].Quantity
			})
			hBBSaleCountPerPlanArray, _ := json.Marshal(saleCount)
			report.HBBSaleCountPerPlanArray = string(hBBSaleCountPerPlanArray)

			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}
			if ind.Payout > threshold {
				ind.Payout = threshold
			}
			ind.Payout = math.Round(ind.Payout*100) / 100

			totalPay += ind.Payout
			totalPay += ind.ExtraPayout

			resp.HBB = ind
			report.HBBPayout = ind.Payout
			report.HBBExtraPayout = ind.ExtraPayout

		case constant.HayyakCategory:

			ind := &models.IndividualIncentive{}

			targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)
			// if targetPerEmp <= 0 {
			// 	ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE HAYYAK TARGET-%v", user, targetPerEmp)
			// 	resp.Hayyak = ind
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).
			// 		Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HAYYAK sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : HAYYAK fetch vop failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : HAYYAK fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : HAYYAK fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : HAYYAK fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).
					Msg("AAdjustedReport : HAYYAK fetch threshold failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}

			report.HayyakTargetCount = targetPerEmp
			report.HayyakSalesCount = totalActive[grp]
			report.HayyakAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			if ind.ExtraPayout > maxExtraPayAllowed {
				ind.ExtraPayout = maxExtraPayAllowed
			}

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).
				Float64("targetPerEmp", targetPerEmp).Int("achievedCount", totalActive[grp]).
				Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).
				Float64("threshold", threshold).Msg("AAdjustedReport : HAYYAK sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)
			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).
					Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
					Float64("cumulativePay", ind.Payout).Msg("AAdjustedReport : HAYYAK per plan info")
			}

			if ind.Error != "" {
				break
			}

			sort.Slice(saleCount, func(i, j int) bool {
				return saleCount[i].Quantity > saleCount[j].Quantity
			})
			hayyakSaleCountPerPlanArray, _ := json.Marshal(saleCount)
			report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)

			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}
			if ind.Payout > threshold {
				ind.Payout = threshold
			}
			ind.Payout = math.Round(ind.Payout*100) / 100

			totalPay += ind.Payout
			totalPay += ind.ExtraPayout

			resp.Hayyak = ind
			report.HayyakPayout = ind.Payout
			report.HayyakExtraPayout = ind.ExtraPayout

		case constant.BaqatiCategory:

			ind := &models.IndividualIncentive{}

			targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)
			// if targetPerEmp <= 0 {
			// 	ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE BAQATI TARGET-%v", user, targetPerEmp)
			// 	resp.Baqati = ind
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).
			// 		Int("targetPerOutlet", targetPerOutlet).Msg("BAQATI sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : BAQATI fetch vop failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : BAQATI fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : BAQATI fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : BAQATI fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).
					Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).
					Msg("AReport : AAdjustedReport fetch threshold failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}

			report.BaqatiTargetCount = targetPerEmp
			report.BaqatiSalesCount = totalActive[grp]
			report.BaqatiAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			if ind.ExtraPayout > maxExtraPayAllowed {
				ind.ExtraPayout = maxExtraPayAllowed
			}

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).
				Float64("targetPerEmp", targetPerEmp).Int("achievedCount", totalActive[grp]).
				Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).
				Float64("threshold", threshold).Msg("AAdjustedReport : BAQATI sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)

			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).
					Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
					Float64("cumulativePay", ind.Payout).Msg("AAdjustedReport : BAQATI per plan info")
			}

			if ind.Error != "" {
				break
			}

			sort.Slice(saleCount, func(i, j int) bool {
				return saleCount[i].Quantity > saleCount[j].Quantity
			})
			baqatiSaleCountPerPlanArray, _ := json.Marshal(saleCount)
			report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)

			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}
			if ind.Payout > threshold {
				ind.Payout = threshold
			}
			ind.Payout = math.Round(ind.Payout*100) / 100

			totalPay += ind.Payout
			totalPay += ind.ExtraPayout

			resp.Baqati = ind
			report.BaqatiPayout = ind.Payout
			report.BaqatiExtraPayout = ind.ExtraPayout
		}
	}

	ind := &models.IndividualIncentive{}

	targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport,
		constant.UpgradeCategory, user.Region, user.Outlet)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Err(internalErr).Msg("AAdjustedReport : fetch target failed")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	if targetPerOutlet <= 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Msg("AAdjustedReport : zero target found")
		ind.Error = "NO TARGET FOUND FOR OUTLET - " + user.Outlet
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}

	targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
	achieved := float64(totalActiveUpgrade) * 100 / float64(targetPerEmp)
	// if targetPerEmp <= 0 {
	// 	ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE UPGRADE TARGET-%v", user, targetPerEmp)
	// 	resp.Upgrade = ind
	// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).
	// 		Int("targetPerOutlet", targetPerOutlet).Msg("UPGRADE sales error")
	// 	return resp, report, nil
	// }

	vop, internalErr := impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : UPGRADE fetch vop failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	extraPayout, internalErr := impl.getExtraPay(quarter, year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Float64("achieved", achieved).Err(internalErr).Msg("AAdjustedReport : UPGRADE fetch extraPayout failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : UPGRADE fetch maxPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Str("position", constant.PromoterCC).Err(internalErr).Msg("AAdjustedReport : UPGRADE fetch maxExtraPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).
			Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).
			Msg("AAdjustedReport : UPGRADE fetch threshold failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}

	report.UpgradeTargetCount = targetPerEmp
	report.UpgradeSalesCount = totalActiveUpgrade
	report.UpradeAchievedPercentage = achieved
	ind.ExtraPayout = math.Round(extraPayout*100) / 100
	if ind.ExtraPayout > maxExtraPayAllowed {
		ind.ExtraPayout = maxExtraPayAllowed
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).
		Float64("targetPerEmp", targetPerEmp).Int("achievedCount", totalActiveUpgrade).
		Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
		Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).
		Float64("threshold", threshold).Msg("AAdjustedReport : UPGRADE sales info")

	saleCount := make([]*models.SaleCountPerPlan, 0)
	for currPlan, prevPlansMap := range upgradedPlans {
		eachPlanSale := &models.SaleCountPerPlan{CurrentPlan: currPlan}
		for prevPln, qty := range prevPlansMap {
			eachPlanSale.Quantity += qty
			currPlanRate, _ := impl.getSIPRatePlan(quarter, year, currPlan)
			prevPlanRate, _ := impl.getSIPRatePlan(quarter, year, prevPln)
			rate := currPlanRate - prevPlanRate
			ind.Payout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("currPlanName", currPlan).
				Str("prevPlanName", prevPln).Float64("currPlanRate", currPlanRate).Float64("prevPlanRate", prevPlanRate).
				Float64("netPlanRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Float64("cumulativePay", ind.Payout).Msg("AAdjustedReport : UPGRADE per plan info")
		}

		saleCount = append(saleCount, eachPlanSale)
	}

	sort.Slice(saleCount, func(i, j int) bool {
		return saleCount[i].Quantity > saleCount[j].Quantity
	})
	upgradeSaleCountPerPlanArray, _ := json.Marshal(saleCount)
	report.UpgradeSaleCountPerPlanArray = string(upgradeSaleCountPerPlanArray)

	if ind.Payout > maxPayAllowed {
		ind.Payout = maxPayAllowed
	}
	if ind.Payout > threshold {
		ind.Payout = threshold
	}
	ind.Payout = math.Round(ind.Payout*100) / 100

	totalPay += ind.Payout
	totalPay += ind.ExtraPayout

	resp.Upgrade = ind
	report.UpgradePayout = ind.Payout
	report.UpgradeExtraPayout = ind.ExtraPayout

	resp.TotalPay = math.Round(totalPay*100) / 100
	report.TotalPayout = resp.TotalPay

	resp.VOC = &models.IndividualIncentive{}

	if quarter == 4 { //|| month == 12

		achievedPercent, err := impl.DbService.FetchVOCUserRecord([]string{user.UserId}, year)
		if err != nil && err.StatusCode != http.StatusBadRequest {
			return nil, nil, err
		}

		if len(achievedPercent) > 0 {
			report.VOCAchievedPercentage = achievedPercent[0].AchievedPercent
			report.VOCExtraPayout = report.VOCAchievedPercentage
			resp.VOC.ExtraPayout = report.VOCAchievedPercentage
		}
	}

	go func(user *models.UserReportDB) {
		if user.HBBTargetCount <= 0 || user.HayyakTargetCount <= 0 || user.BaqatiTargetCount <= 0 ||
			user.UpgradeTargetCount <= 0 {
			return
		}

		err = impl.CacheService.InsertUserReport(chanel, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).
				Err(errors.New(err.Details.(string))).Msg("AAdjustedReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("AAdjustedReport : cache update success")
		}

		// if err := impl.DbService.InsertAgentUserReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
		// 		Msg("AAdjustedReport : force report update")
		// }
	}(report)

	return resp, report, nil
}

func (impl *UserImpl) updateAdjustedOutletManagerReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	chanel string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {
		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
				Int("year", year).Msg("OMAdjustedReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateAdjustedOutletManagerReport: invalid omngr details",
			}
		}
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("OMAdjustedReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateAdjustedOutletManagerReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb {
		//!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).
			Msg("OMAdjustedReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details: "updateAdjustedOutletManagerReport: incorrect omngr position - userId:" +
				omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Msg("OMAdjustedReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" +
				omngr.Region + " outlet:" + omngr.Outlet,
			Details: "updateAdjustedOutletManagerReport: agent blank region or outlet - userId: " +
				omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
		Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("OMAdjustedReport : info")

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(omngr.Region, omngr.Outlet, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(err).Msg("OMAdjustedReport : fetch outlet info failed")
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMAdjustedReport : outlet co-ordinates invalid")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: omngr.UserId + ": OUTLET MANAGER INVALID OUTLET CO-ORDINATES - region:" +
				omngr.Region + " outlet:" + omngr.Outlet,
			Details: omngr.UserId + ": OUTLET MANAGER invalid outlet co-ordintaes - region:" +
				omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(errors.New(err.Details.(string))).Msg("OMAdjustedReport : fetch agents failed")
		return nil, err
	}

	report := &models.UserReportDB{
		UserId:          omngr.UserId,
		UserName:        omngr.Name,
		Position:        omngr.Position,
		ManagerId:       omngr.ManagerId,
		ManagerPosition: omngr.ManagerPosition,
		Region:          omngr.Region,
		Outlet:          omngr.Outlet,
		Latitude:        outletResponse.Lat,
		Longitude:       outletResponse.Lang,
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMReport : zero agents")
		return report, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).
			Msg("OMAdjustedReport : total agents")
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	for _, eachAgent := range userAgents {
		agentReport, err := impl.CacheService.GetUserReport(eachAgent.UserId, chanel, isQuarterReport, month, quarter, year)
		if !(err == nil && agentReport != nil && objecthelper.IsCurrentDay(agentReport.ModifiedAt)) {
			ind, resp, err := impl.updateAdjustedAgentUserReport(ctx, eachAgent, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("OMAdjustedReport : agent report failed")
				// return nil, err
				continue
			}
			rawIncentiveReport, _ := json.Marshal(ind)
			rawUserReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Str("incentive", string(rawIncentiveReport)).Str("record", string(rawUserReport)).
				Msg("OMAdjustedReport : agent log")
			agentReport = resp
		}

		report.HayyakSalesCount += agentReport.HayyakSalesCount
		report.HBBSalesCount += agentReport.HBBSalesCount
		report.BaqatiSalesCount += agentReport.BaqatiSalesCount
		report.UpgradeSalesCount += agentReport.UpgradeSalesCount

		saleCount := make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(agentReport.HBBSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hbbSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(agentReport.HayyakSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hayyakSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(agentReport.BaqatiSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				baqatiSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(agentReport.UpgradeSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				upgradeSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}
	}

	/*
		HBB Updates
	*/
	target, tarErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HBBCategory,
		omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(tarErr).Msg("OMAdjustedReport : HBB target not found")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "HBB TARGET NOT FOUND",
			Details:    fmt.Sprintf("hbb target errr - %v", tarErr),
		}
	}

	achieved := float64(report.HBBSalesCount) * 100 / float64(target)
	vop, _ := impl.getValueOfPoint(quarter, year, constant.HBBCategory, achieved)
	maxPayAllowed, _ := impl.getMaxPayout(quarter, year, constant.OutletManagerCC, constant.HBBCategory)
	threshold, _ := impl.getThresholdPayout(quarter, year, constant.OutletManagerCC, constant.HBBCategory, achieved)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).
		Int("achievedCount", report.HBBSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
		Msg("OMAdjustedReport : HBB sales info")

	for plan, qty := range hbbSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.HBBPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HBBPayout).
			Msg("OMAdjustedReport : HBB per plan info")
	}

	report.HBBTargetCount = float64(target)
	report.HBBAchievedPercentage = achieved
	if report.HBBPayout > maxPayAllowed {
		report.HBBPayout = maxPayAllowed
	}
	if report.HBBPayout > threshold {
		report.HBBPayout = threshold
	}
	//hbb sales count per plan
	saleArray := make([]*models.SaleCountPerPlan, 0)
	for key, val := range hbbSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	hbbSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HBBSaleCountPerPlanArray = string(hbbSaleCountPerPlanArray)

	/*
		HAYYAK Updates
	*/
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HayyakCategory,
		omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(tarErr).Msg("OMAdjustedReport : HAYYAK target not found")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "HAYYAK TARGET NOT FOUND",
			Details:    fmt.Sprintf("hayyak target errr - %v", tarErr),
		}
	}

	achieved = float64(report.HayyakSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.HayyakCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManagerCC, constant.HayyakCategory)
	threshold, _ = impl.getThresholdPayout(quarter, year, constant.OutletManagerCC, constant.HayyakCategory, achieved)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).
		Int("achievedCount", report.HayyakSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
		Msg("OMAdjustedReport : HAYYAK sales info")

	for plan, qty := range hayyakSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.HayyakPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HayyakPayout).
			Msg("OMAdjustedReport : HAYYAK per plan info")
	}

	report.HayyakTargetCount = float64(target)
	report.HayyakAchievedPercentage = achieved
	if report.HayyakPayout > maxPayAllowed {
		report.HayyakPayout = maxPayAllowed
	}
	if report.HayyakPayout > threshold {
		report.HayyakPayout = threshold
	}
	//hayyak sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range hayyakSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	hayyakSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)

	/*
		BAQATI Updates
	*/
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.BaqatiCategory,
		omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(tarErr).Msg("OMAdjustedReport : BAQATI target not found")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "BAQATI TARGET NOT FOUND",
			Details:    fmt.Sprintf("baqati target errr - %v", tarErr),
		}
	}

	achieved = float64(report.BaqatiSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.BaqatiCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManagerCC, constant.BaqatiCategory)
	threshold, _ = impl.getThresholdPayout(quarter, year, constant.OutletManagerCC, constant.BaqatiCategory, achieved)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).
		Int("achievedCount", report.BaqatiSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
		Msg("OMAdjustedReport : BAQATI sales info")

	for plan, qty := range baqatiSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.BaqatiPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.BaqatiPayout).
			Msg("OMAdjustedReport : BAQATI per plan info")
	}

	report.BaqatiTargetCount = float64(target)
	report.BaqatiAchievedPercentage = achieved
	if report.BaqatiPayout > maxPayAllowed {
		report.BaqatiPayout = maxPayAllowed
	}
	if report.BaqatiPayout > threshold {
		report.BaqatiPayout = threshold
	}
	//baqati sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range baqatiSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	baqatiSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)

	/*
		UPGRADE Updates
	*/
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory,
		omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).
			Err(tarErr).Msg("OMAdjustedReport : UPGRADE target not found")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "UPGRADE TARGET NOT FOUND",
			Details:    fmt.Sprintf("upgrade target errr - %v", tarErr),
		}
	}

	achieved = float64(report.UpgradeSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManagerCC, constant.UpgradeCategory)
	threshold, _ = impl.getThresholdPayout(quarter, year, constant.OutletManagerCC, constant.UpgradeCategory, achieved)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).
		Int("achievedCount", report.UpgradeSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
		Msg("OMAdjustedReport : UPGRADE sales info")

	for plan, qty := range upgradeSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.UpgradePayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.UpgradePayout).
			Msg("OMAdjustedReport : UPGRADE per plan info")
	}

	report.UpgradeTargetCount = float64(target)
	report.UpradeAchievedPercentage = achieved
	if report.UpgradePayout > maxPayAllowed {
		report.UpgradePayout = maxPayAllowed
	}
	if report.UpgradePayout > threshold {
		report.UpgradePayout = threshold
	}
	//upgrade sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range upgradeSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	upgradeSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.UpgradeSaleCountPerPlanArray = string(upgradeSaleCountPerPlanArray)

	report.TotalPayout = report.HBBPayout + report.HayyakPayout + report.BaqatiPayout + report.UpgradePayout

	if quarter == 4 { //|| month == 12
		achievedPercent, err := impl.DbService.FetchVOCUserRecord([]string{omngr.UserId}, year)
		if err != nil && err.StatusCode != http.StatusBadRequest {
			return nil, err
		}
		if len(achievedPercent) > 0 {
			report.VOCAchievedPercentage = achievedPercent[0].AchievedPercent
			report.VOCExtraPayout = report.VOCAchievedPercentage
		}
	}

	go func(user *models.UserReportDB) {
		err = impl.CacheService.InsertUserReport(chanel, isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).
				Err(errors.New(err.Details.(string))).Msg("OMAdjustedReport : cache update failed")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("OMAdjustedReport : cache update success")
		}
		// if err := impl.DbService.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
		// 		Msg("OMAdjustedReport : force report update")
		// }
	}(report)

	return report, nil
}

func (impl *UserImpl) updateAdjustedManagerReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB,
	chanel string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).
			Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("MngrAdjustedReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateAdjustedManagerReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).
			Str("position", mngr.Position).Msg("MngrAdjustedReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details: "updateAdjustedManagerReport: incorrect mngr position - userId:" +
				mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).
		Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrAdjustedReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateAdjustedOutletManagerReport(ctx, mngr, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).
			Str("position", mngr.Position).Msg("MngrAdjustedReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	report := &models.UserReportDB{
		UserId:          mngr.UserId,
		UserName:        mngr.Name,
		Position:        mngr.Position,
		ManagerId:       mngr.ManagerId,
		ManagerPosition: mngr.ManagerPosition,
		Region:          mngr.Region,
		Outlet:          mngr.Outlet,
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(mngr.Region, mngr.Outlet, quarter, year)
	if err == nil {
		report.Latitude = outletResponse.Lat
		report.Longitude = outletResponse.Lang
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Err(errors.New(err.Details.(string))).Msg("MngrAdjustedReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrAdjustedReport : zero om")
		return report, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).
			Int("omCount", len(allDistinctOM)).Msg("MngrAdjustedReport : total om")
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	inactiveOM := make([]string, 0)
	emptyTarget := make(map[string][]string, 0)

	for _, om := range allDistinctOM {

		omReport, err := impl.CacheService.GetUserReport(om.UserId, chanel, isQuarterReport, month, quarter, year)

		if !(err == nil && omReport != nil && objecthelper.IsCurrentDay(omReport.ModifiedAt)) {
			resp, err := impl.updateAdjustedOutletManagerReport(ctx, om, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("om", om.UserId).
					Err(err).Msg("MngrAdjustedReport : om report failed")
				continue
				//return nil, err
			}
			rawUserReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", om.UserId).
				Str("record", string(rawUserReport)).Msg("MngrAdjustedReport : agent log")
			omReport = resp
		}

		report.HBBSalesCount += omReport.HBBSalesCount
		report.HayyakSalesCount += omReport.HayyakSalesCount
		report.BaqatiSalesCount += omReport.BaqatiSalesCount
		report.UpgradeSalesCount += omReport.UpgradeSalesCount

		if omReport.HBBTargetCount == 0 {
			if emptyTarget[constant.HBBCategory] == nil {
				emptyTarget[constant.HBBCategory] = []string{omReport.UserId}
			} else {
				emptyTarget[constant.HBBCategory] = append(emptyTarget[constant.HBBCategory], omReport.UserId)
			}
		} else if omReport.HayyakTargetCount == 0 {
			if emptyTarget[constant.HayyakCategory] == nil {
				emptyTarget[constant.HayyakCategory] = []string{omReport.UserId}
			} else {
				emptyTarget[constant.HayyakCategory] = append(emptyTarget[constant.HayyakCategory], omReport.UserId)
			}
		} else if omReport.BaqatiTargetCount == 0 {
			if emptyTarget[constant.BaqatiCategory] == nil {
				emptyTarget[constant.BaqatiCategory] = []string{omReport.UserId}
			} else {
				emptyTarget[constant.BaqatiCategory] = append(emptyTarget[constant.BaqatiCategory], omReport.UserId)
			}
		} else if omReport.UpgradeTargetCount == 0 {
			if emptyTarget[constant.UpgradeCategory] == nil {
				emptyTarget[constant.UpgradeCategory] = []string{omReport.UserId}
			} else {
				emptyTarget[constant.UpgradeCategory] = append(emptyTarget[constant.UpgradeCategory], omReport.UserId)
			}
		}

		if om.IsActive {
			report.HayyakTargetCount += omReport.HayyakTargetCount
			report.HBBTargetCount += omReport.HBBTargetCount
			report.BaqatiTargetCount += omReport.BaqatiTargetCount
			report.UpgradeTargetCount += omReport.UpgradeTargetCount
		} else {
			inactiveOM = append(inactiveOM, om.UserId)
		}

		saleCount := make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(omReport.HBBSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hbbSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(omReport.HayyakSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hayyakSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(omReport.BaqatiSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				baqatiSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}

		saleCount = make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(omReport.UpgradeSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				upgradeSaleMap[sale.CurrentPlan] += sale.Quantity
			}
		}
	}

	logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Interface("inactive", inactiveOM).
		Interface("blank", emptyTarget).Msg("MngrAdjustedReport : om target defaulters")

	/*
		HBB Updates
	*/
	if report.HBBTargetCount == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrAdjustedReport : HBB target not found")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "HBB TARGET NOT FOUND",
			Details:    "mngr hbb target not found",
		}
	}
	achieved := float64(report.HBBSalesCount) * 100 / float64(report.HBBTargetCount)
	report.HBBAchievedPercentage = achieved
	//hbb sales count per plan
	saleArray := make([]*models.SaleCountPerPlan, 0)
	for key, val := range hbbSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	hbbSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HBBSaleCountPerPlanArray = string(hbbSaleCountPerPlanArray)

	/*
		HAYYAK Updates
	*/
	achieved = float64(report.HayyakSalesCount) * 100 / float64(report.HayyakTargetCount)
	report.HayyakAchievedPercentage = achieved
	//hayyak sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range hayyakSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	hayyakSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)

	/*
		BAQATI Updates
	*/
	achieved = float64(report.BaqatiSalesCount) * 100 / float64(report.BaqatiTargetCount)
	report.BaqatiAchievedPercentage = achieved
	//baqati sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range baqatiSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	baqatiSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)

	/*
		UPGRADE Updates
	*/
	achieved = float64(report.UpgradeSalesCount) * 100 / float64(report.UpgradeTargetCount)
	report.UpradeAchievedPercentage = achieved
	//upgrade sales count per plan
	saleArray = make([]*models.SaleCountPerPlan, 0)
	for key, val := range upgradeSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}
	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})
	upgradeSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.UpgradeSaleCountPerPlanArray = string(upgradeSaleCountPerPlanArray)

	if mngr.Position == constant.ManagerAbb || mngr.Position == constant.SeniorManagerAbb {
		//HBB
		var maxPayAllowed float64
		if mngr.Position == constant.ManagerAbb {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.ManagerCC, constant.HBBCategory)
		} else {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.SeniorManagerAbb, constant.HBBCategory)
		}

		vop, _ := impl.getValueOfPoint(quarter, year, constant.HBBCategory, report.HBBAchievedPercentage)
		//maxPayAllowed, _ := impl.getMaxPayout(quarter, year, mngr.Position, constant.HBBCategory)
		threshold, _ := impl.getThresholdPayout(quarter, year, constant.ManagerCC,
			constant.HBBCategory, report.HBBAchievedPercentage)

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HBBTargetCount).
			Int("achievedCount", report.HBBSalesCount).Float64("achievedPercent", report.HBBAchievedPercentage).
			Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
			Msg("MngrAdjustedReport : HBB sales info")

		for plan, qty := range hbbSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HBBPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).
				Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Float64("cumulativePay", report.HBBPayout).Msg("MngrAdjustedReport : HBB per plan info")
		}

		if report.HBBPayout > maxPayAllowed {
			report.HBBPayout = maxPayAllowed
		}
		if mngr.Position == constant.ManagerAbb {
			if report.HBBPayout > threshold {
				report.HBBPayout = threshold
			}
		}

		//HAYYAK
		if mngr.Position == constant.ManagerAbb {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.ManagerCC, constant.HayyakCategory)
		} else {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.SeniorManagerAbb, constant.HayyakCategory)
		}

		vop, _ = impl.getValueOfPoint(quarter, year, constant.HayyakCategory, report.HayyakAchievedPercentage)
		//maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.HayyakCategory)
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager,
			constant.HayyakCategory, report.HayyakAchievedPercentage)

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HayyakTargetCount).
			Int("achievedCount", report.HayyakSalesCount).Float64("achievedPercent", report.HayyakAchievedPercentage).
			Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
			Msg("MngrAdjustedReport : HAYYAK sales info")

		for plan, qty := range hayyakSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HayyakPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).
				Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Float64("cumulativePay", report.HayyakPayout).Msg("MngrAdjustedReport : HAYYAK per plan info")
		}

		if report.HayyakPayout > maxPayAllowed {
			report.HayyakPayout = maxPayAllowed
		}
		if mngr.Position == constant.ManagerAbb {
			if report.HayyakPayout > threshold {
				report.HayyakPayout = threshold
			}
		}

		//BAQATI
		if mngr.Position == constant.ManagerAbb {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.ManagerCC, constant.BaqatiCategory)
		} else {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.SeniorManagerAbb, constant.BaqatiCategory)
		}

		vop, _ = impl.getValueOfPoint(quarter, year, constant.BaqatiCategory, report.BaqatiAchievedPercentage)
		//maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.BaqatiCategory)
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager,
			constant.BaqatiCategory, report.BaqatiAchievedPercentage)

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.BaqatiTargetCount).
			Int("achievedCount", report.BaqatiSalesCount).Float64("achievedPercent", report.BaqatiAchievedPercentage).
			Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
			Msg("MngrAdjustedReport : BAQATI sales info")

		for plan, qty := range baqatiSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.BaqatiPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).
				Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Float64("cumulativePay", report.BaqatiPayout).Msg("MngrAdjustedReport : BAQATI per plan info")
		}

		if report.BaqatiPayout > maxPayAllowed {
			report.BaqatiPayout = maxPayAllowed
		}
		if mngr.Position == constant.ManagerAbb {
			if report.BaqatiPayout > threshold {
				report.BaqatiPayout = threshold
			}
		}

		//UPGRADE
		if mngr.Position == constant.ManagerAbb {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.ManagerCC, constant.UpgradeCategory)
		} else {
			maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.SeniorManagerAbb, constant.UpgradeCategory)
		}

		vop, _ = impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, report.UpradeAchievedPercentage)
		//maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.UpgradeCategory)
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager,
			constant.UpgradeCategory, report.UpradeAchievedPercentage)

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.UpgradeTargetCount).
			Int("achievedCount", report.UpgradeSalesCount).Float64("achievedPercent", report.UpradeAchievedPercentage).
			Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).
			Msg("MngrAdjustedReport : UPGRADE sales info")

		for plan, qty := range upgradeSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.UpgradePayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).
				Float64("planRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Float64("cumulativePay", report.UpgradePayout).Msg("MngrAdjustedReport : UPGRADE per plan info")
		}

		if report.UpgradePayout > maxPayAllowed {
			report.UpgradePayout = maxPayAllowed
		}
		if mngr.Position == constant.ManagerAbb {
			if report.UpgradePayout > threshold {
				report.UpgradePayout = threshold
			}
		}

		report.TotalPayout = report.HBBPayout + report.HayyakPayout + report.BaqatiPayout + report.UpgradePayout
	}

	if quarter == 4 { //|| month == 12

		achievedPercent, err := impl.DbService.FetchVOCUserRecord([]string{mngr.UserId}, year)
		if err != nil && err.StatusCode != http.StatusBadRequest {
			return nil, err
		}

		if len(achievedPercent) > 0 {
			report.VOCAchievedPercentage = achievedPercent[0].AchievedPercent
			report.VOCExtraPayout = report.VOCAchievedPercentage
		}
	}

	go func(user *models.UserReportDB) {

		err = impl.CacheService.InsertUserReport(chanel, isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).
				Err(errors.New(err.Details.(string))).Msg("MngrAdjustedReport : cache update failed")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("MngrAdjustedReport : cache update success")
		}

		// if err := impl.DbService.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).
		// 		Err(errors.New(err.Details.(string))).Msg("MngrAdjustedReport : force report update")
		// }
	}(report)

	return report, nil
}
