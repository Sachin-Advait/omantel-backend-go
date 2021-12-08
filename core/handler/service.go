package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
	//restapi "gitlab.com/rendezvousPoint/incentive_calculator/core/restApi"
)

type DBService interface {
	FetchTotalNonManagerialEmployees(region, outlet string) (int, error)
	FetchUserSales(clientID, month, year, quarter int, isQuarterReport bool) (*db.SaleTransaction, error)
	FetchEmployeeContactNumber(clientID int) (int, *commons.GenericError)
	InsertStaffOTP(otpModel *db.OtpModel) *commons.GenericError
	FetchOTPValidity(clientID int) (*db.OtpModel, *commons.GenericError)
	FetchAllGivenPositionEmpIdAndTheirManager(position string) ([]*db.UserLead, *commons.GenericError)
	InsertEmployeeRecord(usr *db.UserRecord) *commons.GenericError
	FetchAllEmployeesReportingToGivenManagerRole(managerPosition string, managerID int) ([]int, *commons.GenericError)
	FetchIndividualUserRecord(id int) ([]*db.UserRecord, *commons.GenericError)
}

type SalesIncentiveImpl struct {
	DBService DBService
}

func NewSalesIncentiveService(conf *config.Database) *SalesIncentiveImpl {
	dbConn := db.NewDatabaseSetup(conf)

	return &SalesIncentiveImpl{DBService: dbConn}
}

func (impl *SalesIncentiveImpl) FetchPromoterUserIncentive(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.TotalIncentive, error) {

	trx, err := impl.DBService.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	fmt.Println("data", clientID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet)

	resp := &db.TotalIncentive{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &db.IndividualIncentive{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		ind := &db.IndividualIncentive{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	fmt.Println("TotalEmp", totalEmp)

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

		targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)
		//if targetPerOutlet == -1 {
		//	return nil, errors.New("outlet target not found")
		//}

		switch grp {

		case commons.HBBCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HBB sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				break
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HBB: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.HBB = ind

		case commons.HayyakCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HAYYAK sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HAYYAK: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.Hayyak = ind

		case commons.BaqatiCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("BAQATI sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("BAQATI: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.Baqati = ind
		}
	}

	ind := &db.IndividualIncentive{}

	targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.UpgradeCategory, trx.Region, trx.Outlet)
	if targetPerOutlet == -1 {
		ind.Error = "outlet target not found"
		resp.Upgrade = ind
		return resp, nil
	}

	fmt.Println("UPGRADE sales target in "+trx.Outlet, targetPerOutlet)

	target := float64(targetPerOutlet) / float64(totalEmp)
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)

		fmt.Println("UPGRADE: target, achieved, percent, vop", target, totalActiveUpgrade, achieved, vop)

		ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		totalPay += ind.ExtraPayout

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
		totalPay += ind.Payout
	}

	resp.Upgrade = ind
	resp.TotalPay = math.Round(totalPay*100) / 100

	return resp, nil
}

func (impl *SalesIncentiveImpl) FetchOutletManagerIncentive(outletManagerID int, t time.Time, quarter int, isQuarterReport bool) (*db.TotalIncentive, error) {

	trx, err := impl.DBService.FetchUserSales(outletManagerID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	fmt.Println("data", outletManagerID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet)

	resp := &db.TotalIncentive{}

	// if trx.Region == "" || trx.Outlet == "" {
	// 	return nil, errors.New("region and outlet not found")
	// }

	totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		return nil, errors.New("no employees found")
	}

	totalEmp--

	totalPay := 0.0

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

		targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)
		//if targetPerOutlet == -1 {
		//	return nil, errors.New("outlet target not found")
		//}

		switch grp {

		case commons.HBBCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HBB sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				break
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HBB: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.HBB = ind

		case commons.HayyakCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("HAYYAK sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("HAYYAK: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.Hayyak = ind

		case commons.BaqatiCategory:

			ind := &db.IndividualIncentive{}

			fmt.Println("BAQATI sales target in "+trx.Outlet, targetPerOutlet)

			target := float64(targetPerOutlet)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Println("BAQATI: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

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
				totalPay += ind.Payout
			}
			resp.Baqati = ind
		}
	}

	ind := &db.IndividualIncentive{}

	targetPerOutlet := impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, commons.UpgradeCategory, trx.Region, trx.Outlet)
	if targetPerOutlet == -1 {
		ind.Error = "outlet target not found"
		resp.Upgrade = ind
		return resp, nil
	}

	fmt.Println("UPGRADE sales target in "+trx.Outlet, targetPerOutlet)

	target := float64(targetPerOutlet)
	if target <= 0 {
		ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
	} else {
		achieved := float64(totalActiveUpgrade) * 100 / float64(target)
		vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)

		fmt.Println("UPGRADE: target, achieved, percent, vop", target, totalActiveUpgrade, achieved, vop)

		ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		totalPay += ind.ExtraPayout

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
		totalPay += ind.Payout
	}

	resp.Upgrade = ind
	resp.TotalPay = math.Round(totalPay*100) / 100

	return resp, nil
}

func (impl *SalesIncentiveImpl) FetchPromoterTargetChart(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.ChartDataResponse, error) {

	trx, err := impl.DBService.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	fmt.Println("data", clientID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet)

	resp := &db.ChartDataResponse{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &db.TargetAchievementData{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		ind := &db.TargetAchievementData{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	for grp, rcds := range trx.SaleRecords {

		data := &db.TargetAchievementData{}

		for _, v := range rcds {
			data.SalesCount += v.Quantity
		}

		target := float64(impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)) / float64(totalEmp)
		if target <= 0 {
			data.Error = fmt.Errorf("per employee target found is %f", target).Error()
		}

		data.Target = math.Round(target*100) / 100
		data.AchievedPercentage = math.Round((float64(data.SalesCount)*100/target)*100) / 100

		switch grp {
		case commons.HBBCategory:
			resp.HBB = data
			return resp, nil
		case commons.HayyakCategory:
			resp.Hayyak = data
			return resp, nil
		case commons.BaqatiCategory:
			resp.Baqati = data
			return resp, nil
		case commons.UpgradeCategory:
			resp.Upgrade = data
			return resp, nil
		}
	}

	return resp, errors.New("no category found")
}

func (impl *SalesIncentiveImpl) FetchOutletManagerTargetChart(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.ChartDataResponse, error) {

	trx, err := impl.DBService.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	fmt.Println("data", clientID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet)

	resp := &db.ChartDataResponse{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &db.TargetAchievementData{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		ind := &db.TargetAchievementData{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	for grp, rcds := range trx.SaleRecords {

		data := &db.TargetAchievementData{}

		for _, v := range rcds {
			data.SalesCount += v.Quantity
		}

		target := float64(impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet))
		if target <= 0 {
			data.Error = fmt.Errorf("per employee target found is %f", target).Error()
		}

		data.Target = math.Round(target*100) / 100
		data.AchievedPercentage = math.Round((float64(data.SalesCount)*100/target)*100) / 100

		switch grp {
		case commons.HBBCategory:
			resp.HBB = data
			return resp, nil
		case commons.HayyakCategory:
			resp.Hayyak = data
			return resp, nil
		case commons.BaqatiCategory:
			resp.Baqati = data
			return resp, nil
		case commons.UpgradeCategory:
			resp.Upgrade = data
			return resp, nil
		}
	}

	return resp, errors.New("no category found")
}

func (impl *SalesIncentiveImpl) FetchSaleCountPerPlan(clientID int, group string, t time.Time, quarter int, isQuarterReport bool) ([]*db.SaleCountPerPlan, error) {

	trx, err := impl.DBService.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	fmt.Println("FetchSaleCountPerPlan", clientID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet, len(trx.SaleRecords))

	// type temp struct {
	// 	previousPlan string
	// 	qty          int
	// }

	// m := make(map[string]*temp)

	// for grp, rcds := range trx.SaleRecords {

	// 	if grp != group {
	// 		continue
	// 	}

	// 	for _, each := range rcds {
	// 		v, ok := m[each.CurrentProduct]
	// 		if !ok || v.previousPlan != each.PreviousProduct {
	// 			m[each.CurrentProduct] = &temp{
	// 				previousPlan: each.PreviousProduct,
	// 				qty:          each.Quantity,
	// 			}
	// 		} else {
	// 			v.qty += each.Quantity
	// 			m[each.CurrentProduct] = v
	// 		}
	// 	}
	// }

	// resp := make([]*db.SaleCountPerPlan, 0)

	// for k, v := range m {
	// 	resp = append(resp, &db.SaleCountPerPlan{
	// 		CurrentPlan:  k,
	// 		PreviousPlan: v.previousPlan,
	// 		Quantity:     v.qty,
	// 	})
	// }

	// sort.Slice(resp, func(i, j int) bool {
	// 	return resp[i].Quantity < resp[j].Quantity
	// })

	m := make(map[string]int)

	for grp, rcds := range trx.SaleRecords {

		if group != commons.UpgradeCategory && grp != group {
			continue
		}

		for _, each := range rcds {

			v, ok := m[each.CurrentProduct]
			if !ok {
				m[each.CurrentProduct] = each.Quantity
			}

			if grp == commons.UpgradeCategory && each.PreviousProduct != "" {
				m[each.CurrentProduct] = v + each.Quantity
			}

			if grp != commons.UpgradeCategory && each.PreviousProduct == "" {
				m[each.CurrentProduct] = v + each.Quantity
			}
		}
	}

	resp := make([]*db.SaleCountPerPlan, 0)

	for k, v := range m {
		resp = append(resp, &db.SaleCountPerPlan{
			CurrentPlan: k,
			Quantity:    v,
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Quantity > resp[j].Quantity
	})

	return resp, nil
}

func (impl *SalesIncentiveImpl) GenerateOTP(clientID int, requestID string) *commons.GenericError {

	otpModel, err := impl.DBService.FetchOTPValidity(clientID)
	if err != nil {
		return err
	}

	contactNumber, err := impl.DBService.FetchEmployeeContactNumber(clientID)
	if err != nil {
		return err
	}

	fmt.Println(contactNumber)

	if !otpModel.ExpiryTime.After(time.Now().Local()) {

		otpString, err := impl.generateOTP(6)
		if err != nil {
			return &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}

		otpString = "123456"

		otpModel.OTP, _ = strconv.Atoi(otpString)
		otpModel.Retry = 3
		otpModel.StaffID = clientID
		otpModel.CreatedTime = time.Now().Local()
		otpModel.ExpiryTime = otpModel.CreatedTime.Add(5 * time.Minute)

		//store OTP in db
		if err := impl.DBService.InsertStaffOTP(otpModel); err != nil {
			return err
		}
	}

	//reqSMS := &restapi.SMSParams{
	//	HeaderRequestID:             requestID,
	//	HeaderStaffID:               strconv.Itoa(clientID),
	//	HeaderTimeStamp:             otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	//	NotificationDateTime:        otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	//	NotificaitionValidityPeriod: otpModel.CreatedTime.Add(3 * time.Minute).Format("2006-01-02 15:04:05"),
	//	Message:                     "Your OTP is: " + strconv.Itoa(otpModel.OTP),
	//	ContactNumber:               strconv.Itoa(contactNumber),
	//}
	//
	//if err = restapi.SendOTP(reqSMS); err != nil {
	//	return err
	//}

	return nil
}

func (impl *SalesIncentiveImpl) VefiryOTP(clientID, otp int) (bool, *commons.GenericError) {

	otpModel, err := impl.DBService.FetchOTPValidity(clientID)
	if err != nil {
		return false, err
	}

	if otp != otpModel.OTP {
		return false, &commons.GenericError{Code: 1022, ErrorMessage: "otp invalid"}
	}

	if !otpModel.ExpiryTime.After(time.Now().Local()) {
		return false, &commons.GenericError{Code: 1021, ErrorMessage: "otp has expired"}
	}

	//retry count

	return true, nil
}

func (impl *SalesIncentiveImpl) generateOTP(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", err
	}

	if n != max {
		return "", errors.New("incorrect otp length")
	}

	for i := 0; i < len(b); i++ {
		b[i] = commons.Table[int(b[i])%len(commons.Table)]
	}

	return string(b), nil
}
