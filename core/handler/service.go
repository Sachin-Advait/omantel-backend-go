package handler

import (
	"errors"
	"fmt"
	"math"

	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
)

type DBService interface {
	FetchTotalNonManagerialEmployees(region, outlet string) (int, error)
	FetchUserSales(clientID int) (*db.SaleTransaction, error)
}

type SalesIncentiveImpl struct {
	DBService DBService
}

func NewSalesIncentiveService(conf *config.Database) *SalesIncentiveImpl {
	dbConn := db.NewDatabaseSetup(conf)

	return &SalesIncentiveImpl{DBService: dbConn}
}

func (impl *SalesIncentiveImpl) FetchSingleUserIncentive(clientID int) (*db.TotalIncentive, error) {

	trx, err := impl.DBService.FetchUserSales(clientID)
	if err != nil {
		return nil, err
	}

	totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		return nil, errors.New("no employee found")
	}

	fmt.Println("TotalEmp", totalEmp)

	totalEmp--

	resp := &db.TotalIncentive{}
	// resp.HBB.Error = "nil"
	// resp.Hayyak.Error = "nil"
	// resp.Baqati.Error = "nil"
	// resp.Upgrade.Error = "nil"

	grpPlans := make(map[string]map[string]int)
	upgradedPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)

	totalActiveUpgrade := 0

	for grp, rcds := range trx.SaleRecords {

		fmt.Println("Total Sales in "+grp, len(rcds))

		for _, v := range rcds {
			if v.Category != commons.UpgradeCategory {
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

		switch grp {

		case commons.HBBCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HBB sales target in"+trx.Outlet, impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet))

			target := float64(impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet)) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				break
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HBB: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				resp.TotalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Println("HBB planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				resp.TotalPay += ind.Payout
			}
			resp.HBB = ind

		case commons.HayyakCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HAYYAK sales target in"+trx.Outlet, impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet))

			target := float64(impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet)) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HAYYAK: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				resp.TotalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Println("HAYYAK planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}

				ind.Payout = math.Round(ind.Payout*100) / 100
				resp.TotalPay += ind.Payout
			}
			resp.Hayyak = ind

		case commons.BaqatiCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("BAQATI sales target in"+trx.Outlet, impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet))

			target := float64(impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet)) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("BAQATI: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				resp.TotalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Println("BAQATI planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
				ind.Payout = math.Round(ind.Payout*100) / 100
				resp.TotalPay += ind.Payout
			}
			resp.Baqati = ind
		}
	}

	ind := &db.IndividualIncentive{}

	fmt.Println("UPGRADE sales target in"+trx.Outlet, impl.getSaleTargetPerOutlet(commons.UpgradeCategory, trx.Region, trx.Outlet))

	target := float64(impl.getSaleTargetPerOutlet(commons.UpgradeCategory, trx.Region, trx.Outlet)) / float64(totalEmp)
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)

		fmt.Println("UPGRADE: target, achieved, percent, vop", target, totalActiveUpgrade, achieved, vop)

		ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		resp.TotalPay += ind.ExtraPayout

		for currPlan, prevPlansMap := range upgradedPlans {
			for prevPln, qty := range prevPlansMap {
				rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
				fmt.Println("UPGRADE currPlanName prevPlanName  currPlanRate  prevPlanRate rate qty pay: "+currPlan+" "+prevPln, impl.getSIPRatePlan(currPlan), impl.getSIPRatePlan(prevPln), rate, qty, float64(qty)*vop*rate)
				ind.Payout += float64(qty) * vop * rate
			}
		}

		max := impl.getMaxPayout(commons.UpgradeCategory)
		if ind.Payout > max {
			ind.Payout = max
		}
		ind.Payout = math.Round(ind.Payout*100) / 100
		resp.TotalPay += ind.Payout
	}

	resp.Upgrade = ind

	return resp, nil
}
