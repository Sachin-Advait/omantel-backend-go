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

func (impl *IncentiveImpl) getValueOfPoint(year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if year != models.ValOfPoint.Year {
		return 0, errors.New("VOP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[grp] == nil {
		return 0, errors.New("VOP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent < 75:
		return models.ValOfPoint.VopMap[grp][75].Point, nil
	case achievedPercent < 96:
		return models.ValOfPoint.VopMap[grp][95].Point, nil
	case achievedPercent <= 111:
		return models.ValOfPoint.VopMap[grp][110].Point, nil
	case achievedPercent < 140:
		return models.ValOfPoint.VopMap[grp][139].Point, nil
	case achievedPercent >= 140:
		return models.ValOfPoint.VopMap[grp][140].Point, nil
	}

	return 0, errors.New("VOP ACHIEVEMENT PERCENT INVALID")
}

func (impl *IncentiveImpl) getExtraPay(year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if year != models.ValOfPoint.Year {
		return 0, errors.New("VOEP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[grp] == nil {
		return 0, errors.New("VOEP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return models.ValOfPoint.VopMap[grp][159].ExtraPay, nil
	case achievedPercent >= 160 && achievedPercent < 180:
		return models.ValOfPoint.VopMap[grp][179].ExtraPay, nil
	case achievedPercent >= 180 && achievedPercent < 190:
		return models.ValOfPoint.VopMap[grp][189].ExtraPay, nil
	case achievedPercent >= 190 && achievedPercent < 200:
		return models.ValOfPoint.VopMap[grp][199].ExtraPay, nil
	case achievedPercent >= 200:
		return models.ValOfPoint.VopMap[grp][200].ExtraPay, nil
	}

	return 0, nil
}

func (impl *IncentiveImpl) getSIPRatePlan(year int, planName string) (float64, error) {

	planID := ""

	models.PlanRates.Mutex.RLock()
	defer models.PlanRates.Mutex.RUnlock()

	if year != models.PlanRates.Year {
		return 0, errors.New("PLAN-NAME YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	for id, v := range models.PlanRates.PlanRateMap {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != "" {
		return models.PlanRates.PlanRateMap[planID].SIPRatePlan, nil
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

	if year != models.SalesTrgt.Year {
		return 0, errors.New("TARGET YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	rMap, ok := models.SalesTrgt.SaleTargetMap[region]
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
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].HBB
			}
			return result, nil

		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Hayyak
			}
			return result, nil

		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Baqati
			}
			return result, nil

		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Upgrade
			}
			return result, nil

		}
	} else {
		//month target
		switch grp {
		case constant.HBBCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].HBB, nil
		case constant.HayyakCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Hayyak, nil
		case constant.BaqatiCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Baqati, nil
		case constant.UpgradeCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Upgrade, nil
		}
	}

	return 0, errors.New("TARGET GROUP INVALID - " + grp)
}

func (impl *IncentiveImpl) getMaxPayout(year int, grp string) (float64, error) {
	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	if year != models.SalesTrgt.Year {
		return 0, errors.New("MAX PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	val, ok := models.MaxPay.Payoutmap[grp]
	if !ok {
		return 0, errors.New("MAX PAYOUT GROUP INVALID - " + grp)
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

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "invalid agent details",
			}
		}
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "invalid agent details - userId:" + user.UserId,
		}
	}

	if !(user.Position == constant.Promoter || user.Position == constant.CSR) {
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		return nil, nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    user.UserId + ": AGENT invalid outlet co-ordintaes - region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.PsqlDb.FetchActiveUserSales(user.UserId, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, nil, err
	}

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

	totalEmp, err := impl.PsqlDb.FetchTotalNonManagerialEmployees(user.Region, user.Outlet, quarter, year)
	if err != nil {
		return nil, nil, err
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

		//only 3 groups for now Hayyak, Hbb, Baqati
		logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Str("category", grp).Int("grpSaleCount", len(rcds)).Msg("group sales")

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

	for grp, plans := range grpPlans {

		targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, grp, user.Region, user.Outlet)
		if internalErr != nil {
			return nil, nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Message:    internalErr.Error(),
				Details:    internalErr.Error(),
			}
		}

		if targetPerOutlet == 0 {
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
			if targetPerEmp <= 0 {
				ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE HBB TARGET-%v", user, targetPerEmp)
				resp.HBB = ind
				logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HBB sales error")
				break
			}

			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)

			vop, internalErr := impl.getValueOfPoint(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}

			extraPayout, internalErr := impl.getExtraPay(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}

			maxPayAllowed, internalErr := impl.getMaxPayout(year, grp)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.HBB = ind
				break
			}

			report.HBBTargetCount = targetPerEmp
			report.HBBSalesCount = totalActive[grp]
			report.HBBAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			totalPay += ind.ExtraPayout

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("outlet", user.Outlet).Str("region", user.Region).
				Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
				Float64("extraPay", extraPayout).Msg("GetSingleUserIncentive : HBB sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)

			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, internalErr := impl.getSIPRatePlan(year, plan)
				if internalErr != nil {
					ind.Error = internalErr.Error()
					resp.HBB = ind
					break
				}

				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HBB plans info")
				ind.Payout += float64(qty) * vop * rate
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

			ind.Payout = math.Round(ind.Payout*100) / 100
			totalPay += ind.Payout

			resp.HBB = ind
			report.HBBPayout = ind.Payout
			report.HBBExtraPayout = ind.ExtraPayout

		case constant.HayyakCategory:

			ind := &models.IndividualIncentive{}

			targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
			if targetPerEmp <= 0 {
				ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE HAYYAK TARGET-%v", user, targetPerEmp)
				resp.Hayyak = ind
				logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HAYYAK sales error")
				break
			}

			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)

			vop, internalErr := impl.getValueOfPoint(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}

			extraPayout, internalErr := impl.getExtraPay(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}

			maxPayAllowed, internalErr := impl.getMaxPayout(year, grp)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Hayyak = ind
				break
			}

			report.HayyakTargetCount = targetPerEmp
			report.HayyakSalesCount = totalActive[grp]
			report.HayyakAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			totalPay += ind.ExtraPayout

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("outlet", user.Outlet).Str("region", user.Region).
				Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
				Float64("extraPay", extraPayout).Msg("GetSingleUserIncentive : HAYYAK sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)

			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, internalErr := impl.getSIPRatePlan(year, plan)
				if internalErr != nil {
					ind.Error = internalErr.Error()
					resp.Hayyak = ind
					break
				}

				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HAYYAK plans info")
				ind.Payout += float64(qty) * vop * rate
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

			ind.Payout = math.Round(ind.Payout*100) / 100
			totalPay += ind.Payout

			resp.Hayyak = ind
			report.HayyakPayout = ind.Payout
			report.HayyakExtraPayout = ind.ExtraPayout

		case constant.BaqatiCategory:

			ind := &models.IndividualIncentive{}

			targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
			if targetPerEmp <= 0 {
				ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE BAQATI TARGET-%v", user, targetPerEmp)
				resp.Baqati = ind
				logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("BAQATI sales error")
				break
			}

			achieved := float64(totalActive[grp]) * 100 / float64(targetPerEmp)

			vop, internalErr := impl.getValueOfPoint(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}

			extraPayout, internalErr := impl.getExtraPay(year, grp, achieved)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}

			maxPayAllowed, internalErr := impl.getMaxPayout(year, grp)
			if internalErr != nil {
				ind.Error = internalErr.Error()
				resp.Baqati = ind
				break
			}

			report.BaqatiTargetCount = targetPerEmp
			report.BaqatiSalesCount = totalActive[grp]
			report.BaqatiAchievedPercentage = achieved
			ind.ExtraPayout = math.Round(extraPayout*100) / 100
			totalPay += ind.ExtraPayout

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("outlet", user.Outlet).Str("region", user.Region).
				Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
				Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
				Float64("extraPay", extraPayout).Msg("BAQATI sales info")

			saleCount := make([]*models.SaleCountPerPlan, 0)

			for plan, qty := range plans {
				saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
				rate, internalErr := impl.getSIPRatePlan(year, plan)
				if internalErr != nil {
					ind.Error = internalErr.Error()
					resp.Baqati = ind
					break
				}

				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("planName", plan).Float64("planRate", rate).
					Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("BAQATI plans info")
				ind.Payout += float64(qty) * vop * rate
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

			ind.Payout = math.Round(ind.Payout*100) / 100
			totalPay += ind.Payout

			resp.Baqati = ind
			report.BaqatiPayout = ind.Payout
			report.BaqatiExtraPayout = ind.ExtraPayout
		}
	}

	ind := &models.IndividualIncentive{}

	targetPerOutlet, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory, user.Region, user.Outlet)
	if internalErr != nil {
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	if targetPerOutlet == 0 {
		ind.Error = "NO TARGET FOUND FOR OUTLET - " + user.Outlet
		resp.Upgrade = ind
		return resp, report, nil
	}

	targetPerEmp := float64(targetPerOutlet) / float64(totalEmp)
	if targetPerEmp <= 0 {
		ind.Error = fmt.Sprintf("OUTLET-%v, PER EMPLOYEE UPGRADE TARGET-%v", user, targetPerEmp)
		resp.Upgrade = ind
		logger.FileLogger.CtxTrace(ctx).Str("clientID", user.UserId).Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("UPGRADE sales error")
		return resp, report, nil
	}

	achieved := float64(totalActiveUpgrade) * 100 / float64(targetPerEmp)

	vop, internalErr := impl.getValueOfPoint(year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		return resp, report, nil
	}

	extraPayout, internalErr := impl.getExtraPay(year, constant.UpgradeCategory, achieved)
	if internalErr != nil {
		ind.Error = internalErr.Error()
		resp.Upgrade = ind
		return resp, report, nil
	}

	maxPayAllowed, internalErr := impl.getMaxPayout(year, constant.UpgradeCategory)
	if internalErr != nil {
		ind.Error = internalErr.Error()
		resp.Baqati = ind
		return resp, report, nil
	}

	report.UpgradeTargetCount = targetPerEmp
	report.UpgradeSalesCount = totalActiveUpgrade
	report.UpradeAchievedPercentage = achieved
	ind.ExtraPayout = math.Round(extraPayout*100) / 100
	totalPay += ind.ExtraPayout

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("outlet", user.Outlet).Str("region", user.Region).
		Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("targetPerEmp", targetPerEmp).
		Int("achievedCount", totalActiveUpgrade).Float64("achievedPercent", achieved).Float64("vop", vop).
		Float64("extraPay", extraPayout).Msg("GetSingleUserIncentive : UPGRADE sales info")

	saleCount := make([]*models.SaleCountPerPlan, 0)

	for currPlan, prevPlansMap := range upgradedPlans {
		eachPlanSale := &models.SaleCountPerPlan{CurrentPlan: currPlan}
		for prevPln, qty := range prevPlansMap {
			eachPlanSale.Quantity += qty
			currPlanRate, internalErr := impl.getSIPRatePlan(year, currPlan)
			if internalErr != nil {
				ind.Error = "CURRENT:" + internalErr.Error()
				resp.Upgrade = ind
				return resp, report, nil
			}

			prevPlanRate, internalErr := impl.getSIPRatePlan(year, prevPln)
			if internalErr != nil {
				ind.Error = "PREVIOUS:" + internalErr.Error()
				resp.Upgrade = ind
				return resp, report, nil
			}

			rate := currPlanRate - prevPlanRate

			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("currPlanName", currPlan).Str("prevPlanName", prevPln).
				Float64("currPlanRate", currPlanRate).Float64("prevPlanRate", prevPlanRate).
				Float64("rate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
				Msg("GetSingleUserIncentive : UPGRADE plans info")
			ind.Payout += float64(qty) * vop * rate
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

	ind.Payout = math.Round(ind.Payout*100) / 100
	totalPay += ind.Payout

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
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("updateUserReport: force report update")
		}
	}(report)

	return resp, report, nil
}

func (impl *IncentiveImpl) updateOutletManagerReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {
		if omngr == nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "invalid omngr details",
			}
		}
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if !(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader) {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
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

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(omngr.Region, omngr.Outlet, quarter, year)
	if err != nil {
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": OUTLET MANAGER INVALID OUTLET CO-ORDINATES - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    omngr.UserId + ": OUTLET MANAGER invalid outlet co-ordintaes - region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
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

	userAgents, err := impl.PsqlDb.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(userAgents) == 0 {

		return report, nil

		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    "NO AGENTS FOUND FOR OM - " + omngr.UserId,
		// 	Details:    "no agents found for OM - " + omngr.UserId,
		// }
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	for _, eachAgent := range userAgents {

		agentReport, err := impl.PsqlDb.FetchAgentUserReport(eachAgent.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && agentReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				_, resp, err := impl.updateAgentUserReport(ctx, eachAgent, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				//TODO: test this out for struct assignment
				agentReport = resp
			}
		}

		report.HayyakSalesCount += agentReport.HayyakSalesCount
		report.HBBSalesCount += agentReport.HBBSalesCount
		report.BaqatiSalesCount += agentReport.BaqatiSalesCount
		report.UpgradeSalesCount += agentReport.UpgradeSalesCount

		report.HayyakTargetCount += agentReport.HayyakTargetCount
		report.HBBTargetCount += agentReport.HBBTargetCount
		report.BaqatiTargetCount += agentReport.BaqatiTargetCount
		report.UpgradeTargetCount += agentReport.UpgradeTargetCount

		report.HayyakAchievedPercentage += agentReport.HayyakAchievedPercentage
		report.HBBAchievedPercentage += agentReport.HBBAchievedPercentage
		report.BaqatiAchievedPercentage += agentReport.BaqatiAchievedPercentage
		report.UpradeAchievedPercentage += agentReport.UpradeAchievedPercentage

		report.HayyakPayout += agentReport.HayyakPayout
		report.HBBPayout += agentReport.HBBPayout
		report.BaqatiPayout += agentReport.BaqatiPayout
		report.UpgradePayout += agentReport.UpgradePayout

		report.HayyakExtraPayout += agentReport.HayyakExtraPayout
		report.HBBExtraPayout += agentReport.HBBExtraPayout
		report.BaqatiExtraPayout += agentReport.BaqatiExtraPayout
		report.UpgradeExtraPayout += agentReport.UpgradeExtraPayout

		report.TotalPayout += report.TotalPayout

		saleCount := make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(report.HBBSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hbbSaleMap[sale.CurrentPlan] += sale.Quantity
			}

			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.HayyakSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hayyakSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.BaqatiSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				baqatiSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.UpgradeSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				upgradeSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}
	}

	report.HayyakAchievedPercentage = report.HayyakAchievedPercentage / float64(len(userAgents))
	report.HBBAchievedPercentage = report.HBBAchievedPercentage / float64(len(userAgents))
	report.BaqatiAchievedPercentage = report.BaqatiAchievedPercentage / float64(len(userAgents))
	report.UpradeAchievedPercentage = report.UpradeAchievedPercentage / float64(len(userAgents))
	report.HayyakPayout = report.HayyakPayout / float64(len(userAgents))
	report.HBBPayout = report.HBBPayout / float64(len(userAgents))
	report.BaqatiPayout = report.BaqatiPayout / float64(len(userAgents))
	report.UpgradePayout = report.UpgradePayout / float64(len(userAgents))
	report.HayyakExtraPayout = report.HayyakExtraPayout / float64(len(userAgents))
	report.HBBExtraPayout = report.HBBExtraPayout / float64(len(userAgents))
	report.BaqatiExtraPayout = report.BaqatiExtraPayout / float64(len(userAgents))
	report.UpgradeExtraPayout = report.UpgradeExtraPayout / float64(len(userAgents))
	report.TotalPayout += report.TotalPayout / float64(len(userAgents))

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
	saleArray = saleArray[:0]

	//hayyak sales count per plan
	for key, val := range hayyakSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	hayyakSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)
	saleArray = saleArray[:0]

	//baqati sales count per plan
	for key, val := range baqatiSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	baqatiSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)
	saleArray = saleArray[:0]

	//upgrade sales count per plan
	for key, val := range upgradeSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	upgradeSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.UpgradeSaleCountPerPlanArray = string(upgradeSaleCountPerPlanArray)

	go func(user *models.UserReportDB) {

		if err := impl.PsqlDb.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("updateOMReport: force report update")
		}
	}(report)

	return report, nil
}

func (impl *IncentiveImpl) updateManagerReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if mngr == nil || mngr.UserId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "invalid mngr details",
		}
	}

	if _, ok := constant.ManagerMap[mngr.Position]; !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManager || mngr.Position == constant.TeamLeader {
		report, err := impl.updateOutletManagerReport(ctx, mngr, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresident && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
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
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		return report, nil

		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    "NO OM FOUND FOR MANAGER - " + mngr.UserId,
		// 	Details:    "no om found for manager - " + mngr.UserId,
		// }
	}

	hbbSaleMap := make(map[string]int, 0)
	hayyakSaleMap := make(map[string]int, 0)
	baqatiSaleMap := make(map[string]int, 0)
	upgradeSaleMap := make(map[string]int, 0)

	for _, om := range allDistinctOM {

		omReport, err := impl.PsqlDb.FetchManagerReport(om.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && omReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				resp, err := impl.updateOutletManagerReport(ctx, om, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				//TODO: test this out for struct assignment
				omReport = resp
			}
		}

		report.HayyakSalesCount += omReport.HayyakSalesCount
		report.HBBSalesCount += omReport.HBBSalesCount
		report.BaqatiSalesCount += omReport.BaqatiSalesCount
		report.UpgradeSalesCount += omReport.UpgradeSalesCount

		report.HayyakTargetCount += omReport.HayyakTargetCount
		report.HBBTargetCount += omReport.HBBTargetCount
		report.BaqatiTargetCount += omReport.BaqatiTargetCount
		report.UpgradeTargetCount += omReport.UpgradeTargetCount

		report.HayyakAchievedPercentage += omReport.HayyakAchievedPercentage
		report.HBBAchievedPercentage += omReport.HBBAchievedPercentage
		report.BaqatiAchievedPercentage += omReport.BaqatiAchievedPercentage
		report.UpradeAchievedPercentage += omReport.UpradeAchievedPercentage

		report.HayyakPayout += omReport.HayyakPayout
		report.HBBPayout += omReport.HBBPayout
		report.BaqatiPayout += omReport.BaqatiPayout
		report.UpgradePayout += omReport.UpgradePayout

		report.HayyakExtraPayout += omReport.HayyakExtraPayout
		report.HBBExtraPayout += omReport.HBBExtraPayout
		report.BaqatiExtraPayout += omReport.BaqatiExtraPayout
		report.UpgradeExtraPayout += omReport.UpgradeExtraPayout

		report.TotalPayout += report.TotalPayout

		saleCount := make([]*models.SaleCountPerPlan, 0)
		if err := json.Unmarshal([]byte(report.HBBSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hbbSaleMap[sale.CurrentPlan] += sale.Quantity
			}

			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.HayyakSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				hayyakSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.BaqatiSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				baqatiSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}

		if err := json.Unmarshal([]byte(report.UpgradeSaleCountPerPlanArray), &saleCount); err == nil {
			for _, sale := range saleCount {
				upgradeSaleMap[sale.CurrentPlan] += sale.Quantity
			}
			saleCount = saleCount[:0]
		}
	}

	report.HayyakAchievedPercentage = report.HayyakAchievedPercentage / float64(len(allDistinctOM))
	report.HBBAchievedPercentage = report.HBBAchievedPercentage / float64(len(allDistinctOM))
	report.BaqatiAchievedPercentage = report.BaqatiAchievedPercentage / float64(len(allDistinctOM))
	report.UpradeAchievedPercentage = report.UpradeAchievedPercentage / float64(len(allDistinctOM))
	report.HayyakPayout = report.HayyakPayout / float64(len(allDistinctOM))
	report.HBBPayout = report.HBBPayout / float64(len(allDistinctOM))
	report.BaqatiPayout = report.BaqatiPayout / float64(len(allDistinctOM))
	report.UpgradePayout = report.UpgradePayout / float64(len(allDistinctOM))
	report.HayyakExtraPayout = report.HayyakExtraPayout / float64(len(allDistinctOM))
	report.HBBExtraPayout = report.HBBExtraPayout / float64(len(allDistinctOM))
	report.BaqatiExtraPayout = report.BaqatiExtraPayout / float64(len(allDistinctOM))
	report.UpgradeExtraPayout = report.UpgradeExtraPayout / float64(len(allDistinctOM))
	report.TotalPayout += report.TotalPayout / float64(len(allDistinctOM))

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
	saleArray = saleArray[:0]

	//hayyak sales count per plan
	for key, val := range hayyakSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	hayyakSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)
	saleArray = saleArray[:0]

	//baqati sales count per plan
	for key, val := range baqatiSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	baqatiSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)
	saleArray = saleArray[:0]

	//upgrade sales count per plan
	for key, val := range upgradeSaleMap {
		saleArray = append(saleArray, &models.SaleCountPerPlan{CurrentPlan: key, Quantity: val})
	}

	sort.Slice(saleArray, func(i, j int) bool {
		return saleArray[i].Quantity > saleArray[j].Quantity
	})

	upgradeSaleCountPerPlanArray, _ := json.Marshal(saleArray)
	report.UpgradeSaleCountPerPlanArray = string(upgradeSaleCountPerPlanArray)

	go func(user *models.UserReportDB) {
		if err := impl.PsqlDb.InsertManagerReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).Msg("updateMngrReport: force report update")
		}
	}(report)

	return report, nil
}

func (impl *IncentiveImpl) updateOutletReport(ctx *gin.Context, regionName, outletName string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	if regionName == "" || outletName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION OR OUTLET - region:" + regionName + " outlet:" + outletName,
			Details:    "blank region or outlet - region:" + regionName + " outlet:" + outletName,
		}
	}

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
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	hayyakTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HayyakCategory, regionName, outletName)
	if internalErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	baqatiTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.BaqatiCategory, regionName, outletName)
	if internalErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
	}

	upgradeTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory, regionName, outletName)
	if internalErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    internalErr.Error(),
			Details:    internalErr.Error(),
		}
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

	if regionName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION NAME",
			Details:    "blank region name",
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
