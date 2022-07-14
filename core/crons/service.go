package crons

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	excelize "github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
)

type DBService interface {
	InsertUserMasterRecord(emp *db.Employee) error
	InsertSaleTransaction(trx *db.Transaction) error
	InsertProductPlans(plan *db.ProductPlans) error
	InsertStoreTargets(region, outlet string, month int, trgt *db.ProductGrpTarget) error
	InsertValueOfPoints(category string, percent int, vop float64) (int64, error)
	InsertHierarchy(usr ...*db.UserHierarchy) *commons.GenericError

	FetchSaleTarget() error
	FetchPlans() error
	FetchValueOfPoint() error
	FetchMaxPayout() error

	DeleteStoreTargets() error
	DeleteProductPlans() error
	DeleteHierarchy() error
}

type CronImpl struct {
	DBService DBService
}

func NewCronServices(conf *config.Database) *CronImpl {
	dbConn := db.NewDatabaseSetup(conf)

	return &CronImpl{DBService: dbConn}
}

func (impl *CronImpl) HandleUserMasterInputFile() {

	fileName := "./latest/master.xlsx"

	file, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatalf("could not open file, %v\n", err)
	}

	rows, err := file.Rows("CSR & Promoter")
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

		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			return
		}

		emp := &db.Employee{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				num, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					log.Println("ID conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
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
			log.Printf("ERROR: employee input record %+v\n", emp)
			return
		}

		err = impl.DBService.InsertUserMasterRecord(emp)
		if err != nil {
			log.Println("insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
			return
		}
	}

	fmt.Println("HandleUserMasterInputFile successfull")
}

func (impl *CronImpl) HandleDataInputFile() {
	////fileName := "doc/new_doc/q2_2021_monthly_data_orig.xlsx"
	//fileName := "doc/latest/latest_Q2_2021_Monthly_data_updated.xlsx"
	//
	//file, err := excelize.OpenFile(fileName)
	//if err != nil {
	//	log.Fatalf("could not open file, %v\n", err)
	//}
	//
	//rows, err := file.Rows("Data Q2 2021")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	////records := make([]*db.Transaction, 0)
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
	//	eachRowColumns, err := rows.Columns()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	trx := &db.Transaction{}
	//
	//	//LOOP:
	//	for i, colCell := range eachRowColumns {
	//
	//		switch i {
	//
	//		case 0:
	//			trx.Reqion = strings.ToUpper(strings.TrimSpace(colCell))
	//		case 1:
	//			trx.Outlet = strings.ToUpper(strings.TrimSpace(colCell))
	//		case 2:
	//			num, err := strconv.Atoi(strings.TrimSpace(colCell))
	//			if err != nil {
	//				log.Println("ID conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				//break LOOP
	//				return
	//			}
	//			trx.Id = num
	//		case 3:
	//			trx.CategoryType = strings.ToUpper(strings.TrimSpace(colCell))
	//		case 4:
	//			t, err := time.Parse("01-02-06", strings.TrimSpace(colCell))
	//			if err != nil {
	//				log.Println("date conversion error: " + err.Error() + " index: " + strconv.Itoa(index) + " cell " + strings.TrimSpace(colCell))
	//				//break LOOP
	//				return
	//			}
	//			trx.SaleDate = t
	//			trx.Quarter = impl.getQuarterFromMonth(int(t.Month()))
	//		case 5:
	//			// trx.PreviousProduct = strings.ToUpper(strings.TrimSpace(colCell))
	//			trx.PreviousProduct = strings.TrimSpace(colCell)
	//		case 6:
	//			// trx.CurrentProduct = strings.ToUpper(strings.TrimSpace(colCell))
	//			trx.CurrentProduct = strings.TrimSpace(colCell)
	//		case 7:
	//			num, err := strconv.Atoi(strings.TrimSpace(colCell))
	//			if err != nil {
	//				log.Println("Quantity conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
	//				//break LOOP
	//				return
	//			}
	//			trx.Quantity = num
	//		case 8:
	//			trx.Category = strings.ToUpper(strings.TrimSpace(colCell))
	//		case 9:
	//			trx.ProductGrp = strings.ToUpper(strings.TrimSpace(colCell))
	//			// case 10:
	//			// 	t, err := time.Parse("", strings.TrimSpace(colCell))
	//			// 	if err != nil {
	//			// 		log.Println("Incorrect sales time error: " + err.Error() + " index: " + strconv.Itoa(index))
	//			// 		//break LOOP
	//			// 		return
	//			// 	}
	//
	//			// 	trx.SaleDate = t
	//			// 	trx.Quarter = impl.getQuarterFromMonth(int(t.Month()))
	//		}
	//	}
	//
	//	if trx.Reqion == "" || trx.Outlet == "" || trx.CategoryType == "" || trx.CurrentProduct == "" || trx.Category == "" || trx.ProductGrp == "" || trx.SaleDate.Year() == 1 {
	//		log.Printf("ERROR: sales input record %v\n", trx)
	//		return
	//	}
	//
	//	trx.Quarter = impl.getQuarterFromMonth(int(trx.SaleDate.Month()))
	//
	//	fmt.Printf("sales record: %+v\n", trx)
	//
	//	err = impl.DBService.InsertSaleTransaction(trx)
	//	if err != nil {
	//		log.Println("sales insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
	//		return
	//	}
	//}

	fmt.Println("HandleDataInputFile successfull")
}

func (impl *CronImpl) HandlePlansInputFile() {
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

	if err := impl.DBService.FetchPlans(); err != nil {
		log.Println("productPlan fetching failed: " + err.Error())
	}
	fmt.Println("HandlePlansInputFile successful")
}

func (impl *CronImpl) HandleTargetInputFile() {
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

	if err := impl.DBService.FetchSaleTarget(); err != nil {
		log.Println("sale target fetching failed: " + err.Error())
	}

	fmt.Println("HandleTargetInputFile successfull")
}

func (impl *CronImpl) HandleScaleOrderInputFile() {
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

	impl.DBService.FetchValueOfPoint()
	fmt.Println("HandleScaleOrderInputFile successful")
}

func (impl *CronImpl) HandleMaxPayInputFile() {
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

	impl.DBService.FetchMaxPayout()
	fmt.Println("HandleMaxPayInputFile successful")
}

func (impl *CronImpl) HandleHierarchyFile() {
	////fileName := "./doc/sales_hierarchy.xlsx"
	//	fileName := "./doc/latest/latest_Sales_Hierachies.xlsx"
	//
	//	file, err := excelize.OpenFile(fileName)
	//	if err != nil {
	//		log.Fatalf("could not open file, %v\n", err)
	//	}
	//
	//	rows, err := file.Rows("Sheet1 (2)")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	if err := impl.DBService.DeleteHierarchy(); err != nil {
	//		fmt.Println("delete hierarchy error: "+err.Error())
	//		return
	//	}
	//
	//	isFirstRow := true
	//	index := 0
	//
	//	usrs := make([]*db.UserHierarchy, 0)
	//
	//LOOP:
	//	for rows.Next() {
	//
	//		index++
	//
	//		if isFirstRow {
	//			rows.Columns()
	//			isFirstRow = false
	//			continue
	//		}
	//
	//		usr := &db.UserHierarchy{}
	//
	//		eachRowColumns, err := rows.Columns()
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//
	//		for i, colCell := range eachRowColumns {
	//
	//			switch i {
	//			case 0:
	//				if strings.TrimSpace(colCell) == "" {
	//					continue LOOP
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy id parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.EmployeeID = id
	//			case 1:
	//
	//				if strings.TrimSpace(colCell) == "" {
	//					continue
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy outletManager parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.OutletManager = id
	//			case 2:
	//
	//				if strings.TrimSpace(colCell) == "" {
	//					continue
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy manager parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.Manager = id
	//			case 3:
	//
	//				if strings.TrimSpace(colCell) == "" {
	//					continue
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy seniorManager parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.SeniorManager = id
	//			case 4:
	//
	//				if strings.TrimSpace(colCell) == "" {
	//					continue
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy generalManager parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.GeneralManager = id
	//			case 5:
	//
	//				if strings.TrimSpace(colCell) == "" {
	//					continue
	//				}
	//
	//				id, err := strconv.Atoi(strings.TrimSpace(colCell))
	//				if err != nil {
	//					log.Println("hierarchy vicePresident parsing error: " + err.Error() + " index: " + strconv.Itoa(index))
	//					return
	//				}
	//				usr.VicePresident = id
	//			}
	//		}
	//
	//		if usr.EmployeeID != 0 {
	//			fmt.Println(index, "id:", usr.EmployeeID)
	//			usrs = append(usrs, usr)
	//		}
	//
	//		// insert in batches of 50
	//		if len(usrs)%50 == 0 {
	//
	//			errs := impl.DBService.InsertHierarchy(usrs...)
	//			if errs != nil {
	//				log.Println("hierarchy insertionError: " + errs.ErrorMessage + " index: " + strconv.Itoa(index))
	//				return
	//			}
	//
	//			usrs = usrs[:0]
	//		}
	//	}
	//
	//	errs := impl.DBService.InsertHierarchy(usrs...)
	//	if errs != nil {
	//		log.Println("hierarchy last insertionError: " + errs.ErrorMessage + " index: " + strconv.Itoa(index))
	//		return
	//	}

	fmt.Println("HandleHierarchyFile successful")
}

func (impl *CronImpl) isRegion(region string) bool {

	if region == "BATINAH NORTH" || region == "BATINAH SOUTH" || region == "DAKHLIA & DAHIRAH" || region == "DHOFAR" ||
		region == "MUSCAT NORTH" || region == "MUSCAT SOUTH" || region == "SHARQIYA" || region == "GRAND TOTAL" {
		return true
	}

	return false
}

func (impl *CronImpl) getQuarterFromMonth(month int) int {
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

func (impl *CronImpl) isProduct(product string) bool {

	if product == commons.HBBCategory || product == commons.BaqatiCategory || product == commons.HayyakCategory || product == commons.UpgradeCategory {
		return true
	}

	if product == "UPGRADE (BAQATI + HBB)" {
		return true
	}

	return false
}
