package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

// DB connection
//var dbConn *sql.DB

type PsqlOperation interface {
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	PingContext(ctx context.Context) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DBImplementation struct {
	Operation PsqlOperation
}

//NewDatabaseSetup start a new db connection
func NewDatabaseSetup(conf *config.Database) *DBImplementation {

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.UserName, conf.Password, conf.Name)

	utils.LogObj.Info("Connecting to : " + psqlConf)

	var err error
	var dbConn PsqlOperation
	dbConn, err = sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		utils.LogObj.Fatal("psql server error, err: " + err.Error())
	}

	dbConn.SetMaxIdleConns(conf.MaxIdelConn)
	dbConn.SetMaxOpenConns(conf.MaxOpenConn)
	//dbConn.SetConnMaxLifetime(24 * time.Hour)

	ctx := context.Background()

	err = dbConn.PingContext(ctx)
	if err != nil {
		utils.LogObj.Fatal("ping psql server error, err: " + err.Error())
	}

	return &DBImplementation{Operation: dbConn}
}

// InsertUserMasterRecord reads data from user master input file and inserts into table
func (dbConn *DBImplementation) InsertUserMasterRecord(emp *Employee) error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	query := "insert into User_Master_Table(created, modified, emp_id, name, usr_position, region, outlet, remark) values($1,$2,$3,$4,$5,$6,$7,$8)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), emp.Id, emp.Name, emp.Position, emp.Region, emp.Outlet, emp.Remark)
	if err != nil {
		return err
	}

	return err
}

// InsertSaleTransaction reads data from input file and inserts into table
func (dbConn *DBImplementation) InsertSaleTransaction(trx *Transaction) error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	query := "insert into Sale_Transaction_Table(created, modified, region, outlet, emp_id, category_type, prev_product, curr_product, quantity, category, product_group, sales_date, quarter) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,&12,&13)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), trx.Reqion, trx.Outlet, trx.Id, trx.CategoryType, trx.PreviousProduct,
		trx.CurrentProduct, trx.Quantity, trx.Category, trx.ProductGrp, trx.SaleDate, trx.Quarter)
	if err != nil {
		return err
	}

	return err
}

// DeleteProductPlans removes product plan data from table
func (dbConn *DBImplementation) DeleteProductPlans() error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := dbConn.Operation.ExecContext(ctx, "DELETE FROM Product_Plans_Table")
	if err != nil {
		return err
	}

	return nil
}

// InsertProductPlans reads data from plan input file and inserts into table
func (dbConn *DBImplementation) InsertProductPlans(plan *ProductPlans) error {

	// ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// dbConn.Operation.ExecContext(ctx, "DELETE FROM Product_Plans_Table")

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	query := "insert into Product_Plans_Table(created, modified, plan_name, plan_price, sip_plan_rate) values($1,$2,$3,$4,$5)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), plan.PlanName, plan.PlanPrice, plan.SIPRatePlan)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStoreTargets removes store target data from table
func (dbConn *DBImplementation) DeleteStoreTargets() error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := dbConn.Operation.ExecContext(ctx, "DELETE FROM Sales_Target_Table")
	if err != nil {
		return err
	}

	return nil
}

// InsertStoreTargets reads data from target input file and inserts into table
func (dbConn *DBImplementation) InsertStoreTargets(region, outlet string, month int, trgt *ProductGrpTarget) error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	query := "insert into Sales_Target_Table(created, modified, region, outlet, months, Baqati, HBB, Hayyak, Upgrade) values($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), region, outlet, month, trgt.Baqati, trgt.HBB, trgt.Hayyak, trgt.Upgrade)
	if err != nil {
		return err
	}

	return nil
}

// InsertValueOfPoints reads data from ScakeOrder input file and inserts into table
func (dbConn *DBImplementation) InsertValueOfPoints(category string, percent int, vop float64) (int64, error) {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	query := "insert into Scale_Order_Table(category, percentage, vop) values($1,$2,$3)"

	rslt, err := dbConn.Operation.ExecContext(ctx, query, category, percent, vop)
	if err != nil {
		return -1, err
	}

	return rslt.LastInsertId()
}

// FetchSaleTarget loads the global variable PlanRates with the rate plan for each plan type
func (dbConn *DBImplementation) FetchPlans() error {

	rows, err := dbConn.Operation.Query("select id, plan_name, plan_price, sip_plan_rate from Product_Plans_Table")
	if err != nil {
		return err
	}

	var (
		id            int
		plan_name     string
		plan_price    float64
		sip_plan_rate float64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &plan_name, &plan_price, &sip_plan_rate); err != nil {
			return err
		}

		PlanRates.Mutex.Lock()
		if PlanRates.PlanRateMap == nil {
			PlanRates.PlanRateMap = make(map[int]*ProductPlans)
		}

		PlanRates.PlanRateMap[id] = &ProductPlans{
			PlanName:    plan_name,
			PlanPrice:   plan_price,
			SIPRatePlan: sip_plan_rate,
		}

		PlanRates.Mutex.Unlock()

	}

	return nil
}

// FetchSaleTarget loads the global variable SalesTrgt with the target for each outlet
func (dbConn *DBImplementation) FetchSaleTarget() error {

	SalesTrgt.SaleTargetMap = make(map[string]map[string][12]*ProductGrpTarget)

	rows, err := dbConn.Operation.Query("select region, outlet, months, Baqati, HBB, Hayyak, Upgrade from Sales_Target_Table")
	if err != nil {
		return err
	}

	var (
		region  string
		outlet  string
		month   int
		baqati  int
		hBB     int
		hayyak  int
		upgrade int
	)

	for rows.Next() {
		if err := rows.Scan(&region, &outlet, &month, &baqati, &hBB, &hayyak, &upgrade); err != nil {
			return err
		}

		fmt.Println(region, outlet, month, baqati, hBB, hayyak, upgrade)

		SalesTrgt.Mutex.Lock()

		if SalesTrgt.SaleTargetMap[region] == nil {
			fmt.Println("inside2")
			SalesTrgt.SaleTargetMap[region] = make(map[string][12]*ProductGrpTarget)
		}

		array, ok := SalesTrgt.SaleTargetMap[region][outlet]
		if !ok {
			var temp [12]*ProductGrpTarget
			SalesTrgt.SaleTargetMap[region][outlet] = temp
		}

		array[month-1] = &ProductGrpTarget{
			Baqati:  baqati,
			HBB:     hBB,
			Hayyak:  hayyak,
			Upgrade: upgrade,
		}

		SalesTrgt.SaleTargetMap[region][outlet] = array

		SalesTrgt.Mutex.Unlock()

	}

	fmt.Printf("map: %+v /n", SalesTrgt.SaleTargetMap)

	return nil
}

// FetchValueOfPoint loads the global variable ValOfPoint with the value of point for each product type
func (dbConn *DBImplementation) FetchValueOfPoint() error {

	// rows, err := dbConn.Operation.Query("select group, top_percentage, vop, extra_pay from Product_Plans_Table")
	// if err != nil {
	// 	return err
	// }

	// var (
	// 	group          string
	// 	top_percentage int
	// 	vop            float64
	// 	extra_pay      float64
	// )

	// for rows.Next() {
	// 	if err := rows.Scan(&group, &top_percentage, &vop, &extra_pay); err != nil {
	// 		return err
	// 	}

	// 	ValOfPoint.Mutex.Lock()
	// 	if ValOfPoint.VopMap == nil {
	// 		ValOfPoint.VopMap = make(map[string]map[int]*VOP)
	// 	}

	// 	if ValOfPoint.VopMap[group] == nil {
	// 		ValOfPoint.VopMap[group] = make(map[int]*VOP)
	// 	}

	// 	ValOfPoint.VopMap[group][top_percentage] = &VOP{
	// 		Point:    vop,
	// 		ExtraPay: extra_pay,
	// 	}

	// 	ValOfPoint.Mutex.Unlock()

	// }

	ValOfPoint.Mutex.Lock()

	ValOfPoint.VopMap = make(map[string]map[int]*VOP)
	ValOfPoint.VopMap["HBB"] = make(map[int]*VOP)
	ValOfPoint.VopMap["BAQATI"] = make(map[int]*VOP)
	ValOfPoint.VopMap["HAYYAK"] = make(map[int]*VOP)
	ValOfPoint.VopMap["UPGRADE"] = make(map[int]*VOP)

	ValOfPoint.VopMap["HBB"][75] = &VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HBB"][95] = &VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HBB"][110] = &VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HBB"][139] = &VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HBB"][140] = &VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HBB"][159] = &VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	ValOfPoint.VopMap["HBB"][179] = &VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	ValOfPoint.VopMap["HBB"][189] = &VOP{
		Point:    0,
		ExtraPay: 52.5,
	}
	ValOfPoint.VopMap["HBB"][199] = &VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	ValOfPoint.VopMap["HBB"][200] = &VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	ValOfPoint.VopMap["BAQATI"][75] = &VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["BAQATI"][95] = &VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["BAQATI"][110] = &VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["BAQATI"][139] = &VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["BAQATI"][140] = &VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["BAQATI"][159] = &VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	ValOfPoint.VopMap["BAQATI"][179] = &VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	ValOfPoint.VopMap["BAQATI"][189] = &VOP{
		Point:    0,
		ExtraPay: 53,
	}
	ValOfPoint.VopMap["BAQATI"][199] = &VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	ValOfPoint.VopMap["BAQATI"][200] = &VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	ValOfPoint.VopMap["HAYYAK"][75] = &VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HAYYAK"][95] = &VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HAYYAK"][110] = &VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HAYYAK"][139] = &VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HAYYAK"][140] = &VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["HAYYAK"][159] = &VOP{
		Point:    0,
		ExtraPay: 11,
	}
	ValOfPoint.VopMap["HAYYAK"][179] = &VOP{
		Point:    0,
		ExtraPay: 22,
	}
	ValOfPoint.VopMap["HAYYAK"][189] = &VOP{
		Point:    0,
		ExtraPay: 33,
	}
	ValOfPoint.VopMap["HAYYAK"][199] = &VOP{
		Point:    0,
		ExtraPay: 44,
	}
	ValOfPoint.VopMap["HAYYAK"][200] = &VOP{
		Point:    0,
		ExtraPay: 55,
	}

	ValOfPoint.VopMap["UPGRADE"][75] = &VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["UPGRADE"][95] = &VOP{
		Point:    0.1,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["UPGRADE"][110] = &VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["UPGRADE"][139] = &VOP{
		Point:    0.3,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["UPGRADE"][140] = &VOP{
		Point:    0.5,
		ExtraPay: 0,
	}
	ValOfPoint.VopMap["UPGRADE"][159] = &VOP{
		Point:    0,
		ExtraPay: 4,
	}
	ValOfPoint.VopMap["UPGRADE"][179] = &VOP{
		Point:    0,
		ExtraPay: 8,
	}
	ValOfPoint.VopMap["UPGRADE"][189] = &VOP{
		Point:    0,
		ExtraPay: 12,
	}
	ValOfPoint.VopMap["UPGRADE"][199] = &VOP{
		Point:    0,
		ExtraPay: 16,
	}
	ValOfPoint.VopMap["UPGRADE"][200] = &VOP{
		Point:    0,
		ExtraPay: 20,
	}

	ValOfPoint.Mutex.Unlock()

	return nil
}

// FetchMaxPayout loads the gloabal variable MaxPay with the maximum payout limit
func (dbConn *DBImplementation) FetchMaxPayout() error {

	// rows, err := dbConn.Operation.Query("select group, percentage, max_pay from Max_Payout_Table")
	// if err != nil {
	// 	return err
	// }

	// var (
	// 	group      string
	// 	percentage int
	// 	max_pay    int
	// )

	// for rows.Next() {
	// 	if err := rows.Scan(&group, &percentage, &max_pay); err != nil {
	// 		return err
	// 	}

	// 	MaxPay.Mutex.Lock()
	// 	if MaxPay.Payoutmap == nil {
	// 		MaxPay.Payoutmap = make(map[string]float64)
	// 	}

	// 	MaxPay.Payoutmap[group] = float64(percentage) * float64(max_pay) / 100

	// 	MaxPay.Mutex.Unlock()

	// }

	MaxPay.Mutex.Lock()
	MaxPay.Payoutmap = make(map[string]float64)
	MaxPay.Payoutmap["HBB"] = 350
	MaxPay.Payoutmap["BAQATI"] = 350
	MaxPay.Payoutmap["HAYYAK"] = 220
	MaxPay.Payoutmap["UPGRADE"] = 80
	MaxPay.Mutex.Unlock()

	return nil
}

// FetchTotalNonManagerialEmployees returns all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchTotalNonManagerialEmployees(region, outlet string) (int, error) {

	// rows, err := dbConn.Operation.Query("select count(*) from User_Master_Table where region= $1 and outlet= $2 and position!= $3", region, outlet, commons.OutletManager)
	rows, err := dbConn.Operation.Query("select count(*) from User_Master_Table where region= $1 and outlet= $2 and usr_position!='BACK OFFICE'", region, outlet)
	if err != nil {
		return -1, err
	}

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return -1, err
		}
	}

	return count, nil
}

// FetchUserSales returns all sale transaction of a given user
func (dbConn *DBImplementation) FetchUserSales(clientID, month, year, quarter int, isQuarterReport bool) (*SaleTransaction, error) {

	var (
		rows *sql.Rows
		err  error
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.Query("select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and b.quarter= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, quarter, year)
	} else {
		rows, err = dbConn.Operation.Query("select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and EXTRACT('month' from b.sales_date)= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, month, year)
	}

	if err != nil {
		return nil, err
	}

	var (
		position      string
		outlet        string
		region        string
		prev_product  string
		curr_product  string
		quantity      int
		product_group string
		category      string
	)

	trx := &SaleTransaction{}
	trx.SaleRecords = make(map[string][]*SaleRecord)

	for rows.Next() {
		if err := rows.Scan(&position, &outlet, &region, &prev_product, &curr_product, &quantity, &product_group, &category); err != nil {
			return trx, err
		}

		trx.EmpID = clientID
		trx.Position = position
		trx.Outlet = outlet
		trx.Region = region

		rcrd := &SaleRecord{
			Quantity:        quantity,
			Category:        category,
			PreviousProduct: prev_product,
			CurrentProduct:  curr_product,
		}

		if trx.SaleRecords[product_group] == nil {
			trx.SaleRecords[product_group] = make([]*SaleRecord, 0)
		}

		trx.SaleRecords[product_group] = append(trx.SaleRecords[product_group], rcrd)
	}

	return trx, nil
}
