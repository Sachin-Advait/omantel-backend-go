package business

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
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

func (impl *IncentiveImpl) getValueOfPoint(grp string, achievedPercent float64) float64 {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent < 75:
		return models.ValOfPoint.VopMap[grp][75].Point
	case achievedPercent < 96:
		return models.ValOfPoint.VopMap[grp][95].Point
	case achievedPercent <= 111:
		return models.ValOfPoint.VopMap[grp][110].Point
	case achievedPercent < 140:
		return models.ValOfPoint.VopMap[grp][139].Point
	case achievedPercent >= 140:
		return models.ValOfPoint.VopMap[grp][140].Point
	}

	return 0
}

func (impl *IncentiveImpl) getExtraPay(grp string, achievedPercent float64) float64 {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return models.ValOfPoint.VopMap[grp][159].ExtraPay
	case achievedPercent >= 160 && achievedPercent < 180:
		return models.ValOfPoint.VopMap[grp][179].ExtraPay
	case achievedPercent >= 180 && achievedPercent < 190:
		return models.ValOfPoint.VopMap[grp][189].ExtraPay
	case achievedPercent >= 190 && achievedPercent < 200:
		return models.ValOfPoint.VopMap[grp][199].ExtraPay
	case achievedPercent >= 200:
		return models.ValOfPoint.VopMap[grp][200].ExtraPay
	}

	return 0
}

func (impl *IncentiveImpl) getSIPRatePlan(planName string) float64 {

	planID := ""

	models.PlanRates.Mutex.RLock()
	defer models.PlanRates.Mutex.RUnlock()

	for id, v := range models.PlanRates.PlanRateMap {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != "" {
		return models.PlanRates.PlanRateMap[planID].SIPRatePlan
	}

	return 0
}

func (impl *IncentiveImpl) getSaleTargetPerOutlet(month, quarter int, isQuarterReport bool, grp, region, outlet string) int {

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
			return -1
		}
	}

	models.SalesTrgt.Mutex.RLock()
	defer models.SalesTrgt.Mutex.RUnlock()

	rMap, ok := models.SalesTrgt.SaleTargetMap[region]
	if !ok {
		return -1
	}

	if _, ok := rMap[outlet]; !ok {
		return -1
	}

	if isQuarterReport {
		//quarter target
		result := 0
		switch grp {
		case constant.HBBCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].HBB
			}
			return result
		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Hayyak
			}
			return result
		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Baqati
			}
			return result
		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Upgrade
			}
			return result
		}
	} else {
		//month target
		switch grp {
		case constant.HBBCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].HBB
		case constant.HayyakCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Hayyak
		case constant.BaqatiCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Baqati
		case constant.UpgradeCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][month-1].Upgrade
		}
	}

	return -1
}

func (impl *IncentiveImpl) getMaxPayout(grp string) float64 {
	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	return models.MaxPay.Payoutmap[grp]
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

func (impl *IncentiveImpl) generateOTP(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", err
	}

	if n != max {
		return "", errors.New("incorrect otp length")
	}

	for i := 0; i < len(b); i++ {
		b[i] = constant.Table[int(b[i])%len(constant.Table)]
	}

	return string(b), nil
}

func (impl *IncentiveImpl) isRegion(region string) bool {

	if region == "BATINAH NORTH" || region == "BATINAH SOUTH" || region == "DAKHLIA & DAHIRAH" || region == "DHOFAR" ||
		region == "MUSCAT NORTH" || region == "MUSCAT SOUTH" || region == "SHARQIYA" || region == "GRAND TOTAL" {
		return true
	}

	return false
}

func (impl *IncentiveImpl) isProduct(product string) bool {

	if product == constant.HBBCategory || product == constant.BaqatiCategory || product == constant.HayyakCategory || product == constant.UpgradeCategory {
		return true
	}

	if product == "UPGRADE (BAQATI + HBB)" {
		return true
	}

	return false
}

func (impl *IncentiveImpl) updateUserReport(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.TotalIncentive, *models.UserReportDB, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userId)
	if err != nil {
		return nil, nil, err //error - userId not in masterDB, outlet-region not in outletDB, outletManager not in hierarchyDB
	}

	trx, err := impl.PsqlDb.FetchUserSales(userId, month, year, quarter, isQuarterReport)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("total Sales", len(trx.SaleRecords))
	resp := &models.TotalIncentive{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &models.IndividualIncentive{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil, nil
	}

	totalEmp, err := impl.PsqlDb.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, nil, err
	} else if totalEmp == 0 {
		ind := &models.IndividualIncentive{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil, nil
	}

	totalEmp--

	totalPay := 0.0
	// resp.HBB.Error = "nil"
	// resp.Hayyak.Error = "nil"
	// resp.Baqati.Error = "nil"
	// resp.Upgrade.Error = "nil"

	report := &models.UserReportDB{
		UserId:      userId,
		ManagerId:   user.OutletManageId,
		Isquarterly: isQuarterReport,
		Month:       month,
		Quarter:     quarter,
		Year:        year,
		CreateAt:    time.Now(),
	}

	grpPlans := make(map[string]map[string]int)
	upgradedPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)

	totalActiveUpgrade := 0

	for grp, rcds := range trx.SaleRecords {

		//only 3 groups for now Hayyak, Hbb, Baqati
		logger.FileLogger.CtxInfo(ctx).Str("api", ctx.Request.URL.Path).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Str("category", grp).Int("saleCount", len(rcds)).Msg("group sales")

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

		targetPerOutlet := impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)
		if targetPerOutlet == -1 {
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    trx.Outlet + " target not found",
				Details:    trx.Outlet + " target not found",
			}
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

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HBB sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				report.Hbb_target = target
				report.Hbb_sales_count = totalActive[grp]
				report.Hbb_achieved_percentage = achieved

				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : HBB sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				saleCount := make([]*models.SaleCountPerPlan, 0)

				for plan, qty := range plans {
					saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HBB plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				sort.Slice(saleCount, func(i, j int) bool {
					return saleCount[i].Quantity > saleCount[j].Quantity
				})

				hBBSaleCountPerPlanArray, _ := json.Marshal(saleCount)
				report.HBBSaleCountPerPlanArray = string(hBBSaleCountPerPlanArray)

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.HBB = ind
			report.Hbb_payout = ind.Payout
			report.Hbb_extra_payout = ind.ExtraPayout

		case constant.HayyakCategory:

			ind := &models.IndividualIncentive{}

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HAYYAK sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				report.Hayyak_target = target
				report.Hayyak_sales_count = totalActive[grp]
				report.Hayyak_achieved_percentage = achieved

				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : HAYYAK sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				saleCount := make([]*models.SaleCountPerPlan, 0)

				for plan, qty := range plans {
					saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HAYYAK plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				sort.Slice(saleCount, func(i, j int) bool {
					return saleCount[i].Quantity > saleCount[j].Quantity
				})

				hayyakSaleCountPerPlanArray, _ := json.Marshal(saleCount)
				report.HayyakSaleCountPerPlanArray = string(hayyakSaleCountPerPlanArray)

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.Hayyak = ind
			report.Hayyak_payout = ind.Payout
			report.Hayyak_extra_payout = ind.ExtraPayout

		case constant.BaqatiCategory:

			ind := &models.IndividualIncentive{}

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : BAQATI sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				report.Baqati_target = target
				report.Baqati_sales_count = totalActive[grp]
				report.Baqati_achieved_percentage = achieved

				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : BAQATI sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				saleCount := make([]*models.SaleCountPerPlan, 0)

				for plan, qty := range plans {
					saleCount = append(saleCount, &models.SaleCountPerPlan{CurrentPlan: plan, Quantity: qty})
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HAYYAK plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				sort.Slice(saleCount, func(i, j int) bool {
					return saleCount[i].Quantity > saleCount[j].Quantity
				})

				baqatiSaleCountPerPlanArray, _ := json.Marshal(saleCount)
				report.BaqatiSaleCountPerPlanArray = string(baqatiSaleCountPerPlanArray)

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.Baqati = ind
			report.Baqati_payout = ind.Payout
			report.Baqati_extra_payout = ind.ExtraPayout
		}
	}

	ind := &models.IndividualIncentive{}

	targetPerOutlet := impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.UpgradeCategory, trx.Region, trx.Outlet)
	if targetPerOutlet == -1 {
		ind.Error = "outlet target not found"
		resp.Upgrade = ind
		return resp, nil, nil
	}

	target := float64(targetPerOutlet) / float64(totalEmp)
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
		logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : UPGRADE sales error")
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(constant.UpgradeCategory, achieved)

		report.Hayyak_target = target
		report.Hayyak_sales_count = totalActiveUpgrade
		report.Hayyak_achieved_percentage = achieved

		logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
			Int("achievedCount", totalActiveUpgrade).Float64("achievedPercent", achieved).Float64("vop", vop).
			Float64("extraPay", impl.getExtraPay(constant.UpgradeCategory, achieved)).Msg("GetSingleUserIncentive : UPGRADE sales info")

		ind.ExtraPayout = impl.getExtraPay(constant.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		totalPay += ind.ExtraPayout

		saleCount := make([]*models.SaleCountPerPlan, 0)

		for currPlan, prevPlansMap := range upgradedPlans {
			eachPlanSale := &models.SaleCountPerPlan{CurrentPlan: currPlan}
			for prevPln, qty := range prevPlansMap {
				eachPlanSale.Quantity += qty
				rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("currPlanName", currPlan).Str("prevPlanName", prevPln).
					Float64("currPlanRate", impl.getSIPRatePlan(currPlan)).Float64("prevPlanRate", impl.getSIPRatePlan(prevPln)).
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

		max := impl.getMaxPayout(constant.UpgradeCategory)
		if ind.Payout > max {
			ind.Payout = max
		}
		ind.Payout = math.Round(ind.Payout*100) / 100
		totalPay += ind.Payout
	}

	resp.Upgrade = ind
	report.Upgrade_payout = ind.Payout
	report.Upgrade_extra_payout = ind.ExtraPayout

	resp.TotalPay = math.Round(totalPay*100) / 100
	report.Total_payout = resp.TotalPay

	go func(user *models.UserReportDB) {
		if user.Hbb_target <= 0 || user.Hayyak_target <= 0 || user.Baqati_target <= 0 || user.Upgrade_target <= 0 {
			return
		}

		if err := impl.PsqlDb.InsertEmployeeReport(user, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", userId).Err(errors.New(err.Details.(string))).Msg("updateUserReport: force report update")
		}
	}(report)

	return resp, report, nil
}

func (impl *IncentiveImpl) updateManagerReport(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.ManagerReportDB, *errors.Error) {

	manager, err := impl.PsqlDb.FetchManagerBasicDetails(userId)
	if err != nil {
		return nil, err //error - userId not in masterDB, outlet-region not in outletDB, outletManager not in hierarchyDB
	}

	if manager.Position == constant.OutletManager {

		hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.PsqlDb.FetchOutletSalesCount(manager.Region, manager.Outlet, month, year, quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(manager.Region, manager.Outlet)
		if err != nil {
			return nil, err
		}

		oManagerReport := &models.ManagerReportDB{
			UserId:             manager.UserId,
			Position:           constant.EmployeePositionMap[constant.OutletManager],
			ManagerId:          manager.ManagerId,
			Region:             manager.Region,
			Outlet:             manager.Outlet,
			Latitude:           outletResponse.Lat,
			Longitude:          outletResponse.Lang,
			HBBSalesCount:      hbbSales,
			HayyakSalesCount:   hayyakSales,
			BaqatiSalesCount:   baqatiSales,
			UpgradeSalesCount:  upgradeSales,
			HBBTargetCount:     impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HBBCategory, manager.Region, manager.Outlet),
			HayyakTargetCount:  impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HayyakCategory, manager.Region, manager.Outlet),
			BaqatiTargetCount:  impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.BaqatiCategory, manager.Region, manager.Outlet),
			UpgradeTargetCount: impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.UpgradeCategory, manager.Region, manager.Outlet),
			CreatedAt:          time.Now(),
		}

		if err := impl.PsqlDb.InsertManagerReport(oManagerReport, isQuarterReport, month, year, quarter); err != nil {
			return nil, err
		}

		return oManagerReport, nil
	}

	// position higher than outlet manager

	mangerPosition := constant.ManagerMap[manager.Position]

	allDistinctOM, err := impl.PsqlDb.FetchAllOutletManagersReportingToGivenManagerRole(mangerPosition, userId)
	if err != nil {
		return nil, err
	}

	managerReport := &models.ManagerReportDB{
		UserId:    manager.UserId,
		Position:  constant.EmployeePositionMap[manager.Position],
		ManagerId: manager.ManagerId,
		Region:    manager.Region,
		Outlet:    manager.Outlet,
		CreatedAt: time.Now(),
	}

	for _, om := range allDistinctOM {

		hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.PsqlDb.FetchOutletSalesCount(om.Region, om.Outlet, month, year, quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		managerReport.HBBSalesCount += hbbSales
		managerReport.HayyakSalesCount += hayyakSales
		managerReport.BaqatiSalesCount += baqatiSales
		managerReport.UpgradeSalesCount += upgradeSales

		managerReport.HBBTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HBBCategory, om.Region, om.Outlet)
		managerReport.HayyakTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HayyakCategory, om.Region, om.Outlet)
		managerReport.BaqatiTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.BaqatiCategory, om.Region, om.Outlet)
		managerReport.UpgradeTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.UpgradeCategory, om.Region, om.Outlet)
	}

	go func(manager *models.ManagerReportDB) {
		if err := impl.PsqlDb.InsertManagerReport(manager, isQuarterReport, month, year, quarter); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", userId).Err(errors.New(err.Details.(string))).Msg("updateManagerReport: force report update")
		}
	}(managerReport)

	return managerReport, nil
}

func (impl *IncentiveImpl) updateOutletReport(ctx *gin.Context, regionName, outletName string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	outletResponse, err := impl.PsqlDb.FetchOutletDetailsViaName(regionName, outletName)
	if err != nil {
		return nil, err
	}

	hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.PsqlDb.FetchOutletSalesCount(regionName, outletName, month, year, quarter, isQuarterReport)
	if err != nil {
		return nil, err
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
		HBBTargetCount:     impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HBBCategory, regionName, outletName),
		HayyakTargetCount:  impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HayyakCategory, regionName, outletName),
		BaqatiTargetCount:  impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.BaqatiCategory, regionName, outletName),
		UpgradeTargetCount: impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.UpgradeCategory, regionName, outletName),
		CreatedAt:          time.Now(),
	}

	go func(report *models.OutletReportDB) {
		if err := impl.PsqlDb.InsertOutletReport(report, isQuarterReport, quarter, month, year); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("outlet", outletName).Err(errors.New(err.Details.(string))).Msg("updateOutletReport: force report update")
		}
	}(outletReport)

	return outletReport, nil
}

func (impl *IncentiveImpl) updateRegionReport(ctx *gin.Context, regionName string, isQuarterReport bool, quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	outlets, err := impl.PsqlDb.FetchAllOutletNames(regionName)
	if err != nil {
		return nil, err
	}

	regionReport := &models.RegionReportDB{
		Region:    regionName,
		CreatedAt: time.Now(),
	}

	for _, eachOutlet := range outlets {

		hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.PsqlDb.FetchOutletSalesCount(regionName, eachOutlet, month, year, quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		regionReport.HBBSalesCount += hbbSales
		regionReport.HayyakSalesCount += hayyakSales
		regionReport.BaqatiSalesCount += baqatiSales
		regionReport.UpgradeSalesCount += upgradeSales

		regionReport.HBBTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HBBCategory, regionName, eachOutlet)
		regionReport.HayyakTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.HayyakCategory, regionName, eachOutlet)
		regionReport.BaqatiTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.BaqatiCategory, regionName, eachOutlet)
		regionReport.UpgradeTargetCount += impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, constant.UpgradeCategory, regionName, eachOutlet)
	}

	go func(report *models.RegionReportDB) {
		if err := impl.PsqlDb.InsertRegionReport(report, isQuarterReport, quarter, month, year); err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("region", regionName).Err(errors.New(err.Details.(string))).Msg("updateRegionReport: force report update")
		}
	}(regionReport)

	return regionReport, nil
}
