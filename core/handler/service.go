package handler

import (
	"crypto/rand"
	"encoding/json"
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
	FetchPromotersPresentYearSales(clientID int) (map[int]*db.SaleTransaction, map[int]*db.SaleTransaction, error)
	FetchEmployeeContactNumber(clientID int) (int, *commons.GenericError)
	InsertStaffOTP(otpModel *db.OtpModel) *commons.GenericError
	FetchOTPValidity(clientID int) (*db.OtpModel, *commons.GenericError)
	FetchAllGivenPositionEmpIdAndTheirManager(position string) ([]*db.UserLead, *commons.GenericError)
	InsertEmployeeRecord(usr *db.UserRecord) *commons.GenericError
	InsertEmployeeMultipleRecords(usr []*db.UserRecord) *commons.GenericError
	FetchAllEmployeesReportingToGivenManagerRole(managerPosition string, managerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchIndividualUserRecord(id int) ([]*db.UserRecord, *commons.GenericError)
	FetchIndividualUserPayoutReport(id, month, quarter int, isQuarterly bool) (*db.UserRecord, *commons.GenericError)
	FetchIndividualUserTargetReport(id, month, quarter int, isQuarterly bool) (*db.UserRecord, *commons.GenericError)
	FetchIndividualUserSaleCountPerPlanReport(id, month, quarter int, isQuarterly bool) (*db.UserRecord, *commons.GenericError)
	FetchIndividualUserOverAllReport(id int, month int, quarter int, isQuarterly bool) (*db.UserRecord, *commons.GenericError)
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

	fmt.Printf("fetchUserSales: id=%d, month=%d, year=%d quarter=%d isQuarterReport=%v, region=%s, outlet=%s\n", clientID, int(t.Month()), t.Year(), quarter, isQuarterReport, trx.Region, trx.Outlet)

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

	totalEmp--

	fmt.Printf("fetchTotalNonManagerialEmployees: total=%d\n", totalEmp)

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
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Printf("HBB: target=%v, achieved=%d, percent=%v, vop=%v, extraPay=%v\n", target, totalActive[grp], achieved, vop, impl.getExtraPay(grp, achieved))

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Printf("HBB: planName=%s, planRate=%v, qty=%v pay=%v\n",plan, rate, qty, float64(qty)*vop*rate)
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

				fmt.Printf("HAYYAK: target=%v, achieved=%d, percent=%v, vop=%v, extraPay=%v\n", target, totalActive[grp], achieved, vop, impl.getExtraPay(grp, achieved))

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Printf("HAYYAK: planName=%s, planRate=%v, qty=%v pay=%v\n",plan, rate, qty, float64(qty)*vop*rate)
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

			fmt.Println("BAQATI sales target in outlet:"+trx.Outlet+ " targetPerOutlet: ", targetPerOutlet)

			target := float64(targetPerOutlet) / float64(totalEmp)
			if target <= 0 {
				ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
			} else {
				achieved := float64(totalActive[grp]) * 100 / float64(target)
				vop := impl.getValueOfPoint(grp, achieved)

				fmt.Printf("BAQATI: target=%v, achieved=%d, percent=%v, vop=%v, extraPay=%v\n", target, totalActive[grp], achieved, vop, impl.getExtraPay(grp, achieved))

				ind.ExtraPayout = impl.getExtraPay(grp, achieved)
				ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
				totalPay += ind.ExtraPayout

				for plan, qty := range plans {
					rate := impl.getSIPRatePlan(plan)
					fmt.Printf("BAQATI: planName=%s, planRate=%v, qty=%v pay=%v\n",plan, rate, qty, float64(qty)*vop*rate)
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

		fmt.Printf("UPGRADE: target=%v, achieved=%d, percent=%v, vop=%v, extraPay=%v\n", target, totalActiveUpgrade, achieved, vop, impl.getExtraPay(commons.UpgradeCategory, achieved))

		ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
		ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
		totalPay += ind.ExtraPayout

		for currPlan, prevPlansMap := range upgradedPlans {
			for prevPln, qty := range prevPlansMap {
				rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
				fmt.Printf("UPGRADE: currPlanName=%s, prevPlanName=%s, currPlanRate=%s, prevPlanRate=%v, rate=%v, qty=%v pay=%v\n",currPlan, prevPln, impl.getSIPRatePlan(currPlan), impl.getSIPRatePlan(prevPln), rate, qty, float64(qty)*vop*rate)
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

	totalEmp--

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
		case commons.HayyakCategory:
			resp.Hayyak = data
		case commons.BaqatiCategory:
			resp.Baqati = data
		case commons.UpgradeCategory:
			resp.Upgrade = data
		default:
			return nil, errors.New("no category found")
		}
	}

	return resp, nil
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

func (impl *SalesIncentiveImpl) FetchUserOverAllDetails(clientID, month, quarter int, isQuarterReport bool) (*db.UserRecord, *commons.GenericError) {

	return impl.DBService.FetchIndividualUserOverAllReport(clientID, month, quarter, isQuarterReport)
}

func (impl *SalesIncentiveImpl) FetchUserReportingToOutletManager(outletManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError) {

	return impl.DBService.FetchAllEmployeesReportingToGivenManagerRole("outlet_manager", outletManagerID)
}

func (impl *SalesIncentiveImpl) FetchUserReportingToManager(managerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError) {

	return impl.DBService.FetchAllEmployeesReportingToGivenManagerRole("manager", managerID)
}

func (impl *SalesIncentiveImpl) FetchUserReportingToSeniorManager(seniortManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError) {

	return impl.DBService.FetchAllEmployeesReportingToGivenManagerRole("senior_manager", seniortManagerID)
}

func (impl *SalesIncentiveImpl) FetchUserReportingToGeneralManager(generaltManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError) {

	return impl.DBService.FetchAllEmployeesReportingToGivenManagerRole("general_manager", generaltManagerID)
}

func (impl *SalesIncentiveImpl) FetchUserReportingToVicePresident(vicePresidentID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError) {

	return impl.DBService.FetchAllEmployeesReportingToGivenManagerRole("vice_president", vicePresidentID)
}

func (impl *SalesIncentiveImpl) FetchUserPayoutReport(id int, t time.Time, quarter int, isQuarterly bool) (*db.TotalIncentive, *commons.GenericError) {

	usrRcrd, err := impl.DBService.FetchIndividualUserPayoutReport(id, int(t.Month()), quarter, isQuarterly)
	if err != nil {
		return nil, err
	}

	resp := &db.TotalIncentive{
		HBB: &db.IndividualIncentive{
			Payout:      usrRcrd.Hbb_payout,
			ExtraPayout: usrRcrd.Hbb_extra_payout,
		},
		Hayyak: &db.IndividualIncentive{
			Payout:      usrRcrd.Hayyak_payout,
			ExtraPayout: usrRcrd.Hayyak_extra_payout,
		},
		Baqati: &db.IndividualIncentive{
			Payout:      usrRcrd.Baqati_payout,
			ExtraPayout: usrRcrd.Baqati_extra_payout,
		},
		Upgrade: &db.IndividualIncentive{
			Payout:      usrRcrd.Upgrade_payout,
			ExtraPayout: usrRcrd.Upgrade_extra_payout,
		},
		TotalPay: usrRcrd.Total_payout,
	}

	return resp, nil
}

func (impl *SalesIncentiveImpl) FetchUserTargetReport(id int, t time.Time, quarter int, isQuarterly bool) (*db.ChartDataResponse, *commons.GenericError) {

	usrRcrd, err := impl.DBService.FetchIndividualUserTargetReport(id, int(t.Month()), quarter, isQuarterly)
	if err != nil {
		return nil, err
	}

	resp := &db.ChartDataResponse{
		HBB: &db.TargetAchievementData{
			SalesCount:         usrRcrd.Hbb_sales_count,
			Target:             usrRcrd.Hbb_target,
			AchievedPercentage: usrRcrd.Hbb_achieved_percentage,
		},
		Hayyak: &db.TargetAchievementData{
			SalesCount:         usrRcrd.Hayyak_sales_count,
			Target:             usrRcrd.Hayyak_target,
			AchievedPercentage: usrRcrd.Hayyak_achieved_percentage,
		},
		Baqati: &db.TargetAchievementData{
			SalesCount:         usrRcrd.Baqati_sales_count,
			Target:             usrRcrd.Baqati_target,
			AchievedPercentage: usrRcrd.Baqati_achieved_percentage,
		},
		Upgrade: &db.TargetAchievementData{
			SalesCount:         usrRcrd.Upgrade_sales_count,
			Target:             usrRcrd.Upgrade_target,
			AchievedPercentage: usrRcrd.Uprade_achieved_percentage,
		},
	}

	return resp, nil
}

func (impl *SalesIncentiveImpl) FetchSalesCountPerPlanReport(id int, group string, t time.Time, quarter int, isQuarterly bool) ([]*db.SaleCountPerPlan, *commons.GenericError) {

	usrRcrd, genricErr := impl.DBService.FetchIndividualUserSaleCountPerPlanReport(id, int(t.Month()), quarter, isQuarterly)
	if genricErr != nil {
		return nil, genricErr
	}

	str := ""
	response := make([]*db.SaleCountPerPlan, 0)

	switch group {
	case commons.HayyakCategory:
		str = usrRcrd.HayyakSaleCountPerPlanArray
	case commons.HBBCategory:
		str = usrRcrd.HBBSaleCountPerPlanArray
	case commons.BaqatiCategory:
		str = usrRcrd.BaqatiSaleCountPerPlanArray
	case commons.UpgradeCategory:
		str = usrRcrd.UpgradeSaleCountPerPlanArray
	}

	var m map[string]int

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	for key, val := range m {
		response = append(response, &db.SaleCountPerPlan{
			CurrentPlan: key,
			Quantity:    val,
		})
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Quantity > response[j].Quantity
	})

	return response, nil
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

func (impl *SalesIncentiveImpl) VefiryOTP(clientID, otp int) (*db.OtpModel, *commons.GenericError) {

	otpModel, err := impl.DBService.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	if otp != otpModel.OTP {
		return nil, &commons.GenericError{Code: 1022, ErrorMessage: "otp invalid"}
	}

	if !otpModel.ExpiryTime.After(time.Now().Local()) {
		return nil, &commons.GenericError{Code: 1021, ErrorMessage: "otp has expired"}
	}

	//retry count

	return otpModel, nil
}

func (impl *SalesIncentiveImpl) FetchMonthlyAndQuarterlyInfoAboutPromoter(clientID int) (map[int]*db.UserAllInfo, map[int]*db.UserAllInfo, error) {

	//TODO: not jst for same year but for everytime new sales record are added this should be updated

	monthly, quarterly, err := impl.DBService.FetchPromotersPresentYearSales(clientID)
	if err != nil {
		return nil, nil, err
	}

	monthlyRecords := make(map[int]*db.UserAllInfo)
	quarterlyRecords := make(map[int]*db.UserAllInfo)

	for month, trx := range monthly {

		currentMonth := time.Date(time.Now().Year(), time.Month(month), 1, 1, 1, 1, 1, time.Local)

		usrInfo := &db.UserAllInfo{
			TotalIncentive:               &db.TotalIncentive{},
			ChartDataResponse:            &db.ChartDataResponse{},
			HBBSaleCountPerPlanArray:     make(map[string]int),
			HayyakSaleCountPerPlanArray:  make(map[string]int),
			BaqatiSaleCountPerPlanArray:  make(map[string]int),
			UpgradeSaleCountPerPlanArray: make(map[string]int),
		}

		if trx.Region == "" || trx.Outlet == "" {
			// ind := &db.IndividualIncentive{Error: "no sales record found for this period"}
			// indTarget := &db.TargetAchievementData{}
			// usrInfo.TotalIncentive.HBB = ind
			// usrInfo.TotalIncentive.Baqati = ind
			// usrInfo.TotalIncentive.Hayyak = ind
			// usrInfo.TotalIncentive.Upgrade = ind
			// usrInfo.ChartDataResponse.HBB = indTarget
			// usrInfo.ChartDataResponse.Hayyak = indTarget
			// usrInfo.ChartDataResponse.Baqati = indTarget
			// usrInfo.ChartDataResponse.Upgrade = indTarget
			return nil, nil, errors.New(fmt.Sprintf("%v : region and outlet not found", clientID))
		}

		totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
		if err != nil {
			return nil, nil, err
		} else if totalEmp == 0 {
			// ind := &db.IndividualIncentive{Error: "no employee found for this outlet - " + trx.Outlet}
			// usrInfo.TotalIncentive.HBB = ind
			// usrInfo.TotalIncentive.Baqati = ind
			// usrInfo.TotalIncentive.Hayyak = ind
			// usrInfo.TotalIncentive.Upgrade = ind
			return nil, nil, errors.New(fmt.Sprintf("%v : no employee found for this outlet %v", clientID, trx.Outlet))
		}

		totalEmp--

		totalPay := 0.0

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

			targetPerOutlet := impl.getSaleTargetPerOutlet(currentMonth, 0, false, grp, trx.Region, trx.Outlet)
			//if targetPerOutlet == -1 {
			//	return nil, errors.New("outlet target not found")
			//}

			switch grp {

			case commons.HBBCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("HBB sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("HBB: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("HBB planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}

					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.HBBSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.HBB = ind
				usrInfo.ChartDataResponse.HBB = indivChartData

			case commons.HayyakCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("HAYYAK sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("HAYYAK: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("HAYYAK planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}

					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.HayyakSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.Hayyak = ind
				usrInfo.ChartDataResponse.Hayyak = indivChartData

			case commons.BaqatiCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("BAQATI sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("BAQATI: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("BAQATI planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}
					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.BaqatiSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.Baqati = ind
				usrInfo.ChartDataResponse.Baqati = indivChartData
			}
		}

		ind := &db.IndividualIncentive{}
		indivChartData := &db.TargetAchievementData{}

		targetPerOutlet := impl.getSaleTargetPerOutlet(currentMonth, 0, false, commons.UpgradeCategory, trx.Region, trx.Outlet)
		// if targetPerOutlet == -1 {
		// 	ind.Error = "outlet target not found"
		// 	usrInfo.TotalIncentive.Upgrade = ind
		// 	return nil, nil,
		// }

		// fmt.Println("UPGRADE sales target in "+trx.Outlet, targetPerOutlet)

		target := float64(targetPerOutlet) / float64(totalEmp)
		if target <= 0 {
			ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
		} else {
			achieved := float64(totalActiveUpgrade) * 100 / float64(target)
			vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)

			// fmt.Println("UPGRADE: target, achieved, percent, vop", target, totalActiveUpgrade, achieved, vop)

			indivChartData.Target = target
			indivChartData.SalesCount = totalActiveUpgrade
			indivChartData.AchievedPercentage = math.Round((float64(totalActiveUpgrade)*100/target)*100) / 100

			ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
			ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
			totalPay += ind.ExtraPayout

			for currPlan, prevPlansMap := range upgradedPlans {
				for prevPln, qty := range prevPlansMap {
					rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
					// fmt.Println("UPGRADE currPlanName prevPlanName  currPlanRate  prevPlanRate rate qty pay: "+currPlan+" "+prevPln, impl.getSIPRatePlan(currPlan), impl.getSIPRatePlan(prevPln), rate, qty, float64(qty)*vop*rate)
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

		for currPlan, prevPlansMap := range upgradedPlans {

			currPlanQty := 0

			for _, qty := range prevPlansMap {
				currPlanQty += qty
			}

			usrInfo.UpgradeSaleCountPerPlanArray[currPlan] += currPlanQty

		}

		usrInfo.TotalIncentive.Upgrade = ind
		usrInfo.TotalIncentive.TotalPay = math.Round(totalPay*100) / 100

		usrInfo.ChartDataResponse.Upgrade = indivChartData

		monthlyRecords[month] = usrInfo
	}

	/*
		Quarterly Report
	*/

	for quarter, trx := range quarterly {

		currentYear := time.Date(time.Now().Year(), 1, 1, 1, 1, 1, 1, time.Local)

		usrInfo := &db.UserAllInfo{
			TotalIncentive:               &db.TotalIncentive{},
			ChartDataResponse:            &db.ChartDataResponse{},
			HBBSaleCountPerPlanArray:     make(map[string]int),
			HayyakSaleCountPerPlanArray:  make(map[string]int),
			BaqatiSaleCountPerPlanArray:  make(map[string]int),
			UpgradeSaleCountPerPlanArray: make(map[string]int),
		}

		if trx.Region == "" || trx.Outlet == "" {
			// ind := &db.IndividualIncentive{Error: "no sales record found for this period"}
			// usrInfo.TotalIncentive.HBB = ind
			// usrInfo.TotalIncentive.Baqati = ind
			// usrInfo.TotalIncentive.Hayyak = ind
			// usrInfo.TotalIncentive.Upgrade = ind
			return nil, nil, errors.New("no sales record found for this period")
		}

		totalEmp, err := impl.DBService.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
		if err != nil {
			return nil, nil, err
		} else if totalEmp == 0 {
			// ind := &db.IndividualIncentive{Error: "no employee found for this outlet - " + trx.Outlet}
			// usrInfo.TotalIncentive.HBB = ind
			// usrInfo.TotalIncentive.Baqati = ind
			// usrInfo.TotalIncentive.Hayyak = ind
			// usrInfo.TotalIncentive.Upgrade = ind
			return nil, nil, errors.New("no employee found for this outlet - " + trx.Outlet)
		}

		totalEmp--

		totalPay := 0.0

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

			targetPerOutlet := impl.getSaleTargetPerOutlet(currentYear, quarter, true, grp, trx.Region, trx.Outlet)
			//if targetPerOutlet == -1 {
			//	return nil, errors.New("outlet target not found")
			//}

			switch grp {

			case commons.HBBCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("HBB sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for HBB", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("HBB: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("HBB planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}

					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.HBBSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.HBB = ind
				usrInfo.ChartDataResponse.HBB = indivChartData

			case commons.HayyakCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("HAYYAK sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for HAYYAK", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("HAYYAK: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("HAYYAK planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}

					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.HayyakSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.Hayyak = ind
				usrInfo.ChartDataResponse.Hayyak = indivChartData

			case commons.BaqatiCategory:

				ind := &db.IndividualIncentive{}
				indivChartData := &db.TargetAchievementData{}

				// fmt.Println("BAQATI sales target in "+trx.Outlet, targetPerOutlet)

				target := float64(targetPerOutlet) / float64(totalEmp)
				if target <= 0 {
					ind.Error = fmt.Errorf("per employee target found is %f for BAQATI", target).Error()
				} else {
					achieved := float64(totalActive[grp]) * 100 / float64(target)
					vop := impl.getValueOfPoint(grp, achieved)

					// fmt.Println("BAQATI: target, achieved, percent, vop", target, totalActive[grp], achieved, vop)

					indivChartData.Target = target
					indivChartData.SalesCount = totalActive[grp]
					indivChartData.AchievedPercentage = math.Round((float64(totalActive[grp])*100/target)*100) / 100

					ind.ExtraPayout = impl.getExtraPay(grp, achieved)
					ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
					totalPay += ind.ExtraPayout

					for plan, qty := range plans {
						rate := impl.getSIPRatePlan(plan)
						// fmt.Println("BAQATI planName  planRate qty pay: "+plan, rate, qty, float64(qty)*vop*rate)
						ind.Payout += float64(qty) * vop * rate
					}

					max := impl.getMaxPayout(grp)
					if ind.Payout > max {
						ind.Payout = max
					}
					ind.Payout = math.Round(ind.Payout*100) / 100
					totalPay += ind.Payout
				}

				for plan, qty := range plans {
					usrInfo.BaqatiSaleCountPerPlanArray[plan] += qty
				}

				usrInfo.TotalIncentive.Baqati = ind
				usrInfo.ChartDataResponse.Baqati = indivChartData
			}
		}

		ind := &db.IndividualIncentive{}
		indivChartData := &db.TargetAchievementData{}

		targetPerOutlet := impl.getSaleTargetPerOutlet(currentYear, quarter, true, commons.UpgradeCategory, trx.Region, trx.Outlet)
		// if targetPerOutlet == -1 {
		// 	ind.Error = "outlet target not found"
		// 	usrInfo.TotalIncentive.Upgrade = ind
		// 	return nil, nil
		// }

		// fmt.Println("UPGRADE sales target in "+trx.Outlet, targetPerOutlet)

		target := float64(targetPerOutlet) / float64(totalEmp)
		if target <= 0 {
			ind.Error = fmt.Errorf("per employee target found is %f for Upgrade", target).Error()
		} else {
			achieved := float64(totalActiveUpgrade) * 100 / float64(target)
			vop := impl.getValueOfPoint(commons.UpgradeCategory, achieved)

			// fmt.Println("UPGRADE: target, achieved, percent, vop", target, totalActiveUpgrade, achieved, vop)

			indivChartData.Target = target
			indivChartData.SalesCount = totalActiveUpgrade
			indivChartData.AchievedPercentage = math.Round((float64(totalActiveUpgrade)*100/target)*100) / 100

			ind.ExtraPayout = impl.getExtraPay(commons.UpgradeCategory, achieved)
			ind.ExtraPayout = math.Round(ind.ExtraPayout*100) / 100
			totalPay += ind.ExtraPayout

			for currPlan, prevPlansMap := range upgradedPlans {
				for prevPln, qty := range prevPlansMap {
					rate := impl.getSIPRatePlan(currPlan) - impl.getSIPRatePlan(prevPln)
					// fmt.Println("UPGRADE currPlanName prevPlanName  currPlanRate  prevPlanRate rate qty pay: "+currPlan+" "+prevPln, impl.getSIPRatePlan(currPlan), impl.getSIPRatePlan(prevPln), rate, qty, float64(qty)*vop*rate)
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

		for currPlan, prevPlansMap := range upgradedPlans {

			currPlanQty := 0

			for _, qty := range prevPlansMap {
				currPlanQty += qty
			}

			usrInfo.UpgradeSaleCountPerPlanArray[currPlan] += currPlanQty

		}

		usrInfo.TotalIncentive.Upgrade = ind
		usrInfo.TotalIncentive.TotalPay = math.Round(totalPay*100) / 100

		usrInfo.ChartDataResponse.Upgrade = indivChartData

		quarterlyRecords[quarter] = usrInfo
	}

	return monthlyRecords, quarterlyRecords, nil
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
