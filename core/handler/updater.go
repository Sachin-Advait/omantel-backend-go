package handler

import (
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
)

func (impl *SalesIncentiveImpl) UpdateEachPromoterData() *commons.GenericError {

	array, err := impl.DBService.FetchAllGivenPositionEmpIdAndTheirManager(commons.Promoter)
	if err != nil {
		return err
	}

	for _, emp := range array {

		for _, month := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      false,
				Month:            month,
			}

			totalInc, err := impl.FetchPromoterUserIncentive(emp.EmpID, time.Date(time.Now().Year(), time.Month(month), 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			//fmt.Printf("user1: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Month, totalInc)

			targetChart, err := impl.FetchPromoterTargetChart(emp.EmpID, time.Date(time.Now().Year(), time.Month(month), 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			//fmt.Printf("user2: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Month, targetChart)

			if targetChart.HBB != nil {
				rcrd.Hbb_target = targetChart.HBB.Target
				rcrd.Hbb_sales_count = targetChart.HBB.SalesCount
				rcrd.Hbb_achieved_percentage = targetChart.HBB.AchievedPercentage
			}

			if totalInc.HBB != nil {
				rcrd.Hbb_payout = totalInc.HBB.Payout
				rcrd.Hbb_extra_payout = totalInc.HBB.ExtraPayout
			}

			if targetChart.Hayyak != nil {
				rcrd.Hayyak_target = targetChart.Hayyak.Target
				rcrd.Hayyak_sales_count = targetChart.Hayyak.SalesCount
				rcrd.Hayyak_achieved_percentage = targetChart.Hayyak.AchievedPercentage
			}

			if totalInc.Hayyak != nil {
				rcrd.Hayyak_payout = totalInc.Hayyak.Payout
				rcrd.Hayyak_extra_payout = totalInc.Hayyak.ExtraPayout
			}

			if targetChart.Baqati != nil {
				rcrd.Baqati_target = targetChart.Baqati.Target
				rcrd.Baqati_sales_count = targetChart.Baqati.SalesCount
				rcrd.Baqati_achieved_percentage = targetChart.Baqati.AchievedPercentage
			}

			if totalInc.Baqati != nil {
				rcrd.Baqati_payout = totalInc.Baqati.Payout
				rcrd.Baqati_extra_payout = totalInc.Baqati.ExtraPayout
			}

			if targetChart.Upgrade != nil {
				rcrd.Upgrade_target = targetChart.Upgrade.Target
				rcrd.Upgrade_sales_count = targetChart.Upgrade.SalesCount
				rcrd.Uprade_achieved_percentage = targetChart.Upgrade.AchievedPercentage
			}

			if totalInc.Upgrade != nil {
				rcrd.Upgrade_payout = totalInc.Upgrade.Payout
				rcrd.Upgrade_extra_payout = totalInc.Upgrade.ExtraPayout
			}

			rcrd.Total_payout = totalInc.TotalPay

			if err := impl.DBService.InsertEmployeeRecord(rcrd); err != nil {
				return err
			}

		}

		for _, quarter := range []int{1, 2, 3, 4} {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      true,
				Quarter:          quarter,
			}

			totalInc, err := impl.FetchPromoterUserIncentive(emp.EmpID, time.Date(time.Now().Year(), 1, 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			//fmt.Printf("user3: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Quarter, totalInc)

			targetChart, err := impl.FetchPromoterTargetChart(emp.EmpID, time.Date(time.Now().Year(), 1, 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			//fmt.Printf("user4: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Quarter, targetChart)

			if targetChart.HBB != nil {
				rcrd.Hbb_target = targetChart.HBB.Target
				rcrd.Hbb_sales_count = targetChart.HBB.SalesCount
				rcrd.Hbb_achieved_percentage = targetChart.HBB.AchievedPercentage
			}

			if totalInc.HBB != nil {
				rcrd.Hbb_payout = totalInc.HBB.Payout
				rcrd.Hbb_extra_payout = totalInc.HBB.ExtraPayout
			}

			if targetChart.Hayyak != nil {
				rcrd.Hayyak_target = targetChart.Hayyak.Target
				rcrd.Hayyak_sales_count = targetChart.Hayyak.SalesCount
				rcrd.Hayyak_achieved_percentage = targetChart.Hayyak.AchievedPercentage
			}

			if totalInc.Hayyak != nil {
				rcrd.Hayyak_payout = totalInc.Hayyak.Payout
				rcrd.Hayyak_extra_payout = totalInc.Hayyak.ExtraPayout
			}

			if targetChart.Baqati != nil {
				rcrd.Baqati_target = targetChart.Baqati.Target
				rcrd.Baqati_sales_count = targetChart.Baqati.SalesCount
				rcrd.Baqati_achieved_percentage = targetChart.Baqati.AchievedPercentage
			}

			if totalInc.Baqati != nil {
				rcrd.Baqati_payout = totalInc.Baqati.Payout
				rcrd.Baqati_extra_payout = totalInc.Baqati.ExtraPayout
			}

			if targetChart.Upgrade != nil {
				rcrd.Upgrade_target = targetChart.Upgrade.Target
				rcrd.Upgrade_sales_count = targetChart.Upgrade.SalesCount
				rcrd.Uprade_achieved_percentage = targetChart.Upgrade.AchievedPercentage
			}

			if totalInc.Upgrade != nil {
				rcrd.Upgrade_payout = totalInc.Upgrade.Payout
				rcrd.Upgrade_extra_payout = totalInc.Upgrade.ExtraPayout
			}

			rcrd.Total_payout = totalInc.TotalPay

			if err := impl.DBService.InsertEmployeeRecord(rcrd); err != nil {
				return err
			}
		}
	}

	return nil
}

func (impl *SalesIncentiveImpl) UpdateEachOutletManagerData() *commons.GenericError {

	array, err := impl.DBService.FetchAllGivenPositionEmpIdAndTheirManager(commons.OutletManager)
	if err != nil {
		return err
	}

	for _, emp := range array {

		for _, month := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      false,
				Month:            month,
			}

			totalInc, err := impl.FetchOutletManagerIncentive(emp.EmpID, time.Date(time.Now().Year(), time.Month(month), 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			fmt.Printf("user1: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Month, totalInc)

			targetChart, err := impl.FetchOutletManagerTargetChart(emp.EmpID, time.Date(time.Now().Year(), time.Month(month), 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			fmt.Printf("user2: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Month, targetChart)

			if targetChart.HBB != nil {
				rcrd.Hbb_target = targetChart.HBB.Target
				rcrd.Hbb_sales_count = targetChart.HBB.SalesCount
				rcrd.Hbb_achieved_percentage = targetChart.HBB.AchievedPercentage
			}

			if totalInc.HBB != nil {
				rcrd.Hbb_payout = totalInc.HBB.Payout
				rcrd.Hbb_extra_payout = totalInc.HBB.ExtraPayout
			}

			if targetChart.Hayyak != nil {
				rcrd.Hayyak_target = targetChart.Hayyak.Target
				rcrd.Hayyak_sales_count = targetChart.Hayyak.SalesCount
				rcrd.Hayyak_achieved_percentage = targetChart.Hayyak.AchievedPercentage
			}

			if totalInc.Hayyak != nil {
				rcrd.Hayyak_payout = totalInc.Hayyak.Payout
				rcrd.Hayyak_extra_payout = totalInc.Hayyak.ExtraPayout
			}

			if targetChart.Baqati != nil {
				rcrd.Baqati_target = targetChart.Baqati.Target
				rcrd.Baqati_sales_count = targetChart.Baqati.SalesCount
				rcrd.Baqati_achieved_percentage = targetChart.Baqati.AchievedPercentage
			}

			if totalInc.Baqati != nil {
				rcrd.Baqati_payout = totalInc.Baqati.Payout
				rcrd.Baqati_extra_payout = totalInc.Baqati.ExtraPayout
			}

			if targetChart.Upgrade != nil {
				rcrd.Upgrade_target = targetChart.Upgrade.Target
				rcrd.Upgrade_sales_count = targetChart.Upgrade.SalesCount
				rcrd.Uprade_achieved_percentage = targetChart.Upgrade.AchievedPercentage
			}

			if totalInc.Upgrade != nil {
				rcrd.Upgrade_payout = totalInc.Upgrade.Payout
				rcrd.Upgrade_extra_payout = totalInc.Upgrade.ExtraPayout
			}

			rcrd.Total_payout = totalInc.TotalPay

			if err := impl.DBService.InsertEmployeeRecord(rcrd); err != nil {
				return err
			}

		}

		for _, quarter := range []int{1, 2, 3, 4} {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      true,
				Quarter:          quarter,
			}

			totalInc, err := impl.FetchOutletManagerIncentive(emp.EmpID, time.Date(time.Now().Year(), 1, 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			fmt.Printf("user3: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Quarter, totalInc)

			targetChart, err := impl.FetchPromoterTargetChart(emp.EmpID, time.Date(time.Now().Year(), 1, 1, 1, 1, 1, 1, time.Local), rcrd.Quarter, rcrd.Isquarterly)
			if err != nil {
				return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
			}

			fmt.Printf("user4: %v %v %v %+v\n", emp.EmpID, rcrd.Isquarterly, rcrd.Quarter, targetChart)

			if targetChart.HBB != nil {
				rcrd.Hbb_target = targetChart.HBB.Target
				rcrd.Hbb_sales_count = targetChart.HBB.SalesCount
				rcrd.Hbb_achieved_percentage = targetChart.HBB.AchievedPercentage
			}

			if totalInc.HBB != nil {
				rcrd.Hbb_payout = totalInc.HBB.Payout
				rcrd.Hbb_extra_payout = totalInc.HBB.ExtraPayout
			}

			if targetChart.Hayyak != nil {
				rcrd.Hayyak_target = targetChart.Hayyak.Target
				rcrd.Hayyak_sales_count = targetChart.Hayyak.SalesCount
				rcrd.Hayyak_achieved_percentage = targetChart.Hayyak.AchievedPercentage
			}

			if totalInc.Hayyak != nil {
				rcrd.Hayyak_payout = totalInc.Hayyak.Payout
				rcrd.Hayyak_extra_payout = totalInc.Hayyak.ExtraPayout
			}

			if targetChart.Baqati != nil {
				rcrd.Baqati_target = targetChart.Baqati.Target
				rcrd.Baqati_sales_count = targetChart.Baqati.SalesCount
				rcrd.Baqati_achieved_percentage = targetChart.Baqati.AchievedPercentage
			}

			if totalInc.Baqati != nil {
				rcrd.Baqati_payout = totalInc.Baqati.Payout
				rcrd.Baqati_extra_payout = totalInc.Baqati.ExtraPayout
			}

			if targetChart.Upgrade != nil {
				rcrd.Upgrade_target = targetChart.Upgrade.Target
				rcrd.Upgrade_sales_count = targetChart.Upgrade.SalesCount
				rcrd.Uprade_achieved_percentage = targetChart.Upgrade.AchievedPercentage
			}

			if totalInc.Upgrade != nil {
				rcrd.Upgrade_payout = totalInc.Upgrade.Payout
				rcrd.Upgrade_extra_payout = totalInc.Upgrade.ExtraPayout
			}

			rcrd.Total_payout = totalInc.TotalPay

			if err := impl.DBService.InsertEmployeeRecord(rcrd); err != nil {
				return err
			}
		}
	}

	return nil
}

func (impl *SalesIncentiveImpl) UpdateOtherManagersData() *commons.GenericError {

	managerMap := map[string]string{
		commons.OutletManager:  "outlet_manager",
		commons.Manager:        "manager",
		commons.SeniorManager:  "senior_manager",
		commons.GeneralManager: "general_manager",
		commons.VicePresident:  "vice_president",
	}

	for key, value := range managerMap {

		managers, err := impl.DBService.FetchAllGivenPositionEmpIdAndTheirManager(key)
		if err != nil {
			return err
		}

		for _, manager := range managers {

			if manager.EmpID == 0 {
				continue
			}

			array, err := impl.DBService.FetchAllEmployeesReportingToGivenManagerRole(value, manager.EmpID)
			if err != nil {
				return err
			}

			monthlyRcrds := make(map[int]*db.UserRecord)
			quarterlyRcrds := make(map[int]*db.UserRecord)

			for _, id := range array {

				val, err := impl.DBService.FetchIndividualUserRecord(id)
				if err != nil {
					return err
				}

				for _, v := range val {

					if key == commons.OutletManager {
						fmt.Println("user: ", manager.EmpID, v.EmpId, v.Isquarterly, v.Month, v.Quarter)
					}

					if !v.Isquarterly {
						c, ok := monthlyRcrds[v.Month]
						if !ok {
							monthlyRcrds[v.Month] = &db.UserRecord{
								EmpId:                      manager.EmpID,
								Created:                    time.Now(),
								ReportingManager:           manager.ReportTo,
								Isquarterly:                false,
								Month:                      v.Month,
								Quarter:                    v.Quarter,
								Hayyak_sales_count:         v.Hayyak_sales_count,
								Hayyak_target:              v.Hayyak_target,
								Hayyak_payout:              v.Hayyak_payout,
								Hayyak_extra_payout:        v.Hayyak_extra_payout,
								Hayyak_achieved_percentage: v.Hayyak_achieved_percentage,
								Hbb_sales_count:            v.Hbb_sales_count,
								Hbb_target:                 v.Hbb_target,
								Hbb_payout:                 v.Hbb_payout,
								Hbb_extra_payout:           v.Hbb_extra_payout,
								Hbb_achieved_percentage:    v.Hbb_achieved_percentage,
								Baqati_sales_count:         v.Baqati_sales_count,
								Baqati_target:              v.Baqati_target,
								Baqati_payout:              v.Baqati_payout,
								Baqati_extra_payout:        v.Baqati_extra_payout,
								Baqati_achieved_percentage: v.Baqati_achieved_percentage,
								Upgrade_sales_count:        v.Upgrade_sales_count,
								Upgrade_target:             v.Upgrade_target,
								Upgrade_payout:             v.Upgrade_payout,
								Upgrade_extra_payout:       v.Upgrade_extra_payout,
								Uprade_achieved_percentage: v.Uprade_achieved_percentage,
								Total_payout:               v.Total_payout,
							}
						} else {
							c.Hayyak_sales_count += v.Hayyak_sales_count
							c.Hayyak_target += v.Hayyak_target
							c.Hayyak_payout += v.Hayyak_payout
							c.Hayyak_extra_payout += v.Hayyak_extra_payout
							c.Hayyak_achieved_percentage += v.Hayyak_achieved_percentage
							c.Hbb_sales_count += v.Hbb_sales_count
							c.Hbb_target += v.Hbb_target
							c.Hbb_payout += v.Hbb_payout
							c.Hbb_extra_payout += v.Hbb_extra_payout
							c.Hbb_achieved_percentage += v.Hbb_achieved_percentage
							c.Baqati_sales_count += v.Baqati_sales_count
							c.Baqati_target += v.Baqati_target
							c.Baqati_payout += v.Baqati_payout
							c.Baqati_extra_payout += v.Baqati_extra_payout
							c.Baqati_achieved_percentage += v.Baqati_achieved_percentage
							c.Upgrade_sales_count += v.Upgrade_sales_count
							c.Upgrade_target += v.Upgrade_target
							c.Upgrade_payout += v.Upgrade_payout
							c.Upgrade_extra_payout += v.Upgrade_extra_payout
							c.Uprade_achieved_percentage += v.Uprade_achieved_percentage
							c.Total_payout += v.Total_payout
						}
					} else {
						c, ok := quarterlyRcrds[v.Quarter]
						if !ok {
							quarterlyRcrds[v.Month] = &db.UserRecord{
								EmpId:                      manager.EmpID,
								Created:                    time.Now(),
								ReportingManager:           manager.ReportTo,
								Isquarterly:                true,
								Month:                      v.Month,
								Quarter:                    v.Quarter,
								Hayyak_sales_count:         v.Hayyak_sales_count,
								Hayyak_target:              v.Hayyak_target,
								Hayyak_payout:              v.Hayyak_payout,
								Hayyak_extra_payout:        v.Hayyak_extra_payout,
								Hayyak_achieved_percentage: v.Hayyak_achieved_percentage,
								Hbb_sales_count:            v.Hbb_sales_count,
								Hbb_target:                 v.Hbb_target,
								Hbb_payout:                 v.Hbb_payout,
								Hbb_extra_payout:           v.Hbb_extra_payout,
								Hbb_achieved_percentage:    v.Hbb_achieved_percentage,
								Baqati_sales_count:         v.Baqati_sales_count,
								Baqati_target:              v.Baqati_target,
								Baqati_payout:              v.Baqati_payout,
								Baqati_extra_payout:        v.Baqati_extra_payout,
								Baqati_achieved_percentage: v.Baqati_achieved_percentage,
								Upgrade_sales_count:        v.Upgrade_sales_count,
								Upgrade_target:             v.Upgrade_target,
								Upgrade_payout:             v.Upgrade_payout,
								Upgrade_extra_payout:       v.Upgrade_extra_payout,
								Uprade_achieved_percentage: v.Uprade_achieved_percentage,
								Total_payout:               v.Total_payout,
							}
						} else {
							c.Hayyak_sales_count += v.Hayyak_sales_count
							c.Hayyak_target += v.Hayyak_target
							c.Hayyak_payout += v.Hayyak_payout
							c.Hayyak_extra_payout += v.Hayyak_extra_payout
							c.Hayyak_achieved_percentage += v.Hayyak_achieved_percentage
							c.Hbb_sales_count += v.Hbb_sales_count
							c.Hbb_target += v.Hbb_target
							c.Hbb_payout += v.Hbb_payout
							c.Hbb_extra_payout += v.Hbb_extra_payout
							c.Hbb_achieved_percentage += v.Hbb_achieved_percentage
							c.Baqati_sales_count += v.Baqati_sales_count
							c.Baqati_target += v.Baqati_target
							c.Baqati_payout += v.Baqati_payout
							c.Baqati_extra_payout += v.Baqati_extra_payout
							c.Baqati_achieved_percentage += v.Baqati_achieved_percentage
							c.Upgrade_sales_count += v.Upgrade_sales_count
							c.Upgrade_target += v.Upgrade_target
							c.Upgrade_payout += v.Upgrade_payout
							c.Upgrade_extra_payout += v.Upgrade_extra_payout
							c.Uprade_achieved_percentage += v.Uprade_achieved_percentage
							c.Total_payout += v.Total_payout
						}
					}
				}
			}

			for _, rcd := range monthlyRcrds {
				rcd.Hayyak_achieved_percentage = rcd.Hayyak_achieved_percentage / float64(len(array))
				rcd.Hbb_achieved_percentage = rcd.Hbb_achieved_percentage / float64(len(array))
				rcd.Baqati_achieved_percentage = rcd.Baqati_achieved_percentage / float64(len(array))
				rcd.Uprade_achieved_percentage = rcd.Uprade_achieved_percentage / float64(len(array))

				if err := impl.DBService.InsertEmployeeRecord(rcd); err != nil {
					return err
				}
			}

			for _, rcd := range quarterlyRcrds {
				rcd.Hayyak_achieved_percentage = rcd.Hayyak_achieved_percentage / float64(len(array))
				rcd.Hbb_achieved_percentage = rcd.Hbb_achieved_percentage / float64(len(array))
				rcd.Baqati_achieved_percentage = rcd.Baqati_achieved_percentage / float64(len(array))
				rcd.Uprade_achieved_percentage = rcd.Uprade_achieved_percentage / float64(len(array))

				if err := impl.DBService.InsertEmployeeRecord(rcd); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
