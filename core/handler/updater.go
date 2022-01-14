package handler

import (
	"encoding/json"
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

		rcrds := make([]*db.UserRecord, 0)

		monthlyData, quarterlyData, err := impl.FetchMonthlyAndQuarterlyInfoAboutPromoter(emp.EmpID)
		if err != nil {
			return &commons.GenericError{ErrorMessage: err.Error(), Code: 1001}
		}

		for month, usrData := range monthlyData {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      false,
				Month:            month,
			}

			if usrData.TotalIncentive != nil {

				if usrData.TotalIncentive.Baqati != nil {
					rcrd.Baqati_payout = usrData.TotalIncentive.Baqati.Payout
					rcrd.Baqati_extra_payout = usrData.TotalIncentive.Baqati.ExtraPayout
				}

				if usrData.TotalIncentive.Hayyak != nil {
					rcrd.Hayyak_payout = usrData.TotalIncentive.Hayyak.Payout
					rcrd.Hayyak_extra_payout = usrData.TotalIncentive.Hayyak.ExtraPayout
				}

				if usrData.TotalIncentive.HBB != nil {
					rcrd.Hbb_payout = usrData.TotalIncentive.HBB.Payout
					rcrd.Hbb_extra_payout = usrData.TotalIncentive.HBB.ExtraPayout
				}

				if usrData.TotalIncentive.Upgrade != nil {
					rcrd.Upgrade_payout = usrData.TotalIncentive.Upgrade.Payout
					rcrd.Upgrade_extra_payout = usrData.TotalIncentive.Upgrade.ExtraPayout
				}

				rcrd.Total_payout = usrData.TotalIncentive.TotalPay
			}

			if usrData.ChartDataResponse != nil {

				if usrData.ChartDataResponse.Baqati != nil {
					rcrd.Baqati_sales_count = usrData.ChartDataResponse.Baqati.SalesCount
					rcrd.Baqati_target = usrData.ChartDataResponse.Baqati.Target
					rcrd.Baqati_achieved_percentage = usrData.ChartDataResponse.Baqati.AchievedPercentage
				}

				if usrData.ChartDataResponse.Hayyak != nil {
					rcrd.Hayyak_sales_count = usrData.ChartDataResponse.Hayyak.SalesCount
					rcrd.Hayyak_target = usrData.ChartDataResponse.Hayyak.Target
					rcrd.Hayyak_achieved_percentage = usrData.ChartDataResponse.Hayyak.AchievedPercentage
				}

				if usrData.ChartDataResponse.HBB != nil {
					rcrd.Hbb_sales_count = usrData.ChartDataResponse.HBB.SalesCount
					rcrd.Hbb_target = usrData.ChartDataResponse.HBB.Target
					rcrd.Hbb_achieved_percentage = usrData.ChartDataResponse.HBB.AchievedPercentage
				}

				if usrData.ChartDataResponse.Upgrade != nil {
					rcrd.Upgrade_sales_count = usrData.ChartDataResponse.Upgrade.SalesCount
					rcrd.Upgrade_target = usrData.ChartDataResponse.Upgrade.Target
					rcrd.Uprade_achieved_percentage = usrData.ChartDataResponse.Upgrade.AchievedPercentage
				}
			}

			if usrData.BaqatiSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.BaqatiSaleCountPerPlanArray)
				if err != nil {
					rcrd.BaqatiSaleCountPerPlanArray = "{}"
				} else {
					rcrd.BaqatiSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.HayyakSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.HayyakSaleCountPerPlanArray)
				if err != nil {
					rcrd.HayyakSaleCountPerPlanArray = "{}"
				} else {
					rcrd.HayyakSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.HBBSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.HBBSaleCountPerPlanArray)
				if err != nil {
					rcrd.HBBSaleCountPerPlanArray = "{}"
				} else {
					rcrd.HBBSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.UpgradeSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.UpgradeSaleCountPerPlanArray)
				if err != nil {
					rcrd.UpgradeSaleCountPerPlanArray = "{}"
				} else {
					rcrd.UpgradeSaleCountPerPlanArray = string(bytes)
				}
			}

			rcrds = append(rcrds, rcrd)
		}

		/*
			Quarterly Report
		*/

		for quarter, usrData := range quarterlyData {

			rcrd := &db.UserRecord{
				EmpId:            emp.EmpID,
				Created:          time.Now(),
				ReportingManager: emp.ReportTo,
				Isquarterly:      true,
				Quarter:          quarter,
			}

			if usrData.TotalIncentive != nil {

				if usrData.TotalIncentive.Baqati != nil {
					rcrd.Baqati_payout = usrData.TotalIncentive.Baqati.Payout
					rcrd.Baqati_extra_payout = usrData.TotalIncentive.Baqati.ExtraPayout
				}

				if usrData.TotalIncentive.Hayyak != nil {
					rcrd.Hayyak_payout = usrData.TotalIncentive.Hayyak.Payout
					rcrd.Hayyak_extra_payout = usrData.TotalIncentive.Hayyak.ExtraPayout
				}

				if usrData.TotalIncentive.HBB != nil {
					rcrd.Hbb_payout = usrData.TotalIncentive.HBB.Payout
					rcrd.Hbb_extra_payout = usrData.TotalIncentive.HBB.ExtraPayout
				}

				if usrData.TotalIncentive.Upgrade != nil {
					rcrd.Upgrade_payout = usrData.TotalIncentive.Upgrade.Payout
					rcrd.Upgrade_extra_payout = usrData.TotalIncentive.Upgrade.ExtraPayout
				}

				rcrd.Total_payout = usrData.TotalIncentive.TotalPay
			}

			if usrData.ChartDataResponse != nil {

				if usrData.ChartDataResponse.Baqati != nil {
					rcrd.Baqati_sales_count = usrData.ChartDataResponse.Baqati.SalesCount
					rcrd.Baqati_target = usrData.ChartDataResponse.Baqati.Target
					rcrd.Baqati_achieved_percentage = usrData.ChartDataResponse.Baqati.AchievedPercentage
				}

				if usrData.ChartDataResponse.Hayyak != nil {
					rcrd.Hayyak_sales_count = usrData.ChartDataResponse.Hayyak.SalesCount
					rcrd.Hayyak_target = usrData.ChartDataResponse.Hayyak.Target
					rcrd.Hayyak_achieved_percentage = usrData.ChartDataResponse.Hayyak.AchievedPercentage
				}

				if usrData.ChartDataResponse.HBB != nil {
					rcrd.Hbb_sales_count = usrData.ChartDataResponse.HBB.SalesCount
					rcrd.Hbb_target = usrData.ChartDataResponse.HBB.Target
					rcrd.Hbb_achieved_percentage = usrData.ChartDataResponse.HBB.AchievedPercentage
				}

				if usrData.ChartDataResponse.Upgrade != nil {
					rcrd.Upgrade_sales_count = usrData.ChartDataResponse.Upgrade.SalesCount
					rcrd.Upgrade_target = usrData.ChartDataResponse.Upgrade.Target
					rcrd.Uprade_achieved_percentage = usrData.ChartDataResponse.Upgrade.AchievedPercentage
				}
			}

			if usrData.BaqatiSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.BaqatiSaleCountPerPlanArray)
				if err != nil {
					rcrd.BaqatiSaleCountPerPlanArray = "{}"
				} else {
					rcrd.BaqatiSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.HayyakSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.HayyakSaleCountPerPlanArray)
				if err != nil {
					rcrd.HayyakSaleCountPerPlanArray = "{}"
				} else {
					rcrd.HayyakSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.HBBSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.HBBSaleCountPerPlanArray)
				if err != nil {
					rcrd.HBBSaleCountPerPlanArray = "{}"
				} else {
					rcrd.HBBSaleCountPerPlanArray = string(bytes)
				}
			}

			if usrData.UpgradeSaleCountPerPlanArray != nil {
				bytes, err := json.Marshal(&usrData.UpgradeSaleCountPerPlanArray)
				if err != nil {
					rcrd.UpgradeSaleCountPerPlanArray = "{}"
				} else {
					rcrd.UpgradeSaleCountPerPlanArray = string(bytes)
				}
			}

			rcrds = append(rcrds, rcrd)

		}

		if len(rcrds) == 0 {
			fmt.Println(emp.EmpID, "no records found")
			continue
		}

		if err := impl.DBService.InsertEmployeeMultipleRecords(rcrds); err != nil {
			return err
		}
	}

	return nil
}

func (impl *SalesIncentiveImpl) UpdateOtherManagersData() *commons.GenericError {

	for key, value := range commons.ManagerMap {

		managers, err := impl.DBService.FetchAllGivenPositionEmpIdAndTheirManager(key)
		if err != nil {
			return err
		}

		for _, manager := range managers {

			if manager.EmpID == 0 {
				fmt.Println("not found for " + key)
				continue
			}

			monthlySalesCount := make(map[int]map[string]map[string]int)
			quarterlySalesCount := make(map[int]map[string]map[string]int)

			array, err := impl.DBService.FetchAllEmployeesReportingToGivenManagerRole(value, manager.EmpID)
			if err != nil {
				return err
			}

			totalRecords := make([]*db.UserRecord, 0)
			monthlyRcrds := make(map[int]*db.UserRecord)
			quarterlyRcrds := make(map[int]*db.UserRecord)

			for _, usr := range array {

				allRecords, err := impl.DBService.FetchIndividualUserRecord(usr.ID)
				if err != nil {
					return err
				}

				for _, eachRecord := range allRecords {

					if !eachRecord.Isquarterly {
						c, ok := monthlyRcrds[eachRecord.Month]
						if !ok {
							c = &db.UserRecord{
								EmpId:            manager.EmpID,
								Created:          time.Now(),
								ReportingManager: manager.ReportTo,
								Isquarterly:      false,
								Month:            eachRecord.Month,
							}

							monthlySalesCount[eachRecord.Month] = map[string]map[string]int{
								commons.HayyakCategory:  make(map[string]int),
								commons.HBBCategory:     make(map[string]int),
								commons.BaqatiCategory:  make(map[string]int),
								commons.UpgradeCategory: make(map[string]int),
							}
						}
						c.Hayyak_sales_count += eachRecord.Hayyak_sales_count
						c.Hayyak_target += eachRecord.Hayyak_target
						c.Hayyak_payout += eachRecord.Hayyak_payout
						c.Hayyak_extra_payout += eachRecord.Hayyak_extra_payout
						c.Hayyak_achieved_percentage += eachRecord.Hayyak_achieved_percentage
						c.Hbb_sales_count += eachRecord.Hbb_sales_count
						c.Hbb_target += eachRecord.Hbb_target
						c.Hbb_payout += eachRecord.Hbb_payout
						c.Hbb_extra_payout += eachRecord.Hbb_extra_payout
						c.Hbb_achieved_percentage += eachRecord.Hbb_achieved_percentage
						c.Baqati_sales_count += eachRecord.Baqati_sales_count
						c.Baqati_target += eachRecord.Baqati_target
						c.Baqati_payout += eachRecord.Baqati_payout
						c.Baqati_extra_payout += eachRecord.Baqati_extra_payout
						c.Baqati_achieved_percentage += eachRecord.Baqati_achieved_percentage
						c.Upgrade_sales_count += eachRecord.Upgrade_sales_count
						c.Upgrade_target += eachRecord.Upgrade_target
						c.Upgrade_payout += eachRecord.Upgrade_payout
						c.Upgrade_extra_payout += eachRecord.Upgrade_extra_payout
						c.Uprade_achieved_percentage += eachRecord.Uprade_achieved_percentage
						c.Total_payout += eachRecord.Total_payout

						temp := make(map[string]int)
						err := json.Unmarshal([]byte(eachRecord.HayyakSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								monthlySalesCount[eachRecord.Month][commons.HayyakCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.HBBSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								monthlySalesCount[eachRecord.Month][commons.HBBCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.BaqatiSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								monthlySalesCount[eachRecord.Month][commons.BaqatiCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.UpgradeSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								monthlySalesCount[eachRecord.Month][commons.UpgradeCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						monthlyRcrds[eachRecord.Month] = c

					} else {
						c, ok := quarterlyRcrds[eachRecord.Quarter]
						if !ok {
							c = &db.UserRecord{
								EmpId:            manager.EmpID,
								Created:          time.Now(),
								ReportingManager: manager.ReportTo,
								Isquarterly:      true,
								Quarter:          eachRecord.Quarter,
							}

							quarterlySalesCount[eachRecord.Month] = map[string]map[string]int{
								commons.HayyakCategory:  make(map[string]int),
								commons.HBBCategory:     make(map[string]int),
								commons.BaqatiCategory:  make(map[string]int),
								commons.UpgradeCategory: make(map[string]int),
							}
						}
						c.Hayyak_sales_count += eachRecord.Hayyak_sales_count
						c.Hayyak_target += eachRecord.Hayyak_target
						c.Hayyak_payout += eachRecord.Hayyak_payout
						c.Hayyak_extra_payout += eachRecord.Hayyak_extra_payout
						c.Hayyak_achieved_percentage += eachRecord.Hayyak_achieved_percentage
						c.Hbb_sales_count += eachRecord.Hbb_sales_count
						c.Hbb_target += eachRecord.Hbb_target
						c.Hbb_payout += eachRecord.Hbb_payout
						c.Hbb_extra_payout += eachRecord.Hbb_extra_payout
						c.Hbb_achieved_percentage += eachRecord.Hbb_achieved_percentage
						c.Baqati_sales_count += eachRecord.Baqati_sales_count
						c.Baqati_target += eachRecord.Baqati_target
						c.Baqati_payout += eachRecord.Baqati_payout
						c.Baqati_extra_payout += eachRecord.Baqati_extra_payout
						c.Baqati_achieved_percentage += eachRecord.Baqati_achieved_percentage
						c.Upgrade_sales_count += eachRecord.Upgrade_sales_count
						c.Upgrade_target += eachRecord.Upgrade_target
						c.Upgrade_payout += eachRecord.Upgrade_payout
						c.Upgrade_extra_payout += eachRecord.Upgrade_extra_payout
						c.Uprade_achieved_percentage += eachRecord.Uprade_achieved_percentage
						c.Total_payout += eachRecord.Total_payout

						temp := make(map[string]int)
						err := json.Unmarshal([]byte(eachRecord.HayyakSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								quarterlySalesCount[eachRecord.Month][commons.HayyakCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.HBBSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								quarterlySalesCount[eachRecord.Month][commons.HBBCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.BaqatiSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								quarterlySalesCount[eachRecord.Month][commons.BaqatiCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						temp = make(map[string]int)
						err = json.Unmarshal([]byte(eachRecord.UpgradeSaleCountPerPlanArray), &temp)
						if err == nil {
							for plan, qty := range temp {
								quarterlySalesCount[eachRecord.Month][commons.UpgradeCategory][plan] += qty
							}
						} else {
							//return &commons.GenericError{}
							fmt.Println(err.Error())
						}

						quarterlyRcrds[eachRecord.Quarter] = c
					}
				}
			}

			for m, rcd := range monthlyRcrds {
				rcd.Hayyak_achieved_percentage = rcd.Hayyak_achieved_percentage / float64(len(array))
				rcd.Hbb_achieved_percentage = rcd.Hbb_achieved_percentage / float64(len(array))
				rcd.Baqati_achieved_percentage = rcd.Baqati_achieved_percentage / float64(len(array))
				rcd.Uprade_achieved_percentage = rcd.Uprade_achieved_percentage / float64(len(array))

				byteArray, _ := json.Marshal(monthlySalesCount[m][commons.HayyakCategory])
				rcd.HayyakSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(monthlySalesCount[m][commons.HBBCategory])
				rcd.HBBSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(monthlySalesCount[m][commons.BaqatiCategory])
				rcd.BaqatiSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(monthlySalesCount[m][commons.UpgradeCategory])
				rcd.UpgradeSaleCountPerPlanArray = string(byteArray)

				totalRecords = append(totalRecords, rcd)
			}

			for q, rcd := range quarterlyRcrds {
				rcd.Hayyak_achieved_percentage = rcd.Hayyak_achieved_percentage / float64(len(array))
				rcd.Hbb_achieved_percentage = rcd.Hbb_achieved_percentage / float64(len(array))
				rcd.Baqati_achieved_percentage = rcd.Baqati_achieved_percentage / float64(len(array))
				rcd.Uprade_achieved_percentage = rcd.Uprade_achieved_percentage / float64(len(array))

				byteArray, _ := json.Marshal(quarterlySalesCount[q][commons.HayyakCategory])
				rcd.HayyakSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(quarterlySalesCount[q][commons.HBBCategory])
				rcd.HBBSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(quarterlySalesCount[q][commons.BaqatiCategory])
				rcd.BaqatiSaleCountPerPlanArray = string(byteArray)
				byteArray, _ = json.Marshal(quarterlySalesCount[q][commons.UpgradeCategory])
				rcd.UpgradeSaleCountPerPlanArray = string(byteArray)

				totalRecords = append(totalRecords, rcd)
			}

			if err := impl.DBService.InsertEmployeeMultipleRecords(totalRecords); err != nil {
				return err
			}
		}
	}

	return nil
}
