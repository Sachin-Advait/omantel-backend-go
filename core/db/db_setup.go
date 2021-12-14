package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

// DB connection
//var dbConn *sql.DB

type PsqlOperation interface {
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	Begin() (*sql.Tx, error)
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
	dbConn.SetConnMaxIdleTime(30)
	dbConn.SetConnMaxLifetime(90)
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
	query := "insert into User_Master_Table(created, modified, emp_id, name, usr_position, region, outlet, mobile_number ,remark) values($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), emp.Id, emp.Name, emp.Position, emp.Region, emp.Outlet, emp.MobileNos, emp.Remark)
	if err != nil {
		return err
	}

	return err
}

// InsertSaleTransaction reads data from input file and inserts into table
func (dbConn *DBImplementation) InsertSaleTransaction(trx *Transaction) error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	query := "insert into Sale_Transaction_Table(created, modified, region, outlet, emp_id, category_type, prev_product, curr_product, quantity, category, product_group, sales_date, quarter) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)"

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

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

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

		SalesTrgt.Mutex.Lock()

		if SalesTrgt.SaleTargetMap[region] == nil {
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

func (dbConn *DBImplementation) FetchPromotersPresentYearSales(clientID int) (map[int]*SaleTransaction, map[int]*SaleTransaction, error) {

	rows, err := dbConn.Operation.Query("select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category, b.sales_date from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and EXTRACT('year' from b.sales_date)= $2", clientID, time.Now().Year())
	if err != nil {
		return nil, nil, err
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
		sale_date     time.Time
	)

	monthly := make(map[int]*SaleTransaction)
	quarterly := make(map[int]*SaleTransaction)

	for rows.Next() {
		if err := rows.Scan(&position, &outlet, &region, &prev_product, &curr_product, &quantity, &product_group, &category, &sale_date); err != nil {
			return nil, nil, err
		}

		rcrd := &SaleRecord{
			Quantity:        quantity,
			Category:        category,
			PreviousProduct: prev_product,
			CurrentProduct:  curr_product,
			SaleDate:        sale_date,
		}

		v, ok := monthly[int(sale_date.Month())]
		if !ok {
			trx := &SaleTransaction{}
			trx.EmpID = clientID
			trx.Position = position
			trx.Outlet = outlet
			trx.Region = region
			trx.SaleRecords = map[string][]*SaleRecord{
				product_group: {rcrd},
			}

			monthly[int(sale_date.Month())] = trx
		} else {

			if v.SaleRecords[product_group] == nil {
				v.SaleRecords[product_group] = make([]*SaleRecord, 0)
			}

			v.SaleRecords[product_group] = append(v.SaleRecords[product_group], rcrd)

			monthly[int(sale_date.Month())] = v
		}

		v, ok = quarterly[getQuarterFromMonth(int(sale_date.Month()))]
		if !ok {
			trx := &SaleTransaction{}
			trx.EmpID = clientID
			trx.Position = position
			trx.Outlet = outlet
			trx.Region = region
			trx.SaleRecords = map[string][]*SaleRecord{
				product_group: {rcrd},
			}

			quarterly[getQuarterFromMonth(int(sale_date.Month()))] = trx
		} else {

			if v.SaleRecords[product_group] == nil {
				v.SaleRecords[product_group] = make([]*SaleRecord, 0)
			}

			v.SaleRecords[product_group] = append(v.SaleRecords[product_group], rcrd)

			quarterly[getQuarterFromMonth(int(sale_date.Month()))] = v
		}
	}

	return monthly, quarterly, nil
}

// FetchTotalNonManagerialEmployees returns all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchEmployeeContactNumber(clientID int) (int, *commons.GenericError) {

	rows, err := dbConn.Operation.Query("select mobile_number from User_Master_Table where emp_id= $1", clientID)
	if err != nil {
		return -1, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	var (
		count         int
		contactNumber string
	)

	for rows.Next() {

		if count == 1 {
			return -1, &commons.GenericError{Code: 1001, ErrorMessage: "found more than on record"}
		}

		if err := rows.Scan(&contactNumber); err != nil {
			return -1, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}

		count++
	}

	if contactNumber == "" {
		return -1, &commons.GenericError{Code: 1003, ErrorMessage: "contact number not found"}
	}

	mobNum, err := strconv.Atoi(contactNumber)
	if err != nil {
		return -1, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	return mobNum, nil
}

func (dbConn *DBImplementation) InsertStaffOTP(otpModel *OtpModel) *commons.GenericError {
	_, err := dbConn.Operation.Query("DELETE FROM user_authentication_table WHERE emp_id= $1", otpModel.StaffID)
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "delete " + err.Error()}
	}

	_, err = dbConn.Operation.Query("INSERT INTO user_authentication_table(created, emp_id, expirytime, modified, otp, retry) VALUES ($1,$2,$3,$4,$5,$6);",
		otpModel.CreatedTime,
		otpModel.StaffID,
		otpModel.ExpiryTime,
		time.Now().Local(),
		otpModel.OTP,
		otpModel.Retry,
	)

	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "insert " + err.Error()}
	}

	return nil
}

func (dbConn *DBImplementation) FetchOTPValidity(clientID int) (*OtpModel, *commons.GenericError) {

	rows, err := dbConn.Operation.Query("select a.name, a.usr_position, b.otp, b.retry, b.expiryTime from User_Master_Table a inner join user_authentication_table b on a.emp_id=b.emp_id and a.emp_id= $1", clientID)
	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	var count int

	otModel := &OtpModel{}

	for rows.Next() {

		if count == 1 {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: "found more than on record"}
		}

		if err := rows.Scan(&otModel.StaffName, &otModel.StaffPosition, &otModel.OTP, &otModel.Retry, &otModel.ExpiryTime); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}

		count++
	}

	return otModel, nil
}

// FetchAllGivenPositionEmpIdAndTheirManager returns an array of all employees of the given position along with their managers, returns empty array if no records found
func (dbConn *DBImplementation) FetchAllGivenPositionEmpIdAndTheirManager(position string) ([]*UserLead, *commons.GenericError) {

	if !(position == commons.VicePresident || position == commons.GeneralManager || position == commons.SeniorManager || position == commons.Manager || position == commons.OutletManager || position == commons.Promoter) {
		return nil, &commons.GenericError{Code: 1002, ErrorMessage: "incorrect position req- " + position}
	}

	var query string

	switch position {
	case commons.VicePresident:
		query = "select distinct vice_president, 0 from user_hierarchy_table"
	case commons.GeneralManager:
		query = "select distinct general_manager, vice_president from user_hierarchy_table"
	case commons.SeniorManager:
		query = "select distinct senior_manager, general_manager from user_hierarchy_table"
	case commons.Manager:
		query = "select distinct manager, senior_manager from user_hierarchy_table"
	case commons.OutletManager:
		query = "select distinct outlet_manager, manager from user_hierarchy_table"
	case commons.Promoter:
		query = "select distinct emp_id, outlet_manager from user_hierarchy_table"
	}

	rows, err := dbConn.Operation.Query(query)
	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: "fetch emp nd manager err-" + err.Error()}
	}

	var array = make([]*UserLead, 0)

	for rows.Next() {

		usr := &UserLead{}

		if err := rows.Scan(&usr.EmpID, &usr.ReportTo); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}

		array = append(array, usr)
	}

	return array, nil
}

// FetchAllEmployeesReportingToGivenManager returns all employees reporting to given manager, returns empty array if no records found
func (dbConn *DBImplementation) FetchAllEmployeesReportingToGivenManagerRole(managerPosition string, managerID int) ([]*ReportingStaffIDAndName, *commons.GenericError) {

	if !(managerPosition == "vice_president" || managerPosition == "general_manager" || managerPosition == "senior_manager" || managerPosition == "manager" || managerPosition == "outlet_manager") {
		return nil, &commons.GenericError{Code: 1002, ErrorMessage: "incorrect position manager req- " + managerPosition}
	}

	var query string

	switch managerPosition {
	case "vice_president":
		query = "select distinct a.general_manager, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.general_manager=b.emp_id where a.vice_president= $1"
	case "general_manager":
		query = "select distinct a.senior_manager, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.senior_manager=b.emp_id where a.general_manager= $1"
	case "senior_manager":
		query = "select distinct a.manager, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.manager=b.emp_id where a.senior_manager= $1"
	case "manager":
		query = "select distinct a.outlet_manager, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager= $1"
	case "outlet_manager":
		query = "select distinct a.emp_id, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.emp_id=b.emp_id where a.outlet_manager= $1"
	}

	rows, err := dbConn.Operation.Query(query, managerID)
	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	var (
		array = make([]*ReportingStaffIDAndName, 0)
	)

	for rows.Next() {

		usr := &ReportingStaffIDAndName{}

		if err := rows.Scan(&usr.ID, &usr.Name); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: "FetchAllEmployeesReportingToGivenManagerRole next row err: " + err.Error()}
		}

		array = append(array, usr)
	}

	return array, nil
}

// InsertEmployeeRecord inserts each users statistical daily in db, and is updated everyday once
func (dbConn *DBImplementation) InsertEmployeeRecord(usr *UserRecord) *commons.GenericError {

	// fmt.Printf("userRecord: %+v\n", *usr)

	//add where clause
	_, err := dbConn.Operation.Query("INSERT INTO staff_daily_report_table(emp_id, created, reporting_manager, isquarterly, months, quarter, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31);",
		usr.EmpId,
		usr.Created,
		usr.ReportingManager,
		usr.Isquarterly,
		usr.Month,
		usr.Quarter,
		usr.Hayyak_sales_count,
		usr.Hayyak_target,
		usr.Hayyak_achieved_percentage,
		usr.Hayyak_payout,
		usr.Hayyak_extra_payout,
		usr.Hbb_sales_count,
		usr.Hbb_target,
		usr.Hbb_achieved_percentage,
		usr.Hbb_payout,
		usr.Hbb_extra_payout,
		usr.Baqati_sales_count,
		usr.Baqati_target,
		usr.Baqati_achieved_percentage,
		usr.Baqati_payout,
		usr.Baqati_extra_payout,
		usr.Upgrade_sales_count,
		usr.Upgrade_target,
		usr.Uprade_achieved_percentage,
		usr.Upgrade_payout,
		usr.Upgrade_extra_payout,
		usr.Total_payout,
		usr.HayyakSaleCountPerPlanArray,
		usr.HBBSaleCountPerPlanArray,
		usr.BaqatiSaleCountPerPlanArray,
		usr.UpgradeSaleCountPerPlanArray,
	)

	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "insert emp rcrd err-" + err.Error()}
	}

	return nil
}

// InsertEmployeeMultipleRecords inserts each users multiple statistical data daily in db, and is updated everyday once
func (dbConn *DBImplementation) InsertEmployeeMultipleRecords(usr []*UserRecord) *commons.GenericError {

	if len(usr) == 0 {
		return nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}

	count := 0
	for _, w := range usr {
		valueStrings = append(valueStrings, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)", count+1, count+2, count+3, count+4, count+5, count+6, count+7, count+8, count+9, count+10, count+11, count+12, count+13, count+14, count+15, count+16, count+17, count+18, count+19, count+20, count+21, count+22, count+23, count+24, count+25, count+26, count+27, count+28, count+29, count+30, count+31))

		valueArgs = append(valueArgs, w.EmpId)
		valueArgs = append(valueArgs, w.Created)
		valueArgs = append(valueArgs, w.ReportingManager)
		valueArgs = append(valueArgs, w.Isquarterly)
		valueArgs = append(valueArgs, w.Month)
		valueArgs = append(valueArgs, w.Quarter)
		valueArgs = append(valueArgs, w.Hayyak_sales_count)
		valueArgs = append(valueArgs, w.Hayyak_target)
		valueArgs = append(valueArgs, w.Hayyak_achieved_percentage)
		valueArgs = append(valueArgs, w.Hayyak_payout)
		valueArgs = append(valueArgs, w.Hayyak_extra_payout)
		valueArgs = append(valueArgs, w.Hbb_sales_count)
		valueArgs = append(valueArgs, w.Hbb_target)
		valueArgs = append(valueArgs, w.Hbb_achieved_percentage)
		valueArgs = append(valueArgs, w.Hbb_payout)
		valueArgs = append(valueArgs, w.Hbb_extra_payout)
		valueArgs = append(valueArgs, w.Baqati_sales_count)
		valueArgs = append(valueArgs, w.Baqati_target)
		valueArgs = append(valueArgs, w.Baqati_achieved_percentage)
		valueArgs = append(valueArgs, w.Baqati_payout)
		valueArgs = append(valueArgs, w.Baqati_extra_payout)
		valueArgs = append(valueArgs, w.Upgrade_sales_count)
		valueArgs = append(valueArgs, w.Upgrade_target)
		valueArgs = append(valueArgs, w.Uprade_achieved_percentage)
		valueArgs = append(valueArgs, w.Upgrade_payout)
		valueArgs = append(valueArgs, w.Upgrade_extra_payout)
		valueArgs = append(valueArgs, w.Total_payout)
		valueArgs = append(valueArgs, w.HayyakSaleCountPerPlanArray)
		valueArgs = append(valueArgs, w.HBBSaleCountPerPlanArray)
		valueArgs = append(valueArgs, w.BaqatiSaleCountPerPlanArray)
		valueArgs = append(valueArgs, w.UpgradeSaleCountPerPlanArray)

		count += 31
	}

	smt := `INSERT INTO staff_daily_report_table(emp_id, created, reporting_manager, isquarterly, months, quarter, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan)
    VALUES %s ON CONFLICT (emp_id, months, quarter) DO NOTHING`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx, err := dbConn.Operation.Begin()
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "InsertEmployeeMultipleRecords db begin error " + err.Error()}
	}

	_, err = tx.Exec("Delete FROM staff_daily_report_table where emp_id=$1", usr[0].EmpId)
	if err != nil {
		tx.Rollback()
		return &commons.GenericError{Code: 1001, ErrorMessage: "InsertEmployeeMultipleRecords db delete error " + err.Error()}
	}

	if usr[0].EmpId == 34296 {
		for _, v := range usr {
			fmt.Printf("%+v\n", v)
		}
	}

	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return &commons.GenericError{Code: usr[0].EmpId, ErrorMessage: "InsertEmployeeMultipleRecords db insert error " + err.Error()}
	}

	if err := tx.Commit(); err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "InsertEmployeeMultipleRecords db commit error " + err.Error()}
	}

	return nil
}

func (dbConn *DBImplementation) InsertHierarchy(usr ...*UserHierarchy) *commons.GenericError {

	valueStrings := []string{}
	valueArgs := []interface{}{}

	count := 0
	for _, w := range usr {
		valueStrings = append(valueStrings, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v)", count+1, count+2, count+3, count+4, count+5, count+6))
		valueArgs = append(valueArgs, w.EmployeeID)
		valueArgs = append(valueArgs, w.OutletManager)
		valueArgs = append(valueArgs, w.Manager)
		valueArgs = append(valueArgs, w.SeniorManager)
		valueArgs = append(valueArgs, w.GeneralManager)
		valueArgs = append(valueArgs, w.VicePresident)

		count += 6
	}

	smt := `INSERT INTO user_hierarchy_table(emp_id, outlet_manager, manager, senior_manager, general_manager, vice_president)
    VALUES %s ON CONFLICT (emp_id) DO NOTHING`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx, err := dbConn.Operation.Begin()
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "db begin hierarchy error " + err.Error()}
	}

	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return &commons.GenericError{Code: 1001, ErrorMessage: "db insert hierarchy error " + err.Error()}
	}

	if err := tx.Commit(); err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "db commit hierarchy error " + err.Error()}
	}

	return nil
}

func (dbConn *DBImplementation) FetchIndividualUserRecord(id int) ([]*UserRecord, *commons.GenericError) {

	rows, err := dbConn.Operation.Query("select reporting_manager, isquarterly, months, quarter, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from staff_daily_report_table where emp_id= $1", id)
	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	rcrds := make([]*UserRecord, 0)

	for rows.Next() {

		rcrd := &UserRecord{}

		if err := rows.Scan(
			&rcrd.ReportingManager,
			&rcrd.Isquarterly,
			&rcrd.Month,
			&rcrd.Quarter,
			&rcrd.Hayyak_sales_count,
			&rcrd.Hayyak_target,
			&rcrd.Hayyak_achieved_percentage,
			&rcrd.Hayyak_payout,
			&rcrd.Hayyak_extra_payout,
			&rcrd.Hbb_sales_count,
			&rcrd.Hbb_target,
			&rcrd.Hbb_achieved_percentage,
			&rcrd.Hbb_payout,
			&rcrd.Hbb_extra_payout,
			&rcrd.Baqati_sales_count,
			&rcrd.Baqati_target,
			&rcrd.Baqati_achieved_percentage,
			&rcrd.Baqati_payout,
			&rcrd.Baqati_extra_payout,
			&rcrd.Upgrade_sales_count,
			&rcrd.Upgrade_target,
			&rcrd.Uprade_achieved_percentage,
			&rcrd.Upgrade_payout,
			&rcrd.Upgrade_extra_payout,
			&rcrd.Total_payout,
			&rcrd.HayyakSaleCountPerPlanArray,
			&rcrd.HBBSaleCountPerPlanArray,
			&rcrd.BaqatiSaleCountPerPlanArray,
			&rcrd.UpgradeSaleCountPerPlanArray,
		); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}

		rcrds = append(rcrds, rcrd)
	}

	return rcrds, nil
}

func (dbConn *DBImplementation) FetchIndividualUserPayoutReport(id, month, quarter int, isQuarterly bool) (*UserRecord, *commons.GenericError) {

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.Query("select hayyak_payout, hayyak_extra_payout, hbb_payout, hbb_extra_payout, baqati_payout, baqati_extra_payout, upgrade_payout, upgrade_extra_payout, total_payout from staff_daily_report_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.Query("select hayyak_payout, hayyak_extra_payout, hbb_payout, hbb_extra_payout, baqati_payout, baqati_extra_payout, upgrade_payout, upgrade_extra_payout, total_payout from staff_daily_report_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	rcrd := &UserRecord{}

	for rows.Next() {

		if err := rows.Scan(
			&rcrd.Hayyak_payout,
			&rcrd.Hayyak_extra_payout,
			&rcrd.Hbb_payout,
			&rcrd.Hbb_extra_payout,
			&rcrd.Baqati_payout,
			&rcrd.Baqati_extra_payout,
			&rcrd.Upgrade_payout,
			&rcrd.Upgrade_extra_payout,
			&rcrd.Total_payout,
		); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserTargetReport(id, month, quarter int, isQuarterly bool) (*UserRecord, *commons.GenericError) {

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.Query("select hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hbb_sales_count, hbb_target, hbb_achieved_percentage, baqati_sales_count, baqati_target, baqati_achieved_percentage, upgrade_sales_count, upgrade_target, uprade_achieved_percentage from staff_daily_report_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.Query("select hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hbb_sales_count, hbb_target, hbb_achieved_percentage, baqati_sales_count, baqati_target, baqati_achieved_percentage, upgrade_sales_count, upgrade_target, uprade_achieved_percentage from staff_daily_report_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	rcrd := &UserRecord{}

	for rows.Next() {

		if err := rows.Scan(
			&rcrd.Hayyak_sales_count,
			&rcrd.Hayyak_target,
			&rcrd.Hayyak_achieved_percentage,
			&rcrd.Hbb_sales_count,
			&rcrd.Hbb_target,
			&rcrd.Hbb_achieved_percentage,
			&rcrd.Baqati_sales_count,
			&rcrd.Baqati_target,
			&rcrd.Baqati_achieved_percentage,
			&rcrd.Upgrade_sales_count,
			&rcrd.Upgrade_target,
			&rcrd.Uprade_achieved_percentage,
		); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserSaleCountPerPlanReport(id, month, quarter int, isQuarterly bool) (*UserRecord, *commons.GenericError) {

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.Query("select hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from staff_daily_report_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.Query("select hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from staff_daily_report_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	rcrd := &UserRecord{}

	for rows.Next() {

		if err := rows.Scan(
			&rcrd.HayyakSaleCountPerPlanArray,
			&rcrd.HBBSaleCountPerPlanArray,
			&rcrd.BaqatiSaleCountPerPlanArray,
			&rcrd.UpgradeSaleCountPerPlanArray,
		); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserOverAllReport(id, month, quarter int, isQuarterly bool) (*UserRecord, *commons.GenericError) {

	var (
		err  error
		rows *sql.Rows
	)

	if isQuarterly {
		rows, err = dbConn.Operation.Query("select reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout from staff_daily_report_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.Query("select reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout from staff_daily_report_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	rcrd := &UserRecord{}

	for rows.Next() {

		if err := rows.Scan(
			&rcrd.ReportingManager,
			&rcrd.Hayyak_sales_count,
			&rcrd.Hayyak_target,
			&rcrd.Hayyak_achieved_percentage,
			&rcrd.Hayyak_payout,
			&rcrd.Hayyak_extra_payout,
			&rcrd.Hbb_sales_count,
			&rcrd.Hbb_target,
			&rcrd.Hbb_achieved_percentage,
			&rcrd.Hbb_payout,
			&rcrd.Hbb_extra_payout,
			&rcrd.Baqati_sales_count,
			&rcrd.Baqati_target,
			&rcrd.Baqati_achieved_percentage,
			&rcrd.Baqati_payout,
			&rcrd.Baqati_extra_payout,
			&rcrd.Upgrade_sales_count,
			&rcrd.Upgrade_target,
			&rcrd.Uprade_achieved_percentage,
			&rcrd.Upgrade_payout,
			&rcrd.Upgrade_extra_payout,
			&rcrd.Total_payout,
		); err != nil {
			return nil, &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
		}
	}

	return rcrd, nil
}

func getQuarterFromMonth(month int) int {
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
