package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	"math"
	"time"
)

func (impl *IncentiveImpl) GetSingleUserIncentive(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) (*models.TotalIncentive, *errors.Error) {

	trx, err := impl.PsqlDb.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	resp := &models.TotalIncentive{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &models.IndividualIncentive{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp, err := impl.PsqlDb.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		ind := &models.IndividualIncentive{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp--

	totalPay := 0.0
	// resp.HBB.Error = "nil"
	// resp.Hayyak.Error = "nil"
	// resp.Baqati.Error = "nil"
	// resp.Upgrade.Error = "nil"

	grpPlans := make(map[string]map[string]int)
	upgradedPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)

	totalActiveUpgrade := 0

	for grp, rcds := range trx.SaleRecords {

		//only 3 groups for now Hayyak, Hbb, Baqati
		logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Str("category", grp).Int("saleCount", len(rcds)).Msg("GetSingleUserIncentive : groupSalesCount info")

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

		targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)
		//if targetPerOutlet == -1 {
		//	return nil, errors.New("outlet target not found")
		//}

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
				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HBB sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : HBB sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HBB plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.HBB = ind

		case constant.HayyakCategory:

			ind := &models.IndividualIncentive{}

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : HAYYAK sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : HAYYAK sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HAYYAK plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.Hayyak = ind

		case constant.BaqatiCategory:

			ind := &models.IndividualIncentive{}

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : BAQATI sales error")
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
					Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
					Int("achievedCount", totalActive[grp]).Float64("achievedPercent", achieved).Float64("vop", vop).
					Float64("extraPay", impl.getExtraPay(grp, achieved)).Msg("GetSingleUserIncentive : BAQATI sales info")

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("planName", plan).Float64("planRate", rate).
						Int("qty", qty).Float64("pay", float64(qty)*vop*rate).Msg("GetSingleUserIncentive : HAYYAK plans info")
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
				ind.Payout = math.Round(ind.Payout*100) / 100
				totalPay += ind.Payout
			}
			resp.Baqati = ind
		}
	}

	ind := &models.IndividualIncentive{}

	targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, constant.UpgradeCategory, trx.Region, trx.Outlet)
	if targetPerOutlet == -1 {
		ind.Error = "outlet target not found"
		resp.Upgrade = ind
		return resp, nil
	}

	target := float64(targetPerOutlet) / float64(totalEmp)
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
		logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Msg("GetSingleUserIncentive : UPGRADE sales error")
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(constant.UpgradeCategory, achieved)

		logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
			Int("totalEmp", totalEmp).Int("targetPerOutlet", targetPerOutlet).Float64("indTarget", target).
			Int("achievedCount", totalActiveUpgrade).Float64("achievedPercent", achieved).Float64("vop", vop).
			Float64("extraPay", impl.getExtraPay(constant.UpgradeCategory, achieved)).Msg("GetSingleUserIncentive : UPGRADE sales info")

		ind.ExtraPayout = impl.getExtraPay(constant.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		totalPay += ind.ExtraPayout

		for currPlan, prevPlansMap := range upgradedPlans {
			for prevPln, qty := range prevPlansMap {
				rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
				logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("currPlanName", currPlan).Str("prevPlanName", prevPln).
					Float64("currPlanRate", impl.getSIPRatePlan(currPlan)).Float64("prevPlanRate", impl.getSIPRatePlan(prevPln)).
					Float64("rate", rate).Int("qty", qty).Float64("pay", float64(qty)*vop*rate).
					Msg("GetSingleUserIncentive : UPGRADE plans info")
				ind.Payout += float64(qty) * vop * rate
			}
		}

		max := impl.getMaxPayout(constant.UpgradeCategory)
		if ind.Payout > max {
			ind.Payout = max
		}
		ind.Payout = math.Round(ind.Payout*100) / 100
		totalPay += ind.Payout
	}

	resp.Upgrade = ind
	resp.TotalPay = math.Round(totalPay*100) / 100

	return resp, nil
}
