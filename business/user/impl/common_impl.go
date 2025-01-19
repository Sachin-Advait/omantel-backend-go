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

func (impl *UserImpl) getValueOfPoint(quarter, year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if models.ValOfPoint.VopMap[year] == nil {
		return 0, errors.New("VOP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[year][quarter] == nil {
		return 0, errors.New("VOP QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ValOfPoint.VopMap[year][quarter][grp] == nil {
		return 0, errors.New("VOP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent < 75:
		return models.ValOfPoint.VopMap[year][quarter][grp][75].Point, nil
	case achievedPercent < 96:
		return models.ValOfPoint.VopMap[year][quarter][grp][95].Point, nil
	case achievedPercent <= 111:
		return models.ValOfPoint.VopMap[year][quarter][grp][110].Point, nil
	case achievedPercent < 140:
		return models.ValOfPoint.VopMap[year][quarter][grp][139].Point, nil
	case achievedPercent >= 140:
		return models.ValOfPoint.VopMap[year][quarter][grp][140].Point, nil
	}

	return 0, errors.New("VOP ACHIEVEMENT PERCENT INVALID")
}

func (impl *UserImpl) getExtraPay(quarter, year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if models.ValOfPoint.VopMap[year] == nil {
		return 0, errors.New("VOEP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[year][quarter] == nil {
		return 0, errors.New("VOEP QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ValOfPoint.VopMap[year][quarter][grp] == nil {
		return 0, errors.New("VOEP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return models.ValOfPoint.VopMap[year][quarter][grp][159].ExtraPay, nil
	case achievedPercent >= 160 && achievedPercent < 180:
		return models.ValOfPoint.VopMap[year][quarter][grp][179].ExtraPay, nil
	case achievedPercent >= 180 && achievedPercent < 190:
		return models.ValOfPoint.VopMap[year][quarter][grp][189].ExtraPay, nil
	case achievedPercent >= 190 && achievedPercent < 200:
		return models.ValOfPoint.VopMap[year][quarter][grp][199].ExtraPay, nil
	case achievedPercent >= 200:
		return models.ValOfPoint.VopMap[year][quarter][grp][200].ExtraPay, nil
	}

	return 0, nil
}

func (impl *UserImpl) getSIPRatePlan(quarter, year int, planName string) (float64, error) {

	planID := ""

	models.PlanRates.Mutex.RLock()
	defer models.PlanRates.Mutex.RUnlock()

	if models.PlanRates.PlanRateMap[year] == nil {
		return 0, errors.New("PLAN-NAME YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.PlanRates.PlanRateMap[year][quarter] == nil {
		return 0, errors.New("PLAN-NAME QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	for id, v := range models.PlanRates.PlanRateMap[year][quarter] {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != "" {
		return models.PlanRates.PlanRateMap[year][quarter][planID].SIPRatePlan, nil
	}

	return 0, errors.New("PLAN-NAME INVALID - " + planName)
}

func (impl *UserImpl) getSaleTargetPerOutlet(year, month, quarter int, isQuarterReport bool, grp, region, outlet string) (int, error) {

	var (
		start int
		end   int
	)

	if isQuarterReport {
		switch quarter {
		case 1:
			start = 1
			end = 3
		case 2:
			start = 4
			end = 6
		case 3:
			start = 7
			end = 9
		case 4:
			start = 10
			end = 12
		default:
			return 0, errors.New("TARGET QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
		}
	}

	models.SalesTrgt.Mutex.RLock()
	defer models.SalesTrgt.Mutex.RUnlock()

	if models.SalesTrgt.SaleTargetMap[year] == nil {
		return 0, errors.New("TARGET YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	rMap, ok := models.SalesTrgt.SaleTargetMap[year][region]
	if !ok {
		return 0, errors.New("TARGET REGION NOT FOUND - " + region)
	}

	if _, ok := rMap[outlet]; !ok {
		return 0, errors.New("TARGET OUTLET NOT FOUND - " + outlet)
	}

	if isQuarterReport {
		//quarter target
		result := 0
		switch grp {
		case constant.HBBCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET HBB BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].HBB
			}
			return result, nil

		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET HAYYAK BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Hayyak
			}
			return result, nil

		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET BAQATI BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Baqati
			}
			return result, nil

		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET UPGRADE BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Upgrade
			}
			return result, nil

		}
	} else {
		//month target
		switch grp {
		case constant.HBBCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].HBB, nil
		case constant.HayyakCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Hayyak, nil
		case constant.BaqatiCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Baqati, nil
		case constant.UpgradeCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Upgrade, nil
		}
	}

	return 0, errors.New("TARGET GROUP INVALID - " + grp)
}

func (impl *UserImpl) getMaxPayout(quarter, year int, position, grp string) (float64, error) {

	pos := constant.EmployeePositionMap[position]

	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	if models.MaxPay.Payoutmap[year] == nil {
		return 0, errors.New("MAX PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxPay.Payoutmap[year][quarter] == nil {
		return 0, errors.New("MAX PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxPay.Payoutmap[year][quarter][pos] == nil {
		return 0, errors.New("MAX PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxPay.Payoutmap[year][quarter][pos][grp]
	if !ok {
		return 0, errors.New("MAX PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *UserImpl) getMaxExtraPayout(quarter, year int, position, grp string) (float64, error) {

	pos := constant.EmployeePositionMap[position]

	models.MaxExtraPay.Mutex.RLock()
	defer models.MaxExtraPay.Mutex.RUnlock()

	if models.MaxExtraPay.ExtraPayoutmap[year] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter][pos] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxExtraPay.ExtraPayoutmap[year][quarter][pos][grp]
	if !ok {
		return 0, errors.New("MAX EXTRA PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *UserImpl) getVOCCriteria(year int, achievedPercent float64) (float64, error) {

	models.VOCDetail.Mutex.RLock()
	defer models.VOCDetail.Mutex.RUnlock()

	if models.VOCDetail.Criteria[year] == nil {
		return 0, errors.New("VOC CRITERIA YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	switch {
	case achievedPercent < 85:
		return models.VOCDetail.Criteria[year][constant.VOCUnsatisfactory], nil
	case achievedPercent >= 85 && achievedPercent < 89:
		return models.VOCDetail.Criteria[year][constant.VOCAlmostAchieved], nil
	case achievedPercent >= 89 && achievedPercent < 93:
		return models.VOCDetail.Criteria[year][constant.VOCAchieved], nil
	case achievedPercent >= 93 && achievedPercent <= 97:
		return models.VOCDetail.Criteria[year][constant.VOCWellAchieved], nil
	default:
		return models.VOCDetail.Criteria[year][constant.VOCOutstanding], nil
	}

}

func (impl *UserImpl) getVOCMaxPay(year int, position string) (float64, error) {

	models.VOCDetail.Mutex.RLock()
	defer models.VOCDetail.Mutex.RUnlock()

	if models.VOCDetail.MaxPay[year] == nil {
		return 0, errors.New("VOC MAX-PAY YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	return models.VOCDetail.MaxPay[year][position], nil
}

func (impl *UserImpl) getThresholdPayout(quarter, year int, position, grp string, achievedPercent float64) (float64, error) {

	pos := constant.EmployeePositionMap[position]
	achieved := int(math.Round(achievedPercent))

	models.ThresholdCheck.Mutex.RLock()
	defer models.ThresholdCheck.Mutex.RUnlock()

	if achieved < 75 {
		return 0, nil
	}

	if achieved > 140 {
		achieved = 140
	}

	if models.ThresholdCheck.Limit[year] == nil {
		return 0, errors.New("THRESHOLD PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ThresholdCheck.Limit[year][quarter] == nil {
		return 0, errors.New("THRESHOLD PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ThresholdCheck.Limit[year][quarter][pos] == nil {
		return 0, errors.New("THRESHOLD PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	if models.ThresholdCheck.Limit[year][quarter][pos][grp] == nil {
		return 0, errors.New("THRESHOLD PAYOUT GROUP INVALID - " + fmt.Sprintf("%v", grp))
	}

	val, ok := models.ThresholdCheck.Limit[year][quarter][pos][grp][achieved]
	if !ok {
		return 0, errors.New("THRESHOLD PAYOUT ACHIEVED PERCENT INVALID - " + grp)
	}

	return val, nil
}

func (impl *UserImpl) updateAgentUserReport(ctx *gin.Context, user *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("AReport : nil info")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("AReport : invalid info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("AReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("AReport : invalid position")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "updateAgentUserReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Msg("AReport : invalid outlet info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "updateAgentUserReport: agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("AReport : fetch outlet info failed")
		return nil, nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Msg("AReport : outlet co-ordinates invalid")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "updateAgentUserReport: " + user.UserId + ": AGENT invalid outlet co-ordintaes - region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	totalEmp, err := impl.DbService.FetchTotalNonManagerialEmployees(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AReport : fetch empCount failed")
		return nil, nil, err
	}

	allUserSales, err := impl.DbService.FetchActiveUserSales(user.UserId, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AReport : fetch sales failed")
		return nil, nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("agentCount", totalEmp).Int("saleGroupCount", len(allUserSales.SaleRecords)).
		Msg("AReport : total agents and salesGroup")

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
		logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("category", grp).Int("count", qty).Msg("AReport : sale count")
	}
	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Int("count", totalActiveUpgrade).Msg("AReport : sale count")

	for grp, plans := range grpPlans {

		targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, grp, user.Region, user.Outlet)
		if internalErr != nil {
			internalErr = fmt.Errorf("updateAgentUserReport: %v", internalErr.Error())
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Err(internalErr).Msg("AReport : fetch target failed")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Message:    fmt.Sprintf("INVALID TARGET: reg: %v, out: %v, grp: %v", user.Region, user.Outlet, grp),
				Details:    internalErr.Error(),
			}
		}

		if targetPerOutlet <= 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Msg("AReport : zero target found")
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
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HBB sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HBB fetch vop failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HBB fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : HBB fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : HBB fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HBB fetch threshold failed")
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

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Float64("threshold", threshold).Msg("AReport : HBB sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)
			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", ind.Payout).Msg("AReport : HBB per plan info")
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
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HAYYAK sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HAYYAK fetch vop failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HAYYAK fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : HAYYAK fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : HAYYAK fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).Msg("AReport : HAYYAK fetch threshold failed")
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

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Float64("threshold", threshold).Msg("AReport : HAYYAK sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)
			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", ind.Payout).Msg("AReport : HAYYAK per plan info")
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
			// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("BAQATI sales error")
			// 	break
			// }

			vop, internalErr := impl.getValueOfPoint(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : BAQATI fetch vop failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			extraPayout, internalErr := impl.getExtraPay(quarter, year, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Float64("achieved", achieved).Err(internalErr).Msg("AReport : BAQATI fetch extraPayout failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : BAQATI fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : BAQATI fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, grp, achieved)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).Msg("AReport : BAQATI fetch threshold failed")
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

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Float64("threshold", threshold).Msg("AReport : BAQATI sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)

			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, _ := impl.getSIPRatePlan(quarter, year, plan)
				ind.Payout += float64(qty) * vop * rate
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", ind.Payout).Msg("AReport : BAQATI per plan info")
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

	targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory, user.Region, user.Outlet)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Err(internalErr).Msg("AReport : fetch target failed")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	if targetPerOutlet <= 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Msg("AReport : zero target found")
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
	// 	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("UPGRADE sales error")
	// 	return resp, report, nil
	// }

	vop, internalErr := impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Float64("achieved", achieved).Err(internalErr).Msg("AReport : UPGRADE fetch vop failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	extraPayout, internalErr := impl.getExtraPay(quarter, year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Float64("achieved", achieved).Err(internalErr).Msg("AReport : UPGRADE fetch extraPayout failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : UPGRADE fetch maxPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Str("position", constant.PromoterCC).Err(internalErr).Msg("AReport : UPGRADE fetch maxExtraPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		resp.VOC = &models.IndividualIncentive{}
		return resp, report, nil
	}
	threshold, internalErr := impl.getThresholdPayout(quarter, year, constant.PromoterCC, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Str("position", constant.PromoterCC).Float64("achieved", achieved).Err(internalErr).Msg("AReport : UPGRADE fetch threshold failed")
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

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
		Int("achievedCount", totalActiveUpgrade).Float64("achievedPercent", achieved).Float64("vop", vop).Float64("extraPay", extraPayout).
		Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Float64("threshold", threshold).Msg("AReport : UPGRADE sales info")

	saleCount := make([]*models.SaleCountPerPlan, 0)
	for currPlan, prevPlansMap := range upgradedPlans {
		eachPlanSale := &models.SaleCountPerPlan{CurrentPlan: currPlan}
		for prevPln, qty := range prevPlansMap {
			eachPlanSale.Quantity += qty
			currPlanRate, _ := impl.getSIPRatePlan(quarter, year, currPlan)
			prevPlanRate, _ := impl.getSIPRatePlan(quarter, year, prevPln)
			rate := currPlanRate - prevPlanRate
			ind.Payout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("currPlanName", currPlan).Str("prevPlanName", prevPln).
				Float64("currPlanRate", currPlanRate).Float64("prevPlanRate", prevPlanRate).
				Float64("netPlanRate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", ind.Payout).
				Msg("AReport : UPGRADE per plan info")
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
		if user.HBBTargetCount <= 0 || user.HayyakTargetCount <= 0 || user.BaqatiTargetCount <= 0 || user.UpgradeTargetCount <= 0 {
			return
		}

		err = impl.CacheService.InsertUserReport(isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)

		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("AReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("AReport : cache update success")
		}

		// if err := impl.DbService.InsertAgentUserReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("AReport : force report update")
		// }
	}(report)

	return resp, report, nil
}

func (impl *UserImpl) updateOutletManagerReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOutletManagerReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOutletManagerReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "updateOutletManagerReport: incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "updateOutletManagerReport: agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMReport : info")

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(omngr.Region, omngr.Outlet, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(err).Msg("OMReport : fetch outlet info failed")
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMReport : outlet co-ordinates invalid")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": OUTLET MANAGER INVALID OUTLET CO-ORDINATES - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    omngr.UserId + ": OUTLET MANAGER invalid outlet co-ordintaes - region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).Msg("OMReport : fetch agents failed")
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
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMReport : total agents")
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	for _, eachAgent := range userAgents {

		agentReport, err := impl.CacheService.GetUserReport(eachAgent.UserId, isQuarterReport, month, quarter, year)

		if !(err == nil && agentReport != nil && objecthelper.IsCurrentDay(agentReport.ModifiedAt)) {
			ind, resp, err := impl.updateAgentUserReport(ctx, eachAgent, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMReport : agent report failed")
				// return nil, err
				continue
			}
			rawIncentiveReport, _ := json.Marshal(ind)
			rawUserReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("incentive", string(rawIncentiveReport)).
				Str("record", string(rawUserReport)).Msg("OMReport : agent log")
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
	target, tarErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HBBCategory, omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(tarErr).Msg("OMReport : HBB target not found")
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

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.HBBSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("OMReport : HBB sales info")

	for plan, qty := range hbbSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.HBBPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HBBPayout).Msg("OMReport : HBB per plan info")
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
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HayyakCategory, omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(tarErr).Msg("OMReport : HAYYAK target not found")
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

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.HayyakSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("OMReport : HAYYAK sales info")

	for plan, qty := range hayyakSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.HayyakPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HayyakPayout).Msg("OMReport : HAYYAK per plan info")
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
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.BaqatiCategory, omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(tarErr).Msg("OMReport : BAQATI target not found")
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

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.BaqatiSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("OMReport : BAQATI sales info")

	for plan, qty := range baqatiSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.BaqatiPayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.BaqatiPayout).Msg("OMReport : BAQATI per plan info")
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
	target, tarErr = impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory, omngr.Region, omngr.Outlet)
	if tarErr != nil || target == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(tarErr).Msg("OMReport : UPGRADE target not found")
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

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.UpgradeSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("OMReport : UPGRADE sales info")

	for plan, qty := range upgradeSaleMap {
		rate, _ := impl.getSIPRatePlan(quarter, year, plan)
		report.UpgradePayout += float64(qty) * vop * rate
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("planName", plan).Float64("planRate", rate).
			Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.UpgradePayout).Msg("OMReport : UPGRADE per plan info")
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

		err = impl.CacheService.InsertUserReport(isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("OMReport : cache update failed")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("OMReport : cache update success")
		}

		// if err := impl.DbService.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("OMReport : force report update")
		// }
	}(report)

	return report, nil
}

func (impl *UserImpl) updateManagerReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateManagerReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "updateManagerReport: incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerReport(ctx, mngr, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrReport : user's manager info missing")
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
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Err(errors.New(err.Details.(string))).Msg("MngrReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrReport : zero om")
		return report, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).Msg("MngrReport : total om")
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	inactiveOM := make([]string, 0)
	emptyTarget := make(map[string][]string, 0)

	for _, om := range allDistinctOM {

		omReport, err := impl.CacheService.GetUserReport(om.UserId, isQuarterReport, month, quarter, year)

		if !(err == nil && omReport != nil && objecthelper.IsCurrentDay(omReport.ModifiedAt)) {
			resp, err := impl.updateOutletManagerReport(ctx, om, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("om", om.UserId).Err(err).Msg("MngrReport : om report failed")
				continue
				//return nil, err
			}
			rawUserReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", om.UserId).Str("record", string(rawUserReport)).Msg("OMReport : agent log")
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

	logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Interface("inactive", inactiveOM).Interface("blank", emptyTarget).Msg("MngrReport : om target defaulters")

	/*
		HBB Updates
	*/
	if report.HBBTargetCount == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrReport : HBB target not found")
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
		threshold, _ := impl.getThresholdPayout(quarter, year, constant.ManagerCC, constant.HBBCategory, report.HBBAchievedPercentage)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HBBTargetCount).Int("achievedCount", report.HBBSalesCount).
			Float64("achievedPercent", report.HBBAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("MngrReport : HBB sales info")
		for plan, qty := range hbbSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HBBPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HBBPayout).Msg("MngrReport : HBB per plan info")
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
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager, constant.HayyakCategory, report.HayyakAchievedPercentage)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HayyakTargetCount).Int("achievedCount", report.HayyakSalesCount).
			Float64("achievedPercent", report.HayyakAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("MngrReport : HAYYAK sales info")
		for plan, qty := range hayyakSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HayyakPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HayyakPayout).Msg("MngrReport : HAYYAK per plan info")
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
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager, constant.BaqatiCategory, report.BaqatiAchievedPercentage)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.BaqatiTargetCount).Int("achievedCount", report.BaqatiSalesCount).
			Float64("achievedPercent", report.BaqatiAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("MngrReport : BAQATI sales info")
		for plan, qty := range baqatiSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.BaqatiPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.BaqatiPayout).Msg("MngrReport : BAQATI per plan info")
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
		threshold, _ = impl.getThresholdPayout(quarter, year, constant.Manager, constant.UpgradeCategory, report.UpradeAchievedPercentage)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.UpgradeTargetCount).Int("achievedCount", report.UpgradeSalesCount).
			Float64("achievedPercent", report.UpradeAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Float64("threshold", threshold).Msg("MngrReport : UPGRADE sales info")
		for plan, qty := range upgradeSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.UpgradePayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.UpgradePayout).Msg("MngrReport : UPGRADE per plan info")
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

		err = impl.CacheService.InsertUserReport(isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("MngrReport : cache update failed")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("MngrReport : cache update success")
		}

		// if err := impl.DbService.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
		// 	logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("MngrReport : force report update")
		// }
	}(report)

	return report, nil
}

func (impl *UserImpl) updateOutletReport(ctx *gin.Context, regionName, outletName string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if regionName == "" || outletName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION OR OUTLET - region:" + regionName + " outlet:" + outletName,
			Details:    "blank region or outlet - region:" + regionName + " outlet:" + outletName,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Str("outlet", outletName).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("update Outlet Report")

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(regionName, outletName, quarter, year)
	if err != nil {
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "INVALID OUTLET CO-ORDINATES - region:" + regionName + " outlet:" + outletName,
			Details:    "invalid outlet co-ordintaes - region:" + regionName + " outlet:" + outletName,
		}
	}

	hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.DbService.FetchOutletSalesCount(regionName, outletName, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	hbbTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HBBCategory, regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	hayyakTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HayyakCategory, regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	baqatiTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.BaqatiCategory, regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	upgradeTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory, regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	outletReport := &models.OutletReportDB{
		Outlet:             outletName,
		Region:             regionName,
		Latitude:           outletResponse.Lat,
		Longitude:          outletResponse.Lang,
		HBBSalesCount:      hbbSales,
		HayyakSalesCount:   hayyakSales,
		BaqatiSalesCount:   baqatiSales,
		UpgradeSalesCount:  upgradeSales,
		HBBTargetCount:     hbbTarget,
		HayyakTargetCount:  hayyakTarget,
		BaqatiTargetCount:  baqatiTarget,
		UpgradeTargetCount: upgradeTarget,
		CreatedAt:          time.Now(),
		ModifiedAt:         time.Now(),
	}

	go func(report *models.OutletReportDB) {

		err = impl.CacheService.InsertOutletReport(isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("outlet", outletName).Str("region", regionName).Err(errors.New(err.Details.(string))).Msg("OutletReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("outlet", outletName).Msg("OutletReport : cache update success")
		}

	}(outletReport)

	return outletReport, nil
}

func (impl *UserImpl) updateRegionReport(ctx *gin.Context, regionName string, isQuarterReport bool, quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if regionName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION NAME",
			Details:    "blank region name",
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("update Outlet Report")

	allOutlets, err := impl.DbService.FetchAllOutletNames(regionName, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allOutlets) == 0 {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "NO OUTLETS FOUND FOR REGION - " + regionName,
			Details:    "no outlet found for region - " + regionName,
		}
	}

	regionReport := &models.RegionReportDB{
		Region:     regionName,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	for _, eachOutlet := range allOutlets {

		var (
			err          *errors.Error
			outletReport *models.OutletReportDB
		)

		outletReport, err = impl.CacheService.GetOutletReport(eachOutlet, regionName, isQuarterReport, month, quarter, year)
		if !(err == nil && outletReport != nil && objecthelper.IsCurrentDay(outletReport.ModifiedAt)) {
			outletResp, err := impl.updateOutletReport(ctx, regionName, eachOutlet, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("region", regionName).Str("outlet", eachOutlet).Bool("isQuarter", isQuarterReport).
					Int("quarter", quarter).Int("month", month).Int("year", year).Err(err).Msg("RegionReport : outlet report failed")
				//return nil, err
				continue
			}

			logger.FileLogger.CtxInfo(ctx).Interface("outletReport", outletResp).Msg("updateRegionReport : trace")
			//TODO: test this out for struct assignment
			outletReport = outletResp
		}

		regionReport.HayyakSalesCount += outletReport.HayyakSalesCount
		regionReport.HBBSalesCount += outletReport.HBBSalesCount
		regionReport.BaqatiSalesCount += outletReport.BaqatiSalesCount
		regionReport.UpgradeSalesCount += outletReport.UpgradeSalesCount

		regionReport.HayyakTargetCount += outletReport.HayyakTargetCount
		regionReport.HBBTargetCount += outletReport.HBBTargetCount
		regionReport.BaqatiTargetCount += outletReport.BaqatiTargetCount
		regionReport.UpgradeTargetCount += outletReport.UpgradeTargetCount
	}

	go func(report *models.RegionReportDB) {

		err = impl.CacheService.InsertRegionReport(isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("region", regionName).Err(errors.New(err.Details.(string))).Msg("RegionReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Msg("RegionReport` : cache update success")
		}

	}(regionReport)

	return regionReport, nil
}

func (impl *UserImpl) updateAgentUserSalesCountReport(ctx *gin.Context, user *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserSalesCountReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserSalesCountReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("ASalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "ASalesCountReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).
			Msg("ASalesCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "ASalesCountReport: agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchUserSalesCountPerPlan(user.UserId, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserSalesCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("ASalesCountReport : fetch sales failed")
		return nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales)).
		Msg("ASalesCountReport : total agents and salesGroup")

	if product == constant.HBBCategory {
		response := models.HbbUserSalesCountPerPlanDB{
			Product:    constant.HBBCategory,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, v := range allUserSales {
			if v.ProductGroup != constant.HBBCategory {
				continue
			}
			switch v.BusinessUnit {
			case "GOVERNMENT":
				if response.BusinessUnit.Government == nil {
					response.BusinessUnit.Government = &models.SaleCalculation{}
				}
				response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
			case "CORPORATE":
				if response.BusinessUnit.Corporate == nil {
					response.BusinessUnit.Corporate = &models.SaleCalculation{}
				}
				response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					if response.Category.CBUWFBB == nil {
						response.Category.CBUWFBB = &models.SaleCalculation{}
					}
					response.Category.CBUWFBB.Count = response.Category.CBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUWFBB == nil {
						response.SubCategory.CBUWFBB = make(map[string]interface{})
					}
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					total := response.SubCategory.CBUWFBB["total"].(int)
					response.SubCategory.CBUWFBB["total"] = total + v.Quantity
					if response.SubCategory.CBUWFBB[v.SubCategory] == nil {
						response.SubCategory.CBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.CBUWFBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.CBUWFBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.CBUWFBB == nil {
						response.ServiceCategory.CBUWFBB = make(map[string]interface{})
					}
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					total = response.ServiceCategory.CBUWFBB["total"].(int)
					response.ServiceCategory.CBUWFBB["total"] = total + v.Quantity

					//plans
					if response.Plans.CBUWFBB == nil {
						response.Plans.CBUWFBB = make(map[string]int)
					}
					total = response.Plans.CBUWFBB["total"]
					response.Plans.CBUWFBB["total"] = total + v.Quantity
					total = response.Plans.CBUWFBB[v.CurrentPlan]
					response.Plans.CBUWFBB[v.CurrentPlan] = total + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" || v.LineType == "PSTN VOICE SERVICE" {
					//category
					if response.Category.CBUHBB == nil {
						response.Category.CBUHBB = &models.SaleCalculation{}
					}
					response.Category.CBUHBB.Count = response.Category.CBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUHBB == nil {
						response.SubCategory.CBUHBB = make(map[string]interface{})
					}
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}

					total := response.SubCategory.CBUHBB["total"].(int)
					response.SubCategory.CBUHBB["total"] = total + v.Quantity
					if response.SubCategory.CBUHBB[v.SubCategory] == nil {
						response.SubCategory.CBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.CBUHBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.CBUHBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.CBUHBB == nil {
						response.ServiceCategory.CBUHBB = make(map[string]interface{})
					}
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					total = response.ServiceCategory.CBUHBB["total"].(int)
					response.ServiceCategory.CBUHBB["total"] = total + v.Quantity
					if response.ServiceCategory.CBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.CBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					srvcCat := response.ServiceCategory.CBUHBB[v.FinalServiceCategory].(*models.SaleCalculation)
					srvcCat.Count = srvcCat.Count + v.Quantity
					response.ServiceCategory.CBUHBB[v.FinalServiceCategory] = srvcCat

					//plans
					if response.Plans.CBUHBB == nil {
						response.Plans.CBUHBB = make(map[string]int)
					}
					count := response.Plans.CBUHBB["total"]
					response.Plans.CBUHBB["total"] = count + v.Quantity
					count = response.Plans.CBUHBB[v.CurrentPlan]
					response.Plans.CBUHBB[v.CurrentPlan] = count + v.Quantity
				}
			case "CONSUMER":
				if response.BusinessUnit.Consumer == nil {
					response.BusinessUnit.Consumer = &models.SaleCalculation{}
				}
				response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					if response.Category.EBUWFBB == nil {
						response.Category.EBUWFBB = &models.SaleCalculation{}
					}
					response.Category.EBUWFBB.Count = response.Category.EBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUWFBB == nil {
						response.SubCategory.EBUWFBB = make(map[string]interface{})
					}
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					total := response.SubCategory.EBUWFBB["total"].(int)
					response.SubCategory.EBUWFBB["total"] = total + v.Quantity
					if response.SubCategory.EBUWFBB[v.SubCategory] == nil {
						response.SubCategory.EBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.EBUWFBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.EBUWFBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.EBUWFBB == nil {
						response.ServiceCategory.EBUWFBB = make(map[string]interface{})
					}
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					total = response.ServiceCategory.EBUWFBB["total"].(int)
					response.ServiceCategory.EBUWFBB["total"] = total + v.Quantity

					//plans
					if response.Plans.EBUWFBB == nil {
						response.Plans.EBUWFBB = make(map[string]int)
					}
					total = response.Plans.EBUWFBB["total"]
					response.Plans.EBUWFBB["total"] = total + v.Quantity
					total = response.Plans.EBUWFBB[v.CurrentPlan]
					response.Plans.EBUWFBB[v.CurrentPlan] = total + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" || v.LineType == "PSTN VOICE SERVICE" {
					//category
					if response.Category.EBUHBB == nil {
						response.Category.EBUHBB = &models.SaleCalculation{}
					}
					response.Category.EBUHBB.Count = response.Category.EBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUHBB == nil {
						response.SubCategory.EBUHBB = make(map[string]interface{})
					}
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					total := response.SubCategory.EBUHBB["total"].(int)
					response.SubCategory.EBUHBB["total"] = total + v.Quantity
					if response.SubCategory.EBUHBB[v.SubCategory] == nil {
						response.SubCategory.EBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.EBUHBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.EBUHBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.EBUHBB == nil {
						response.ServiceCategory.EBUHBB = make(map[string]interface{})
					}
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					total = response.ServiceCategory.EBUHBB["total"].(int)
					response.ServiceCategory.EBUHBB["total"] = total + v.Quantity
					if response.ServiceCategory.EBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.EBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					srvcCat := response.ServiceCategory.EBUHBB[v.FinalServiceCategory].(*models.SaleCalculation)
					srvcCat.Count = srvcCat.Count + v.Quantity
					response.ServiceCategory.EBUHBB[v.FinalServiceCategory] = srvcCat

					//plans
					if response.Plans.EBUHBB == nil {
						response.Plans.EBUHBB = make(map[string]int)
					}
					count := response.Plans.EBUHBB["total"]
					response.Plans.EBUHBB["total"] = count + v.Quantity
					count = response.Plans.EBUHBB[v.CurrentPlan]
					response.Plans.EBUHBB[v.CurrentPlan] = count + v.Quantity
				}
			default:
			}
		}

		//businesssUnit
		if response.BusinessUnit.Government == nil {
			response.BusinessUnit.Government = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Corporate == nil {
			response.BusinessUnit.Corporate = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Consumer == nil {
			response.BusinessUnit.Consumer = &models.SaleCalculation{}
		}
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		if response.Category.CBUWFBB == nil {
			response.Category.CBUWFBB = &models.SaleCalculation{}
		}
		if response.Category.CBUHBB == nil {
			response.Category.CBUHBB = &models.SaleCalculation{}
		}
		if response.Category.EBUWFBB == nil {
			response.Category.EBUWFBB = &models.SaleCalculation{}
		}
		if response.Category.EBUHBB == nil {
			response.Category.EBUHBB = &models.SaleCalculation{}
		}
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		if response.SubCategory.CBUWFBB == nil {
			response.SubCategory.CBUWFBB = make(map[string]interface{})
		}
		if response.SubCategory.CBUHBB == nil {
			response.SubCategory.CBUHBB = make(map[string]interface{})
		}
		if response.SubCategory.EBUWFBB == nil {
			response.SubCategory.EBUWFBB = make(map[string]interface{})
		}
		if response.SubCategory.EBUHBB == nil {
			response.SubCategory.EBUHBB = make(map[string]interface{})
		}
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		if response.ServiceCategory.CBUWFBB == nil {
			response.ServiceCategory.CBUWFBB = make(map[string]interface{})
		}
		if response.ServiceCategory.CBUHBB == nil {
			response.ServiceCategory.CBUHBB = make(map[string]interface{})
		}
		if response.ServiceCategory.EBUWFBB == nil {
			response.ServiceCategory.EBUWFBB = make(map[string]interface{})
		}
		if response.ServiceCategory.EBUHBB == nil {
			response.ServiceCategory.EBUHBB = make(map[string]interface{})
		}
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("ASalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).Msg("ASalesCountReport : cache update success")
			}
		}(&response)

		return response, nil
	} else {
		response := models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, v := range allUserSales {
			if v.ProductGroup != product {
				continue
			}
			switch v.BusinessUnit {
			case "GOVERNMENT":
				if response.BusinessUnit.Government == nil {
					response.BusinessUnit.Government = &models.SaleCalculation{}
				}
				response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
			case "CORPORATE":
				if response.BusinessUnit.Corporate == nil {
					response.BusinessUnit.Corporate = &models.SaleCalculation{}
				}
				response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

				//subcategory
				if response.SubCategory.CBU == nil {
					response.SubCategory.CBU = make(map[string]interface{})
				}
				if response.SubCategory.CBU["total"] == nil {
					response.SubCategory.CBU["total"] = 0
				}
				total := response.SubCategory.CBU["total"].(int)
				total = total + v.Quantity
				response.SubCategory.CBU["total"] = total
				if response.SubCategory.CBU[v.SubCategory] == nil {
					response.SubCategory.CBU[v.SubCategory] = &models.SaleCalculation{}
				}
				subcategory := response.SubCategory.CBU[v.SubCategory].(*models.SaleCalculation)
				subcategory.Count = subcategory.Count + v.Quantity
				response.SubCategory.CBU[v.SubCategory] = subcategory

				//plans
				if response.Plans.CBU == nil {
					response.Plans.CBU = make(map[string]int)
				}
				total = response.Plans.CBU["total"]
				response.Plans.CBU["total"] = total + v.Quantity
				total = response.Plans.CBU[v.CurrentPlan]
				response.Plans.CBU[v.CurrentPlan] = total + v.Quantity
			case "CONSUMER":
				if response.BusinessUnit.Consumer == nil {
					response.BusinessUnit.Consumer = &models.SaleCalculation{}
				}
				response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

				//subcategory
				if response.SubCategory.EBU == nil {
					response.SubCategory.EBU = make(map[string]interface{})
				}
				if response.SubCategory.EBU["total"] == nil {
					response.SubCategory.EBU["total"] = 0
				}
				total := response.SubCategory.EBU["total"].(int)
				total = total + v.Quantity
				response.SubCategory.EBU["total"] = total
				if response.SubCategory.EBU[v.SubCategory] == nil {
					response.SubCategory.EBU[v.SubCategory] = &models.SaleCalculation{}
				}
				subcategory := response.SubCategory.EBU[v.SubCategory].(*models.SaleCalculation)
				subcategory.Count = subcategory.Count + v.Quantity
				response.SubCategory.EBU[v.SubCategory] = subcategory

				//plans
				if response.Plans.EBU == nil {
					response.Plans.EBU = make(map[string]int)
				}
				total = response.Plans.EBU["total"]
				response.Plans.EBU["total"] = total + v.Quantity
				total = response.Plans.EBU[v.CurrentPlan]
				response.Plans.EBU[v.CurrentPlan] = total + v.Quantity
			default:
			}
		}

		//businesssUnit
		if response.BusinessUnit.Government == nil {
			response.BusinessUnit.Government = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Corporate == nil {
			response.BusinessUnit.Corporate = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Consumer == nil {
			response.BusinessUnit.Consumer = &models.SaleCalculation{}
		}
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		if response.SubCategory.CBU == nil {
			response.SubCategory.CBU = make(map[string]interface{})
		}
		if response.SubCategory.EBU == nil {
			response.SubCategory.EBU = make(map[string]interface{})
		}

		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("ASalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).Msg("ASalesCountReport : cache update success")
			}
		}(&response)

		return response, nil
	}
}

func (impl *UserImpl) updateOutletManagerUserSaleCountReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMSalesCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMUserSalesReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMUserSalesReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "updateOMUserSalesReport: incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMSalesCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "updateOMUserSalesReport: agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMSalesCountReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).Msg("OMSalesCountReport : fetch agents failed")
		return nil, err
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		if len(userAgents) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
			return response, nil
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMSalesCountReport : total agents")
		}
		for _, eachAgent := range userAgents {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
			if !(err == nil && agentReport != nil) {
				userSaleCount, _ := agentReport.(*models.HbbUserSalesCountPerPlanDB)
				if !objecthelper.IsCurrentDay(userSaleCount.ModifiedAt) {
					resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(resp)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
					agentReport = resp
				}
			}
			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.HbbUserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.Category.CBUWFBB == nil {
				response.Category.CBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.CBUHBB == nil {
				response.Category.CBUHBB = &models.SaleCalculation{}
			}
			if response.Category.EBUWFBB == nil {
				response.Category.EBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.EBUHBB == nil {
				response.Category.EBUHBB = &models.SaleCalculation{}
			}
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					response.SubCategory.CBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.CBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}
					response.SubCategory.CBUHBB["total"] = val.(int) + userSaleCount.SubCategory.CBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUHBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					response.SubCategory.EBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.EBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					response.SubCategory.EBUHBB["total"] = val.(int) + userSaleCount.SubCategory.EBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUHBB[subcategory] = reportSubCat
			}

			//serviceCategory
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					response.ServiceCategory.CBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					response.ServiceCategory.CBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.CBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.CBUHBB[srvccategory] = reportSubCat
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					response.ServiceCategory.EBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					response.ServiceCategory.EBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.EBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.EBUHBB[srvccategory] = reportSubCat
			}

			//plans
			if userSaleCount.Plans.CBUWFBB == nil {
				userSaleCount.Plans.CBUWFBB = make(map[string]int)
			}
			if userSaleCount.Plans.CBUHBB == nil {
				userSaleCount.Plans.CBUHBB = make(map[string]int)
			}
			if userSaleCount.Plans.EBUWFBB == nil {
				userSaleCount.Plans.EBUWFBB = make(map[string]int)
			}
			if userSaleCount.Plans.EBUHBB == nil {
				userSaleCount.Plans.EBUHBB = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBUWFBB {
				response.Plans.CBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.CBUHBB {
				response.Plans.CBUHBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUWFBB {
				response.Plans.EBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUHBB {
				response.Plans.EBUHBB[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.ServiceCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.ServiceCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("OMSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	} else {
		response := &models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		if len(userAgents) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
			return response, nil
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMSalesCountReport : total agents")
		}
		for _, eachAgent := range userAgents {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
			if !(err == nil && agentReport != nil && objecthelper.IsCurrentDay(response.ModifiedAt)) {
				resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
				agentReport = resp
			}
			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.UserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//subcategory
			if response.SubCategory.CBU == nil {
				response.SubCategory.CBU = make(map[string]interface{})
			}
			if response.SubCategory.EBU == nil {
				response.SubCategory.EBU = make(map[string]interface{})
			}
			for subcategory, val := range userSaleCount.SubCategory.CBU {
				if subcategory == "total" {
					if response.SubCategory.CBU["total"] == nil {
						response.SubCategory.CBU["total"] = 0
					}
					response.SubCategory.CBU["total"] = val.(int) + userSaleCount.SubCategory.CBU["total"].(int)
					continue
				}
				if response.SubCategory.CBU[subcategory] == nil {
					response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBU[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBU {
				if subcategory == "total" {
					if response.SubCategory.EBU["total"] == nil {
						response.SubCategory.EBU["total"] = 0
					}
					response.SubCategory.EBU["total"] = val.(int) + userSaleCount.SubCategory.EBU["total"].(int)
					continue
				}
				if response.SubCategory.EBU[subcategory] == nil {
					response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBU[subcategory] = reportSubCat
			}

			//plan
			if response.Plans.CBU == nil {
				response.Plans.CBU = make(map[string]int)
			}
			if response.Plans.EBU == nil {
				response.Plans.EBU = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBU {
				response.Plans.CBU[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBU {
				response.Plans.EBU[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("OMSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	}
}

func (impl *UserImpl) updateManagerUserSalesCountReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateMngrUserSalesCountReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "updateMngrUserSalesCountReport: incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrSalesCountReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserSaleCountReport(ctx, mngr, product, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrSalesCountReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Err(errors.New(err.Details.(string))).Msg("MngrSalesCountReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrSalesCountReport : zero om")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).Msg("MngrSalesCountReport : total om")
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, eachAgent := range allDistinctOM {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
			if !(err == nil && agentReport != nil && objecthelper.IsCurrentDay(response.ModifiedAt)) {
				resp, err := impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
				agentReport = resp
			}
			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.HbbUserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.Category.CBUWFBB == nil {
				response.Category.CBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.CBUHBB == nil {
				response.Category.CBUHBB = &models.SaleCalculation{}
			}
			if response.Category.EBUWFBB == nil {
				response.Category.EBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.EBUHBB == nil {
				response.Category.EBUHBB = &models.SaleCalculation{}
			}
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					response.SubCategory.CBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.CBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}
					response.SubCategory.CBUHBB["total"] = val.(int) + userSaleCount.SubCategory.CBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUHBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					response.SubCategory.EBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.EBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					response.SubCategory.EBUHBB["total"] = val.(int) + userSaleCount.SubCategory.EBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUHBB[subcategory] = reportSubCat
			}

			//serviceCategory
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					response.ServiceCategory.CBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					response.ServiceCategory.CBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.CBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.CBUHBB[srvccategory] = reportSubCat
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					response.ServiceCategory.EBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					response.ServiceCategory.EBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.EBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.EBUHBB[srvccategory] = reportSubCat
			}

			//plans
			if userSaleCount.Plans.CBUWFBB == nil {
				userSaleCount.Plans.CBUWFBB = make(map[string]int)
			}
			if userSaleCount.Plans.CBUHBB == nil {
				userSaleCount.Plans.CBUHBB = make(map[string]int)
			}
			if userSaleCount.Plans.EBUWFBB == nil {
				userSaleCount.Plans.EBUWFBB = make(map[string]int)
			}
			if userSaleCount.Plans.EBUHBB == nil {
				userSaleCount.Plans.EBUHBB = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBUWFBB {
				response.Plans.CBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.CBUHBB {
				response.Plans.CBUHBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUWFBB {
				response.Plans.EBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUHBB {
				response.Plans.EBUHBB[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.ServiceCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("MngrSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	} else {
		response := &models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, eachAgent := range allDistinctOM {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
			if !(err == nil && agentReport != nil && objecthelper.IsCurrentDay(response.ModifiedAt)) {
				resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
				agentReport = resp
			}
			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.UserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.SubCategory.CBU == nil {
				response.SubCategory.CBU = make(map[string]interface{})
			}
			if response.SubCategory.EBU == nil {
				response.SubCategory.EBU = make(map[string]interface{})
			}
			for subcategory, val := range userSaleCount.SubCategory.CBU {
				if subcategory == "total" {
					if response.SubCategory.CBU["total"] == nil {
						response.SubCategory.CBU["total"] = 0
					}
					response.SubCategory.CBU["total"] = val.(int) + userSaleCount.SubCategory.CBU["total"].(int)
					continue
				}
				if response.SubCategory.CBU[subcategory] == nil {
					response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBU[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBU {
				if subcategory == "total" {
					if response.SubCategory.EBU["total"] == nil {
						response.SubCategory.EBU["total"] = 0
					}
					response.SubCategory.EBU["total"] = val.(int) + userSaleCount.SubCategory.EBU["total"].(int)
					continue
				}
				if response.SubCategory.EBU[subcategory] == nil {
					response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBU[subcategory] = reportSubCat
			}

			//plan
			if response.Plans.CBU == nil {
				response.Plans.CBU = make(map[string]int)
			}
			if response.Plans.EBU == nil {
				response.Plans.EBU = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBU {
				response.Plans.CBU[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBU {
				response.Plans.EBU[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("MngrSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	}
}

func (impl *UserImpl) updateAgentUserPlanPerCategoryReport(ctx *gin.Context, user *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserPlanPerCategoryReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserPlanPerCategoryReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("APlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "APlanPerCategoryReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).
			Msg("APlanPerCategoryReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "APlanPerCategoryReport: agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchUserSalesCountPerPlan(user.UserId, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserPlanPerCategoryReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("APlanPerCategoryReport : fetch sales failed")
		return nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales)).
		Msg("APlanPerCategoryReport : total agents and salesGroup")

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)
	for _, v := range allUserSales {
		if v.ProductGroup != constant.HBBCategory {
			continue
		}

		switch {
		case product == "CBU" && v.BusinessUnit == "CORPORATE":
		case product == "EBU" && v.BusinessUnit == "CONSUMER":
		default:
			continue
		}

		//only consider HBB (i.e CBU-HBB and EBU-HBB) and ignore WFBB (i.e CBU-WFBB and EBU-WFBB)
		if !(v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" || v.LineType == "PSTN VOICE SERVICE") {
			continue
		}

		if response[v.SubCategory] == nil {
			response[v.SubCategory] = make(map[string]interface{})
		}
		plans := response[v.SubCategory].(map[string]interface{})

		if plans[v.CurrentPlan] == nil {
			plans[v.CurrentPlan] = make(map[string]int)
		}
		srvccategory := plans[v.CurrentPlan].(map[string]int)
		srvccategory[v.FinalServiceCategory] = srvccategory[v.FinalServiceCategory] + v.Quantity

		//total
		srvccategory["total"] = srvccategory["total"] + v.Quantity

		totalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
		totalSrvcCategory[v.FinalServiceCategory] = totalSrvcCategory[v.FinalServiceCategory] + v.Quantity
		totalSrvcCategory["total"] = totalSrvcCategory["total"] + v.Quantity
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(user.UserId, product, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("APlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).Msg("APlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}

func (impl *UserImpl) updateOutletManagerUserPlanPerCategoryReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMUserPlanPerCategoryReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMUserPlanPerCategoryReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMUserPlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "updateOMUserPlanPerCategoryReport: incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMUserPlanPerCategoryReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "updateOMUserPlanPerCategoryReport: agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).Msg("OMUserPlanPerCategoryReport : fetch agents failed")
		return nil, err
	}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMUserPlanPerCategoryReport : zero agents")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMUserPlanPerCategoryReport : total agents")
	}

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)

	for _, eachAgent := range userAgents {
		agentReport, err := impl.CacheService.GetUserPlanCategoryReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
		if !(err == nil && agentReport != nil) {
			report, _ := agentReport.(map[string]interface{})
			modifiedAt, _ := report["ModifiedAt"].(time.Time)
			if objecthelper.IsCurrentDay(modifiedAt) {
				resp, err := impl.updateAgentUserPlanPerCategoryReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMUserPlanPerCategoryReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMUserPlanPerCategoryReport : agent log")
				agentReport = resp
			}
		}
		if agentReport == nil {
			continue
		}
		userPlanCount := agentReport.(map[string]interface{})

		for subcategoryName, subcategoryMap := range userPlanCount {

			if subcategoryName == "GRAND TOTAL" {
				agentGrandTotal, ok := subcategoryMap.(map[string]int)
				if !ok {
					continue
				}
				for name, qty := range agentGrandTotal {
					responseGrandTotalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
					responseGrandTotalSrvcCategory[name] = responseGrandTotalSrvcCategory[name] + qty
					response["GRAND TOTAL"] = responseGrandTotalSrvcCategory
				}
				continue
			}

			if response[subcategoryName] == nil {
				response[subcategoryName] = make(map[string]interface{})
			}

			planMap, _ := subcategoryMap.(map[string]interface{})
			responsePlanMap, _ := response[subcategoryName].(map[string]interface{})

			for planName, plan := range planMap {
				if responsePlanMap[planName] == nil {
					responsePlanMap[planName] = make(map[string]int)
				}

				srvccategoryMap, _ := plan.(map[string]int)
				responseSrvccategoryMap, _ := responsePlanMap[planName].(map[string]int)

				for srvcCategoryName, qty := range srvccategoryMap {
					responseSrvccategoryMap[srvcCategoryName] = responseSrvccategoryMap[srvcCategoryName] + qty
				}
			}
		}
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(omngr.UserId, product, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("OMUserPlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMUserPlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}

func (impl *UserImpl) updateManagerUserPlanPerCategoryReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, product string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrUserPlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateMngrUserPlanPerCategoryReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrUserPlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "updateMngrUserPlanPerCategoryReport: incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrUserPlanPerCategoryReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserPlanPerCategoryReport(ctx, mngr, product, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrUserPlanPerCategoryReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Err(errors.New(err.Details.(string))).Msg("MngrUserPlanPerCategoryReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrUserPlanPerCategoryReport : zero om")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).Msg("MngrUserPlanPerCategoryReport : total om")
	}

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)

	for _, eachAgent := range allDistinctOM {
		agentReport, err := impl.CacheService.GetUserPlanCategoryReport(eachAgent.UserId, product, isQuarterReport, month, quarter, year)
		if !(err == nil && agentReport != nil) {
			report, _ := agentReport.(map[string]interface{})
			modifiedAt, _ := report["ModifiedAt"].(time.Time)
			if objecthelper.IsCurrentDay(modifiedAt) {
				resp, err := impl.updateOutletManagerUserPlanPerCategoryReport(ctx, eachAgent, product, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrUserPlanPerCategoryReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrUserPlanPerCategoryReport : agent log")
				agentReport = resp
			}
		}
		if agentReport == nil {
			continue
		}
		userPlanCount := agentReport.(map[string]interface{})

		for subcategoryName, subcategoryMap := range userPlanCount {

			if subcategoryName == "GRAND TOTAL" {
				agentGrandTotal, ok := subcategoryMap.(map[string]int)
				if !ok {
					continue
				}
				for name, qty := range agentGrandTotal {
					responseGrandTotalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
					responseGrandTotalSrvcCategory[name] = responseGrandTotalSrvcCategory[name] + qty
					response["GRAND TOTAL"] = responseGrandTotalSrvcCategory
				}
				continue
			}

			if response[subcategoryName] == nil {
				response[subcategoryName] = make(map[string]interface{})
			}

			planMap, _ := subcategoryMap.(map[string]interface{})
			responsePlanMap, _ := response[subcategoryName].(map[string]interface{})

			for planName, plan := range planMap {
				if responsePlanMap[planName] == nil {
					responsePlanMap[planName] = make(map[string]int)
				}

				srvccategoryMap, _ := plan.(map[string]int)
				responseSrvccategoryMap, _ := responsePlanMap[planName].(map[string]int)

				for srvcCategoryName, qty := range srvccategoryMap {
					responseSrvccategoryMap[srvcCategoryName] = responseSrvccategoryMap[srvcCategoryName] + qty
				}
			}
		}
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(mngr.UserId, product, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("MngrUserPlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrUserPlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}
