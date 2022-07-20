package business

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"log"
	"strconv"
	"strings"
	"time"
)

func (impl *IncentiveImpl) HandleUserMasterInputFile() {

	fileName := "./latest/master.xlsx"

	file, err := excelize.OpenFile(fileName)
	if err != nil {
		logger.FileLogger.Error().Err(err).Str("filename", fileName).Msg("HandleUserMasterInputFile : open file error")
	}

	rows, err := file.Rows("CSR & Promoter")
	if err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleUserMasterInputFile : get rows error")
		return
	}

	//records := make([]*db.Transaction, 0)

	isFirstRow := true
	index := 0

	for rows.Next() {

		index++

		if isFirstRow {
			rows.Columns()
			isFirstRow = false
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleUserMasterInputFile : get columns error")
			return
		}

		emp := &psqldb.Employee{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				num, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleUserMasterInputFile : ID conversion error")
					//break LOOP
					return
				}
				emp.Id = num
			case 1:
				emp.Name = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				emp.Region = strings.ToUpper(strings.TrimSpace(colCell))
			case 3:
				emp.Outlet = strings.ToUpper(strings.TrimSpace(colCell))
			case 4:
				emp.Position = strings.ToUpper(strings.TrimSpace(colCell))
			case 5:
				emp.Remark = strings.ToUpper(strings.TrimSpace(colCell))
			}
		}

		if emp.Name == "" || emp.Region == "" || emp.Outlet == "" || emp.Position == "" {
			logger.FileLogger.Error().Err(errors.New("invalid record")).Interface("emp", emp).Msg("HandleUserMasterInputFile : employee invalid record")
			return
		}

		err = impl.PsqlDb.InsertUserMasterRecord(emp)
		if err != nil {
			logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleUserMasterInputFile : insertionError")
			return
		}
	}

	logger.FileLogger.Info().Msg("HandleUserMasterInputFile : successful")
}

func (impl *IncentiveImpl) HandleDataInputFile() {
	//fileName := "doc/new_doc/q2_2021_monthly_data_orig.xlsx"
	fileName := "doc/latest/latest_Q2_2021_Monthly_data_updated.xlsx"

	file, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatalf("could not open file, %v\n", err)
	}

	rows, err := file.Rows("Data Q2 2021")
	if err != nil {
		fmt.Println(err)
		return
	}

	//records := make([]*db.Transaction, 0)

	isFirstRow := true
	index := 0

	for rows.Next() {

		index++

		if isFirstRow {
			rows.Columns()
			isFirstRow = false
			continue
		}

		eachRowColumns, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			return
		}

		trx := &psqldb.Transaction{}

		//LOOP:
		for i, colCell := range eachRowColumns {

			switch i {

			case 0:
				trx.Reqion = strings.ToUpper(strings.TrimSpace(colCell))
			case 1:
				trx.Outlet = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				num, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					log.Println("ID conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
					//break LOOP
					return
				}
				trx.Id = num
			case 3:
				trx.CategoryType = strings.ToUpper(strings.TrimSpace(colCell))
			case 4:
				t, err := time.Parse("01-02-06", strings.TrimSpace(colCell))
				if err != nil {
					log.Println("date conversion error: " + err.Error() + " index: " + strconv.Itoa(index) + " cell " + strings.TrimSpace(colCell))
					//break LOOP
					return
				}
				trx.SaleDate = t
				trx.Quarter = impl.getQuarterFromMonth(int(t.Month()))
			case 5:
				// trx.PreviousProduct = strings.ToUpper(strings.TrimSpace(colCell))
				trx.PreviousProduct = strings.TrimSpace(colCell)
			case 6:
				// trx.CurrentProduct = strings.ToUpper(strings.TrimSpace(colCell))
				trx.CurrentProduct = strings.TrimSpace(colCell)
			case 7:
				num, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					log.Println("Quantity conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
					//break LOOP
					return
				}
				trx.Quantity = num
			case 8:
				trx.Category = strings.ToUpper(strings.TrimSpace(colCell))
			case 9:
				trx.ProductGrp = strings.ToUpper(strings.TrimSpace(colCell))
				// case 10:
				// 	t, err := time.Parse("", strings.TrimSpace(colCell))
				// 	if err != nil {
				// 		log.Println("Incorrect sales time error: " + err.Error() + " index: " + strconv.Itoa(index))
				// 		//break LOOP
				// 		return
				// 	}

				// 	trx.SaleDate = t
				// 	trx.Quarter = impl.getQuarterFromMonth(int(t.Month()))
			}
		}

		if trx.Reqion == "" || trx.Outlet == "" || trx.CategoryType == "" || trx.CurrentProduct == "" || trx.Category == "" || trx.ProductGrp == "" || trx.SaleDate.Year() == 1 {
			log.Printf("ERROR: sales input record %v\n", trx)
			return
		}

		trx.Quarter = impl.getQuarterFromMonth(int(trx.SaleDate.Month()))

		fmt.Printf("sales record: %+v\n", trx)

		err = impl.PsqlDb.InsertSaleTransaction(trx)
		if err != nil {
			log.Println("sales insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
			return
		}
	}

	logger.FileLogger.Info().Msg("HandleDataInputFile : successful")
}

func (impl *IncentiveImpl) HandlePlansInputFile() {
	//fileName := "./doc/Input.xlsx"
	//
	//file, err := excelize.OpenFile(fileName)
	//if err != nil {
	//	log.Fatalf("could not open file, %v\n", err)
	//}
	//
	//rows, err := file.Rows("Plans")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println("total rows: ", rows.TotalRows())
	//
	////records := make([]*db.Transaction, 0)
	//impl.DBService.DeleteProductPlans()
	//
	//isFirstRow := true
	//index := 0
	//
	//for rows.Next() {
	//
	//	index++
	//
	//	if isFirstRow {
	//		rows.Columns()
	//		isFirstRow = false
	//		continue
	//	}
	//
	//	row, err := rows.Columns()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	trx := &db.ProductPlans{}
	//
	//	for i, colCell := range row {
	//
	//		switch i {
	//		case 0:
	//			trx.PlanName = strings.TrimSpace(colCell)
	//		case 1:
	//			num, err := strconv.ParseFloat(strings.TrimSpace(colCell), 64)
	//			if err != nil {
	//				log.Println("PlanPrice error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				return
	//			}
	//			trx.PlanPrice = num
	//		case 2:
	//			num, err := strconv.ParseFloat(strings.TrimSpace(colCell), 64)
	//			if err != nil {
	//				log.Println("SIPPlanPrice error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				return
	//			}
	//			trx.SIPRatePlan = num
	//		}
	//	}
	//
	//	err = impl.DBService.InsertProductPlans(trx)
	//	if err != nil {
	//		log.Println("productPlan insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
	//		return
	//	}
	//}

	if err := impl.PsqlDb.FetchPlans(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandlePlansInputFile : productPlan fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("HandlePlansInputFile : successful")
}

func (impl *IncentiveImpl) HandleTargetInputFile() {
	////fileName := "./doc/new_doc/q2_2021_monthly_data_modify.xlsx"
	//	fileName := "./doc/latest/latest_Monthly_Target_2021.xlsx"
	//
	//	file, err := excelize.OpenFile(fileName)
	//	if err != nil {
	//		log.Fatalf("could not open file, %v\n", err)
	//	}
	//
	//	rows, err := file.Rows("Target 2021")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	impl.DBService.DeleteStoreTargets()
	//
	//	index := 0
	//	//
	//	//for rows.Next() {
	//	//
	//	//	index++
	//	//
	//	//	if index ==4 {
	//	//		return
	//	//	}
	//	//
	//	//	cols, err := rows.Columns()
	//	//	if err != nil {
	//	//		fmt.Println(err)
	//	//		return
	//	//	}
	//	//
	//	//	for i, colCell := range cols {
	//	//		fmt.Print("data ", i,  " ", colCell, " ")
	//	//	}
	//	//
	//	//	fmt.Println()
	//	//}
	//
	//	region := ""
	//	product := ""
	//	m := make(map[string]map[string][12]*db.ProductGrpTarget)
	//LOOP:
	//	for rows.Next() {
	//
	//		index++
	//
	//		eachRowColumns, err := rows.Columns()
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//
	//		outlet := ""
	//
	//		for i, colCell := range eachRowColumns {
	//
	//			if strings.TrimSpace(colCell) == "" || strings.Contains(strings.ToUpper(strings.TrimSpace(colCell)), "TOTAL") ||
	//				strings.Contains(strings.ToUpper(strings.TrimSpace(colCell)), "GRAND TOTAL") {
	//				fmt.Println("blank", index)
	//				continue LOOP
	//			}
	//
	//			if impl.isProduct(strings.ToUpper(strings.TrimSpace(colCell))) {
	//				fmt.Println("product", index)
	//				product = strings.ToUpper(strings.TrimSpace(colCell))
	//				continue LOOP
	//			}
	//
	//			if impl.isRegion(strings.ToUpper(strings.TrimSpace(colCell))) {
	//				fmt.Println("region", index)
	//				region = strings.ToUpper(strings.TrimSpace(colCell))
	//				continue LOOP
	//			}
	//
	//			if i > 12 {
	//				return
	//			}
	//
	//			if i == 0 {
	//				outlet = strings.ToUpper(strings.TrimSpace(colCell))
	//				continue
	//			}
	//
	//			num, err := strconv.ParseFloat(strings.TrimSpace(colCell), 64)
	//			if err != nil {
	//				log.Println("conversion failed " + strings.TrimSpace(colCell) + " error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				return
	//			} else if num == 0 {
	//				log.Println("emtpy value found error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				return
	//			}
	//
	//			num = math.Round(num)
	//
	//			if m[region] == nil {
	//				m[region] = make(map[string][12]*db.ProductGrpTarget)
	//			}
	//
	//			_, ok := m[region][outlet]
	//			if !ok {
	//				m[region][outlet] = [12]*db.ProductGrpTarget{}
	//			}
	//
	//			v := m[region][outlet]
	//
	//			if v[i-1] == nil {
	//				v[i-1] = &db.ProductGrpTarget{}
	//			}
	//
	//			switch product {
	//			case commons.HBBCategory:
	//				v[i-1].HBB = int(num)
	//			case commons.BaqatiCategory:
	//				v[i-1].Baqati = int(num)
	//			case commons.HayyakCategory:
	//				v[i-1].Hayyak = int(num)
	//			case commons.UpgradeCategory:
	//				v[i-1].Upgrade = int(num)
	//			case "UPGRADE (BAQATI + HBB)":
	//				v[i-1].Upgrade = int(num)
	//			}
	//
	//			m[region][outlet] = v
	//
	//			//fmt.Println("break", region, outlet, v)
	//		}
	//	}
	//
	//	for regKey, regMap := range m {
	//		for outletKey, outletMap := range regMap {
	//			for month, rcrd := range outletMap {
	//				//fmt.Println("debug: ", regKey, outletKey, month+1, rcrd.HBB, rcrd.Hayyak, rcrd.Baqati, rcrd.Upgrade)
	//				err = impl.DBService.InsertStoreTargets(regKey, outletKey, month+1, rcrd)
	//				if err != nil {
	//					log.Println("Target insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//			}
	//		}
	//	}

	if err := impl.PsqlDb.FetchSaleTarget(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleTargetInputFile : sale target fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("HandleTargetInputFile : successful")
}

func (impl *IncentiveImpl) HandleScaleOrderInputFile() {
	//fileName := "./doc/Calculation logic.xlsx"
	//
	//file, err := excelize.OpenFile(fileName)
	//if err != nil {
	//	log.Fatalf("could not open file, %v\n", err)
	//}
	//
	//hbb, err := file.GetCellValue("Scale Order:", "A18")
	//if err != nil {
	//	return
	//}
	//
	//if hbb != "HBB" {
	//	return
	//}
	//
	//hbb, err := file.GetCellValue("Scale Order:", "B10")
	//if err != nil {
	//	return
	//}
	//
	//impl.DBService.InsertValueOfPoints(hbb, 75)

	if err := impl.PsqlDb.FetchValueOfPoint(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleScaleOrderInputFile : fetch vop failed")
		return
	}

	logger.FileLogger.Info().Msg("HandleScaleOrderInputFile : successful")
}

func (impl *IncentiveImpl) HandleMaxPayInputFile() {
	// fileName := "./doc/Calculation logic.xlsx"

	// file, err := excelize.OpenFile(fileName)
	// if err != nil {
	// 	log.Fatalf("could not open file, %v\n", err)
	// }

	// hbb, err := file.GetCellValue("Scale Order:", "A18")
	// if err != nil {
	// 	return
	// }

	// if hbb != "HBB" {
	// 	return
	// }

	// hbb, err := file.GetCellValue("Scale Order:", "B10")
	// if err != nil {
	// 	return
	// }

	// impl.DBService.InsertValueOfPoints(hbb, 75)

	if err := impl.PsqlDb.FetchMaxPayout(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleMaxPayInputFile : fetch maxPayout failed")
		return
	}

	logger.FileLogger.Info().Msg("HandleMaxPayInputFile : successful")
}

func (impl *IncentiveImpl) HandleHierarchyFile() {

	//fileName := "./doc/sales_hierarchy.xlsx"
	fileName := "./doc/latest/latest_Sales_Hierachies.xlsx"

	file, err := excelize.OpenFile(fileName)
	if err != nil {
		logger.FileLogger.Error().Err(err).Str("filename", fileName).Msg("HandleHierarchyFile : open file error")
	}

	rows, err := file.Rows("Sheet1 (2)")
	if err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleHierarchyFile : get rows error")
		return
	}

	if err := impl.PsqlDb.DeleteHierarchy(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("HandleHierarchyFile : delete hierarchy error")
		return
	}

	isFirstRow := true
	index := 0

	usrs := make([]*psqldb.UserHierarchy, 0)

LOOP:
	for rows.Next() {

		index++

		if isFirstRow {
			rows.Columns()
			isFirstRow = false
			continue
		}

		usr := &psqldb.UserHierarchy{}

		eachRowColumns, err := rows.Columns()
		if err != nil {
			logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : get columns error")
			return
		}

		for i, colCell := range eachRowColumns {

			switch i {
			case 0:
				if strings.TrimSpace(colCell) == "" {
					continue LOOP
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : hierarchy id parsing error")
					return
				}
				usr.EmployeeID = id
			case 1:

				if strings.TrimSpace(colCell) == "" {
					continue
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : outletManager id parsing error")
					return
				}
				usr.OutletManager = id
			case 2:

				if strings.TrimSpace(colCell) == "" {
					continue
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : manager id parsing error")
					return
				}
				usr.Manager = id
			case 3:

				if strings.TrimSpace(colCell) == "" {
					continue
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : seniorManager id parsing error")
					return
				}
				usr.SeniorManager = id
			case 4:

				if strings.TrimSpace(colCell) == "" {
					continue
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : generalManager id parsing error")
					return
				}
				usr.GeneralManager = id
			case 5:

				if strings.TrimSpace(colCell) == "" {
					continue
				}

				id, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : vicePresident id parsing error")
					return
				}
				usr.VicePresident = id
			}
		}

		if usr.EmployeeID != 0 {
			usrs = append(usrs, usr)
		}

		// insert in batches of 50
		if len(usrs)%50 == 0 {

			errs := impl.PsqlDb.InsertHierarchy(usrs...)
			if errs != nil {
				logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : hierarchy insertionError")
				return
			}

			usrs = usrs[:0]
		}
	}

	errs := impl.PsqlDb.InsertHierarchy(usrs...)
	if errs != nil {
		logger.FileLogger.Error().Err(err).Str("index", strconv.Itoa(index)).Msg("HandleHierarchyFile : hierarchy last insertionError")
		return
	}

	logger.FileLogger.Info().Msg("HandleHierarchyFile : successful")
}
