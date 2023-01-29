package business

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

type IncentiveImpl struct {
	Config *configs.IncentiveConfig
	PsqlDb *psqldb.DBImplementation
}

func NewIncentiveService(config *configs.IncentiveConfig, dbConfig *models.DbClientConfig) *IncentiveImpl {

	return &IncentiveImpl{
		Config: config,
		PsqlDb: psqldb.NewDatabaseSetup(dbConfig),
	}
}

func (impl *IncentiveImpl) getValueOfPoint(quarter, year int, grp string, achievedPercent float64) (float64, error) {
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

func (impl *IncentiveImpl) getExtraPay(quarter, year int, grp string, achievedPercent float64) (float64, error) {
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

func (impl *IncentiveImpl) getSIPRatePlan(quarter, year int, planName string) (float64, error) {

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

func (impl *IncentiveImpl) getSaleTargetPerOutlet(year, month, quarter int, isQuarterReport bool, grp, region, outlet string) (int, error) {

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
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].HBB
			}
			return result, nil

		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Hayyak
			}
			return result, nil

		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Baqati
			}
			return result, nil

		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
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

func (impl *IncentiveImpl) getMaxPayout(quarter, year int, position, grp string) (float64, error) {
	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	if models.MaxPay.Payoutmap[year] == nil {
		return 0, errors.New("MAX PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxPay.Payoutmap[year][quarter] == nil {
		return 0, errors.New("MAX PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxPay.Payoutmap[year][quarter][position] == nil {
		return 0, errors.New("MAX PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxPay.Payoutmap[year][quarter][position][grp]
	if !ok {
		return 0, errors.New("MAX PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *IncentiveImpl) getMaxExtraPayout(quarter, year int, position, grp string) (float64, error) {
	models.MaxExtraPay.Mutex.RLock()
	defer models.MaxExtraPay.Mutex.RUnlock()

	if models.MaxExtraPay.ExtraPayoutmap[year] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter][position] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxExtraPay.ExtraPayoutmap[year][quarter][position][grp]
	if !ok {
		return 0, errors.New("MAX EXTRA PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *IncentiveImpl) getQuarterFromMonth(month int) int {
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	case month >= 10 && month <= 12:
		return 4
	}

	return -1
}

func (impl *IncentiveImpl) updateAgentUserReport(ctx *gin.Context, user *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("AReport : nil info")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("AReport : invalid info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("AReport : info")

	if !(user.Position == constant.Promoter || user.Position == constant.CSR) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("AReport : invalid position")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Msg("AReport : invalid outlet info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("AReport : info")

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AReport : fetch outlet info failed")
		return nil, nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Msg("AReport : outlet co-ordinates invalid")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    user.UserId + ": AGENT invalid outlet co-ordintaes - region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	totalEmp, err := impl.PsqlDb.FetchTotalNonManagerialEmployees(user.Region, user.Outlet, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AReport : fetch empCount failed")
		return nil, nil, err
	}

	allUserSales, err := impl.PsqlDb.FetchActiveUserSales(user.UserId, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AReport : fetch sales failed")
		return nil, nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("agentCount", totalEmp).Int("saleGroupCount", len(allUserSales.SaleRecords)).
		Msg("AReport : total agents and salesGroup")

	resp := &models.TotalIncentive{}
	report := &models.UserReportDB{
		UserId:          user.UserId,
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
		return resp, report, nil
	}

	if totalEmp == 0 {
		ind := &models.IndividualIncentive{Error: "NO AGENT FOUND FOR OUTLET - " + user.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
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
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Err(internalErr).Msg("AReport : fetch target failed")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Message:    internalErr.Error(),
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
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : HBB fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : HBB fetch maxExtraPay failed")
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
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Msg("AReport : HBB sales info")

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

			ind.Payout = math.Round(ind.Payout*100) / 100
			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}

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
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : HAYYAK fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : HAYYAK fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
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
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Msg("AReport : HAYYAK sales info")

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

			ind.Payout = math.Round(ind.Payout*100) / 100
			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}

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
			maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : BAQATI fetch maxPay failed")
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}
			maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.Promoter, grp)
			if internalErr != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", grp).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : BAQATI fetch maxExtraPay failed")
				ind.Error = internalErr.Error()
				resp.HBB = ind
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
				Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Msg("AReport : BAQATI sales info")

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
			ind.Payout = math.Round(ind.Payout*100) / 100
			if ind.Payout > maxPayAllowed {
				ind.Payout = maxPayAllowed
			}

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
		return resp, report, nil
	}
	extraPayout, internalErr := impl.getExtraPay(quarter, year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Float64("achieved", achieved).Err(internalErr).Msg("AReport : UPGRADE fetch extraPayout failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		return resp, report, nil
	}
	maxPayAllowed, internalErr := impl.getMaxPayout(quarter, year, constant.Promoter, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : UPGRADE fetch maxPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		return resp, report, nil
	}
	maxExtraPayAllowed, internalErr := impl.getMaxExtraPayout(quarter, year, constant.Promoter, constant.UpgradeCategory)
	if internalErr != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("category", constant.UpgradeCategory).Str("position", constant.Promoter).Err(internalErr).Msg("AReport : UPGRADE fetch maxExtraPay failed")
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
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
		Float64("maxPayAllowed", maxPayAllowed).Float64("maxExtraPayAllowed", maxExtraPayAllowed).Msg("AReport : UPGRADE sales info")

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

	ind.Payout = math.Round(ind.Payout*100) / 100
	if ind.Payout > maxPayAllowed {
		ind.Payout = maxPayAllowed
	}

	totalPay += ind.Payout
	totalPay += ind.ExtraPayout

	resp.Upgrade = ind
	report.UpgradePayout = ind.Payout
	report.UpgradeExtraPayout = ind.ExtraPayout

	resp.TotalPay = math.Round(totalPay*100) / 100
	report.TotalPayout = resp.TotalPay

	go func(user *models.UserReportDB) {
		if user.HBBTargetCount <= 0 || user.HayyakTargetCount <= 0 || user.BaqatiTargetCount <= 0 || user.UpgradeTargetCount <= 0 {
			return
		}

		if err := impl.PsqlDb.InsertAgentUserReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("AReport : force report update")
		}
	}(report)

	return resp, report, nil
}

func (impl *IncentiveImpl) updateOutletManagerReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManager { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxTrace(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMReport : info")

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

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(omngr.Region, omngr.Outlet, quarter, year)
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

	userAgents, err := impl.PsqlDb.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(err).Msg("OMReport : fetch agents failed")
		return nil, err
	}

	report := &models.UserReportDB{
		UserId:          omngr.UserId,
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

		agentReport, err := impl.PsqlDb.FetchAgentUserReport(eachAgent.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMReport : cache failed")
				return nil, err
			}
		} else {
			if err != nil || (err == nil && agentReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				ind, resp, err := impl.updateAgentUserReport(ctx, eachAgent, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMReport : agent report failed")
					return nil, err
				}

				rawIncentiveReport, _ := json.Marshal(ind)
				rawUserReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("incentive", string(rawIncentiveReport)).
					Str("record", string(rawUserReport)).Msg("OMReport : agent log")
				agentReport = resp
			}
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
			Details:    fmt.Errorf("hbb target errr - %v", tarErr),
		}
	}

	achieved := float64(report.HBBSalesCount) * 100 / float64(target)
	vop, _ := impl.getValueOfPoint(quarter, year, constant.HBBCategory, achieved)
	maxPayAllowed, _ := impl.getMaxPayout(quarter, year, constant.OutletManager, constant.HBBCategory)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.HBBSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("OMReport : HBB sales info")

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
			Details:    fmt.Errorf("hayyak target errr - %v", tarErr),
		}
	}

	achieved = float64(report.HayyakSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.HayyakCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManager, constant.HayyakCategory)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.HayyakSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("OMReport : HAYYAK sales info")

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
			Details:    fmt.Errorf("baqati target errr - %v", tarErr),
		}
	}

	achieved = float64(report.BaqatiSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.BaqatiCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManager, constant.BaqatiCategory)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.BaqatiSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("OMReport : BAQATI sales info")

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
			Details:    fmt.Errorf("upgrade target errr - %v", tarErr),
		}
	}

	achieved = float64(report.UpgradeSalesCount) * 100 / float64(target)
	vop, _ = impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, achieved)
	maxPayAllowed, _ = impl.getMaxPayout(quarter, year, constant.OutletManager, constant.UpgradeCategory)

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("target", target).Int("achievedCount", report.UpgradeSalesCount).Float64("achievedPercent", achieved).
		Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("OMReport : UPGRADE sales info")

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

	go func(user *models.UserReportDB) {

		if err := impl.PsqlDb.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("OMReport : force report update")
		}
	}(report)

	return report, nil
}

func (impl *IncentiveImpl) updateManagerReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "invalid mngr details",
		}
	}

	if _, ok := constant.ManagerMap[mngr.Position]; !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}

	logger.FileLogger.CtxTrace(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManager { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerReport(ctx, mngr, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresident && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
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

	report := &models.UserReportDB{
		UserId:          mngr.UserId,
		Position:        mngr.Position,
		ManagerId:       mngr.ManagerId,
		ManagerPosition: mngr.ManagerPosition,
		Region:          mngr.Region,
		Outlet:          mngr.Outlet,
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(mngr.Region, mngr.Outlet, quarter, year)
	if err == nil {
		report.Latitude = outletResponse.Lat
		report.Longitude = outletResponse.Lang
	}

	allDistinctOM, err := impl.PsqlDb.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Err(err).Msg("MngrReport : fetch om failed")
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

		omReport, err := impl.PsqlDb.FetchManagerReport(om.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", om.UserId).Err(err).Msg("MngrReport : cache failed")
				return nil, err
			}
		} else {
			if err != nil || (err == nil && omReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				resp, err := impl.updateOutletManagerReport(ctx, om, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("om", om.UserId).Err(err).Msg("MngrReport : om report failed")
					return nil, err
				}

				rawUserReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", om.UserId).Str("record", string(rawUserReport)).Msg("OMReport : agent log")
				omReport = resp
			}
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

	if mngr.Position == constant.Manager || mngr.Position == constant.SeniorManager {
		//HBB
		vop, _ := impl.getValueOfPoint(quarter, year, constant.HBBCategory, report.HBBAchievedPercentage)
		maxPayAllowed, _ := impl.getMaxPayout(quarter, year, mngr.Position, constant.HBBCategory)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HBBTargetCount).Int("achievedCount", report.HBBSalesCount).
			Float64("achievedPercent", report.HBBAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("MngrReport : HBB sales info")
		for plan, qty := range hbbSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HBBPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HBBPayout).Msg("MngrReport : HBB per plan info")
		}
		if report.HBBPayout > maxPayAllowed {
			report.HBBPayout = maxPayAllowed
		}

		//HAYYAK
		vop, _ = impl.getValueOfPoint(quarter, year, constant.HayyakCategory, report.HayyakAchievedPercentage)
		maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.HayyakCategory)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.HayyakTargetCount).Int("achievedCount", report.HayyakSalesCount).
			Float64("achievedPercent", report.HayyakAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("MngrReport : HAYYAK sales info")
		for plan, qty := range hayyakSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.HayyakPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.HayyakPayout).Msg("MngrReport : HAYYAK per plan info")
		}
		if report.HayyakPayout > maxPayAllowed {
			report.HayyakPayout = maxPayAllowed
		}

		//BAQATI
		vop, _ = impl.getValueOfPoint(quarter, year, constant.BaqatiCategory, report.BaqatiAchievedPercentage)
		maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.BaqatiCategory)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.BaqatiTargetCount).Int("achievedCount", report.BaqatiSalesCount).
			Float64("achievedPercent", report.BaqatiAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("MngrReport : BAQATI sales info")
		for plan, qty := range baqatiSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.BaqatiPayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.BaqatiPayout).Msg("MngrReport : BAQATI per plan info")
		}
		if report.BaqatiPayout > maxPayAllowed {
			report.BaqatiPayout = maxPayAllowed
		}

		//UPGRADE
		vop, _ = impl.getValueOfPoint(quarter, year, constant.UpgradeCategory, report.UpradeAchievedPercentage)
		maxPayAllowed, _ = impl.getMaxPayout(quarter, year, mngr.Position, constant.UpgradeCategory)
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Float64("target", report.UpgradeTargetCount).Int("achievedCount", report.UpgradeSalesCount).
			Float64("achievedPercent", report.UpradeAchievedPercentage).Float64("vop", vop).Float64("maxPayAllowed", maxPayAllowed).Msg("MngrReport : UPGRADE sales info")
		for plan, qty := range upgradeSaleMap {
			rate, _ := impl.getSIPRatePlan(quarter, year, plan)
			report.UpgradePayout += float64(qty) * vop * rate
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("planName", plan).Float64("planRate", rate).
				Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Float64("cumulativePay", report.UpgradePayout).Msg("MngrReport : UPGRADE per plan info")
		}
		if report.UpgradePayout > maxPayAllowed {
			report.UpgradePayout = maxPayAllowed
		}

		report.TotalPayout = report.HBBPayout + report.HayyakPayout + report.BaqatiPayout + report.UpgradePayout
	}

	go func(user *models.UserReportDB) {
		if err := impl.PsqlDb.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("MngrReport : force report update")
		}
	}(report)

	return report, nil
}

func (impl *IncentiveImpl) updateOutletReport(ctx *gin.Context, regionName, outletName string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
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

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(regionName, outletName, quarter, year)
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

	hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.PsqlDb.FetchOutletSalesCount(regionName, outletName, isQuarterReport, quarter, month, year)
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
		if err := impl.PsqlDb.InsertOutletReport(report, isQuarterReport, quarter, month, year); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("outlet", outletName).Err(errors.New(err.Details.(string))).Msg("updateOutletReport: force report update")
		}
	}(outletReport)

	return outletReport, nil
}

func (impl *IncentiveImpl) updateRegionReport(ctx *gin.Context, regionName string, isQuarterReport bool, quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
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

	allOutlets, err := impl.PsqlDb.FetchAllOutletNames(regionName, quarter, year)
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

		outletReport, err := impl.PsqlDb.FetchOutletReport(eachOutlet, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && outletReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				outletResp, err := impl.updateOutletReport(ctx, regionName, eachOutlet, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Interface("outletReport", outletResp).Msg("updateRegionReport : trace")
				//TODO: test this out for struct assignment
				outletReport = outletResp
			}
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
		if err := impl.PsqlDb.InsertRegionReport(report, isQuarterReport, quarter, month, year); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("region", regionName).Err(errors.New(err.Details.(string))).Msg("updateRegionReport: force report update")
		}
	}(regionReport)

	return regionReport, nil
}
