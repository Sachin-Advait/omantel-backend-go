package handler

import (
	"errors"
	"fmt"

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

	totalEmp--

	fmt.Println("totaalEmp", totalEmp)

	resp := &db.TotalIncentive{}

	grpPlans := make(map[string]map[string]int)
	upgradedPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)

	totalActiveUpgrade := 0

	for grp, rcds := range trx.SaleRecords {

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

			target := impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet) / totalEmp
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %d for HBB", target).Error()
				break
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HBB", achieved, vop, totalActive[grp], target)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Println("rate", rate, qty)
					fmt.Println("pay", float64(qty)*vop*rate)
					ind.Payout += float64(qty) * vop * rate
				}

				fmt.Println("payout", ind.Payout)

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
			}
			resp.HBB = ind

		case commons.HayyakCategory:

			ind := &db.IndividualIncentive{}

			target := impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet) / totalEmp
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %d for HAYYAK", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
			}
			resp.Hayyak = ind

		case commons.BaqatiCategory:

			ind := &db.IndividualIncentive{}

			target := impl.getSaleTargetPerOutlet(grp, trx.Region, trx.Outlet) / totalEmp
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %d for BAQATI", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					ind.Payout += float64(qty) * vop * rate
				}

				max := impl.getMaxPayout(grp)
				if ind.Payout > max {
					ind.Payout = max
				}
			}
			resp.Baqati = ind
		}
	}

	ind := &db.IndividualIncentive{}
	target := impl.getSaleTargetPerOutlet(commons.UpgradeCategory, trx.Region, trx.Outlet) / totalEmp
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %d for BAQATI", target).Error()
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)
		ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)

		for currPlan, prevPlansMap := range upgradedPlans {
			for prevPln, qty := range prevPlansMap {
				rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
				ind.Payout += float64(qty) * vop * rate
			}
		}

		max := impl.getMaxPayout(commons.UpgradeCategory)
		if ind.Payout > max {
			ind.Payout = max
		}
	}

	resp.Upgrade = ind

	return resp, nil
}
