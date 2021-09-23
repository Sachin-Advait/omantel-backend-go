package crons

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	excelize "github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
)

type DBService interface {
	InsertUserMasterRecord(emp *db.Employee) error
	InsertSaleTransaction(trx *db.Transaction) error
	InsertProductPlans(plan *db.ProductPlans) error
	InsertStoreTargets(region, outlet string, trgt *db.ProductGrpTarget) error
	InsertValueOfPoints(category string, percent int, vop float64) (int64, error)

	FetchSaleTarget() error
	FetchPlans() error
	FetchValueOfPoint() error
	FetchMaxPayout() error

	DeleteStoreTargets() error
	DeleteProductPlans() error
}

type CronImpl struct {
	DBService DBService
}

func NewCronServices(conf *config.Database) *CronImpl {
	dbConn := db.NewDatabaseSetup(conf)

	return &CronImpl{DBService: dbConn}
}

func (impl *CronImpl) HandleUserMasterInputFile() {
	fmt.Println("input")
	fileName := "./doc/Input.xlsx"

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
}

func (impl *CronImpl) HandleDataInputFile() {
	fmt.Println("input")
	fileName := "./doc/Input.xlsx"

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

		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			return
		}

		trx := &db.Transaction{}

		//LOOP:
		for i, colCell := range row {

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
				// trx.PreviousProduct = strings.ToUpper(strings.TrimSpace(colCell))
				trx.PreviousProduct = strings.TrimSpace(colCell)
			case 5:
				// trx.CurrentProduct = strings.ToUpper(strings.TrimSpace(colCell))
				trx.CurrentProduct = strings.TrimSpace(colCell)
			case 6:
				num, err := strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					log.Println("Quantity conversion error: " + err.Error() + " index: " + strconv.Itoa(index))
					//break LOOP
					return
				}
				trx.Quantity = num
			case 7:
				trx.Category = strings.ToUpper(strings.TrimSpace(colCell))
			case 8:
				trx.ProductGrp = strings.ToUpper(strings.TrimSpace(colCell))
			}
		}

		if trx.Reqion == "" || trx.Outlet == "" || trx.CategoryType == "" || trx.CurrentProduct == "" || trx.Category == "" || trx.ProductGrp == "" {
			log.Printf("ERROR: sales input record %v\n", trx)
			return
		}

		err = impl.DBService.InsertSaleTransaction(trx)
		if err != nil {
			log.Println("sales insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
			return
		}
	}
}

func (impl *CronImpl) HandlePlansInputFile() {
	// fileName := "./doc/Input.xlsx"

	// file, err := excelize.OpenFile(fileName)
	// if err != nil {
	// 	log.Fatalf("could not open file, %v\n", err)
	// }

	// rows, err := file.Rows("Plans")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// //records := make([]*db.Transaction, 0)
	// impl.DBService.DeleteProductPlans()

	// isFirstRow := true
	// index := 0

	// for rows.Next() {

	// 	index++

	// 	if isFirstRow {
	// 		rows.Columns()
	// 		isFirstRow = false
	// 		continue
	// 	}

	// 	row, err := rows.Columns()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	trx := &db.ProductPlans{}

	// 	for i, colCell := range row {

	// 		switch i {
	// 		case 0:
	// 			trx.PlanName = strings.TrimSpace(colCell)
	// 		case 1:
	// 			num, err := strconv.ParseFloat(strings.TrimSpace(colCell), 64)
	// 			if err != nil {
	// 				log.Println("PlanPrice error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 				return
	// 			}
	// 			trx.PlanPrice = num
	// 		case 2:
	// 			num, err := strconv.ParseFloat(strings.TrimSpace(colCell), 64)
	// 			if err != nil {
	// 				log.Println("SIPPlanPrice error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 				return
	// 			}
	// 			trx.SIPRatePlan = num
	// 		}
	// 	}

	// 	err = impl.DBService.InsertProductPlans(trx)
	// 	if err != nil {
	// 		log.Println("productPlan insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
	// 		return
	// 	}
	// }

	if err := impl.DBService.FetchPlans(); err != nil {
		log.Println("productPlan fetching failed: " + err.Error())
	}
}

func (impl *CronImpl) HandleTargetInputFile() {
	// 	fileName := "./doc/Input.xlsx"

	// 	file, err := excelize.OpenFile(fileName)
	// 	if err != nil {
	// 		log.Fatalf("could not open file, %v\n", err)
	// 	}

	// 	rows, err := file.Rows("Target")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	//records := make([]*db.Transaction, 0)

	// 	impl.DBService.DeleteStoreTargets()

	// 	isFirstRow := true
	// 	index := 0

	// 	region := ""
	// LOOP:
	// 	for rows.Next() {

	// 		index++

	// 		if isFirstRow {
	// 			rows.Columns()
	// 			isFirstRow = false
	// 			continue
	// 		}

	// 		row, err := rows.Columns()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}

	// 		trx := &db.ProductGrpTarget{}
	// 		outlet := ""

	// 		for i, colCell := range row {

	// 			if impl.isRegion(strings.ToUpper(strings.TrimSpace(colCell))) {
	// 				region = strings.ToUpper(strings.TrimSpace(colCell))
	// 				continue LOOP
	// 			}

	// 			switch i {
	// 			case 0:
	// 				outlet = strings.ToUpper(strings.TrimSpace(colCell))
	// 			case 1:
	// 				num, err := strconv.Atoi(strings.TrimSpace(colCell))
	// 				if err != nil {
	// 					log.Println("Baqati target error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 					return
	// 				}
	// 				trx.Baqati = num
	// 			case 2:
	// 				num, err := strconv.Atoi(strings.TrimSpace(colCell))
	// 				if err != nil {
	// 					log.Println("HBB target error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 					return
	// 				}
	// 				trx.HBB = num
	// 			case 3:
	// 				num, err := strconv.Atoi(strings.TrimSpace(colCell))
	// 				if err != nil {
	// 					log.Println("Hayyak target error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 					return
	// 				}
	// 				trx.Hayyak = num
	// 			case 4:
	// 				num, err := strconv.Atoi(strings.TrimSpace(colCell))
	// 				if err != nil {
	// 					log.Println("Upgrade target error: " + err.Error() + " index: " + strconv.Itoa(index))
	// 					return
	// 				}
	// 				trx.Upgrade = num
	// 			}
	// 		}

	// 		if trx.HBB == 0 || trx.Baqati == 0 || trx.Hayyak == 0 || trx.Upgrade == 0 {
	// 			continue LOOP
	// 		}

	// 		err = impl.DBService.InsertStoreTargets(region, outlet, trx)
	// 		if err != nil {
	// 			log.Println("Target insertionError: " + err.Error() + " index: " + strconv.Itoa(index))
	// 			return
	// 		}
	// 	}

	if err := impl.DBService.FetchSaleTarget(); err != nil {
		log.Println("sale target fetching failed: " + err.Error())
	}
}

func (impl *CronImpl) HandleScaleOrderInputFile() {
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

	impl.DBService.FetchValueOfPoint()
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
}

func (impl *CronImpl) isRegion(region string) bool {

	if region == "BATINAH NORTH" || region == "BATINAH SOUTH" || region == "DAKHLIA & DAHIRAH " || region == "DHOFAR" ||
		region == "MUSCAT NORTH" || region == "MUSCAT SOUTH" || region == "SHARQIYA" || region == "GRAND TOTAL" {
		return true
	}

	return false
}
