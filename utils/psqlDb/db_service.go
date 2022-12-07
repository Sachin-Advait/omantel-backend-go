package psqldb

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

// FetchPlans loads the global variable PlanRates with the rate plan for each plan type
func (dbConn *DBImplementation) FetchPlans() *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchPlans start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchPlans done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select plan_code, plan_name, plan_price, sip_plan_rate from plans_table")
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchPlans Query: failed",
			Details:    "FetchPlans Query: " + err.Error(),
		}
	}

	var (
		id            string
		plan_name     string
		plan_price    float64
		sip_plan_rate float64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &plan_name, &plan_price, &sip_plan_rate); err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPlans Scan: failed",
				Details:    "FetchPlans Scan: " + err.Error(),
			}
		}

		models.PlanRates.Mutex.Lock()
		if models.PlanRates.PlanRateMap == nil {
			models.PlanRates.PlanRateMap = make(map[string]*models.ProductPlans)
		}

		models.PlanRates.PlanRateMap[id] = &models.ProductPlans{
			PlanName:    plan_name,
			PlanPrice:   plan_price,
			SIPRatePlan: sip_plan_rate,
		}

		models.PlanRates.Mutex.Unlock()

	}

	return nil
}

// FetchSaleTarget loads the global variable SalesTrgt with the target for each outlet
func (dbConn *DBImplementation) FetchSaleTarget() *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchSaleTarget start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchSaleTarget done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	models.SalesTrgt.SaleTargetMap = make(map[string]map[string][12]*models.ProductGrpTarget)

	rows, err := dbConn.Operation.QueryContext(ctx, "select region, outlet, month, baqati, hbb, hayyak, upgrade from targets_table")
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchSaleTarget Query: failed",
			Details:    "FetchSaleTarget Query: " + err.Error(),
		}
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
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchSaleTarget Scan: failed",
				Details:    "FetchSaleTarget Scan: " + err.Error(),
			}
		}

		models.SalesTrgt.Mutex.Lock()

		if models.SalesTrgt.SaleTargetMap[region] == nil {
			models.SalesTrgt.SaleTargetMap[region] = make(map[string][12]*models.ProductGrpTarget)
		}

		array, ok := models.SalesTrgt.SaleTargetMap[region][outlet]
		if !ok {
			var temp [12]*models.ProductGrpTarget
			models.SalesTrgt.SaleTargetMap[region][outlet] = temp
		}

		array[month-1] = &models.ProductGrpTarget{
			Baqati:  baqati,
			HBB:     hBB,
			Hayyak:  hayyak,
			Upgrade: upgrade,
		}

		models.SalesTrgt.SaleTargetMap[region][outlet] = array

		models.SalesTrgt.Mutex.Unlock()

	}

	return nil
}

// FetchValueOfPoint loads the global variable ValOfPoint with the value of point for each product type
func (dbConn *DBImplementation) FetchValueOfPoint() *errors.Error {

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

	models.ValOfPoint.Mutex.Lock()

	models.ValOfPoint.VopMap = make(map[string]map[int]*models.VOP)
	models.ValOfPoint.VopMap["HBB"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap["BAQATI"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap["HAYYAK"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap["UPGRADE"] = make(map[int]*models.VOP)

	models.ValOfPoint.VopMap["HBB"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HBB"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HBB"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HBB"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HBB"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HBB"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	models.ValOfPoint.VopMap["HBB"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	models.ValOfPoint.VopMap["HBB"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 52.5,
	}
	models.ValOfPoint.VopMap["HBB"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	models.ValOfPoint.VopMap["HBB"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	models.ValOfPoint.VopMap["BAQATI"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["BAQATI"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["BAQATI"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["BAQATI"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["BAQATI"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["BAQATI"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	models.ValOfPoint.VopMap["BAQATI"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	models.ValOfPoint.VopMap["BAQATI"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 53,
	}
	models.ValOfPoint.VopMap["BAQATI"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	models.ValOfPoint.VopMap["BAQATI"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	models.ValOfPoint.VopMap["HAYYAK"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HAYYAK"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HAYYAK"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HAYYAK"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HAYYAK"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["HAYYAK"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 11,
	}
	models.ValOfPoint.VopMap["HAYYAK"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 22,
	}
	models.ValOfPoint.VopMap["HAYYAK"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 33,
	}
	models.ValOfPoint.VopMap["HAYYAK"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 44,
	}
	models.ValOfPoint.VopMap["HAYYAK"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 55,
	}

	models.ValOfPoint.VopMap["UPGRADE"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["UPGRADE"][95] = &models.VOP{
		Point:    0.1,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["UPGRADE"][110] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["UPGRADE"][139] = &models.VOP{
		Point:    0.3,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["UPGRADE"][140] = &models.VOP{
		Point:    0.5,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap["UPGRADE"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 4,
	}
	models.ValOfPoint.VopMap["UPGRADE"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 8,
	}
	models.ValOfPoint.VopMap["UPGRADE"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 12,
	}
	models.ValOfPoint.VopMap["UPGRADE"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 16,
	}
	models.ValOfPoint.VopMap["UPGRADE"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 20,
	}

	models.ValOfPoint.Mutex.Unlock()

	return nil
}

// FetchMaxPayout loads the gloabal variable MaxPay with the maximum payout limit
func (dbConn *DBImplementation) FetchMaxPayout() *errors.Error {

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

	models.MaxPay.Mutex.Lock()
	models.MaxPay.Payoutmap = make(map[string]float64)
	models.MaxPay.Payoutmap["HBB"] = 350
	models.MaxPay.Payoutmap["BAQATI"] = 350
	models.MaxPay.Payoutmap["HAYYAK"] = 220
	models.MaxPay.Payoutmap["UPGRADE"] = 80
	models.MaxPay.Mutex.Unlock()

	return nil
}

////////////////

func (dbConn *DBImplementation) UpsertRegionRecord(dbReq *models.NewRegionRequest) *errors.Error {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertRegionRecord start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO regions_table(region_code, region_name, created, modified, comment) VALUES ($1, $2, $3, $4, $5)`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.RegionCode, dbReq.RegionName, time.Now(), time.Now(), "insert api")
	default:
		query := `UPDATE regions_table SET region_name=$2, modified=$3, comment=$4 WHERE region_code=$1`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.RegionCode, dbReq.RegionName, time.Now(), "update api")
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "UpsertRegionRecord : failed duplicate entry",
				Details:    "UpsertRegionRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertRegionRecord : failed ",
			Details:    "UpsertRegionRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertRegionRecord : rowsAffected failed",
			Details:    "UpsertRegionRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertRegionRecord : rowsAffected not 1",
			Details:    "UpsertRegionRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) UpsertRegionRecordViaFile(dbReq []*models.NewRegionRequest) []*models.NewRegionFailedResponse {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertRegionRecordViaFile start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertRegionRecordViaFile done")

	var failureResponse []*models.NewRegionFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO regions_table(region_code, region_name, created, modified, comment) VALUES ($1, $2, $3, $4, $5)`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.RegionName, time.Now(), time.Now(), "insert file")
		default:
			query := `UPDATE regions_table SET region_name=$2, modified=$3, comment=$4 WHERE region_code=$1`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.RegionName, time.Now(), "update file")
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.NewRegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.NewRegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.NewRegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.NewRegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.NewRegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertOutletRecord(dbReq *models.NewOutletRequest) *errors.Error {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertOutletRecord start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var regionName string

	rows, err := dbConn.Operation.QueryContext(ctx, "select region_name from regions_table where region_code=$1", dbReq.RegionCode)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertOutletRecord Select Query: failed",
			Details:    "UpsertOutletRecord Select Query: " + err.Error(),
		}
	}

	if err := rows.Scan(&regionName); err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertOutletRecord Select Scan: failed",
			Details:    "UpsertOutletRecord Select Scan: " + err.Error(),
		}
	}

	if regionName == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "RegionCode Not Found",
			Details:    "UpsertOutletRecord: blank regionName found",
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result sql.Result

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO outlets_table(outlet_code, outlet_name, region_code, region_name, latitude, longitude, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.OutletCode, dbReq.OutletName, dbReq.RegionCode, regionName, dbReq.Lat, dbReq.Lang, time.Now(), time.Now(), "insert api")
	default:
		query := `UPDATE outlets_table SET outlet_name=$2, region_code=$3, latitude=$4, longitude=$5, modified=$6, comment=$7  WHERE outlet_code=$1`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.OutletCode, dbReq.OutletName, dbReq.RegionCode, dbReq.Lat, dbReq.Lang, time.Now(), "update api")
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "UpsertOutletRecord : failed duplicate entry",
				Details:    "UpsertOutletRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertOutletRecord : failed",
			Details:    "UpsertOutletRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertOutletRecord : rowsAffected failed",
			Details:    "UpsertOutletRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertOutletRecord : rowsAffected not 1",
			Details:    "UpsertOutletRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) UpsertOutletRecordViaFile(dbReq []*models.NewOutletRequest) []*models.NewOutletFailedResponse {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertOutletRecordViaFile start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertOutletRecordViaFile done")

	var failureResponse []*models.NewOutletFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var regionName string

		rows, err := dbConn.Operation.QueryContext(ctx, "select region_name from regions_table where region_code=$1", eachRecord.RegionCode)
		if err != nil {
			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: fetch regionName query failed: " + err.Error()})
			continue
		}

		if err := rows.Scan(&regionName); err != nil {
			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: fetch regionName scan failed: " + err.Error()})
			continue
		}

		if regionName == "" {
			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: blank regionName found"})
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO outlets_table(outlet_code, outlet_name, region_code, region_name, latitude, longitude, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.OutletName, eachRecord.RegionCode, regionName, eachRecord.Lat, eachRecord.Lang, time.Now(), time.Now(), "insert file")
		default:
			query := `UPDATE outlets_table SET outlet_name=$2, region_code=$3, latitude=$4, longitude=$5, modified=$6, comment=$7  WHERE outlet_code=$1`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.OutletName, eachRecord.RegionCode, eachRecord.Lat, eachRecord.Lang, time.Now(), "update file")
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.NewOutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertProductRecord(dbReq *models.NewProductRequest) *errors.Error {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertProductRecord start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO products_table(product_code, product_name, weightage, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6)`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.ProductCode, dbReq.ProductName, dbReq.Weightage, time.Now(), time.Now(), "insert api")
	default:
		query := `UPDATE products_table SET product_name=$2, weightage=$3, modified=$4, comment=$5  WHERE product_code=$1`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.ProductCode, dbReq.ProductName, dbReq.Weightage, time.Now(), "update api")
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "UpsertProductRecord : failed duplicate entry",
				Details:    "UpsertProductRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertProductRecord : failed",
			Details:    "UpsertProductRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertProductRecord : rowsAffected failed",
			Details:    "UpsertProductRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertProductRecord : rowsAffected not 1",
			Details:    "UpsertProductRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) UpsertProductRecordViaFile(dbReq []*models.NewProductRequest) []*models.NewProductFailedResponse {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertProductRecordViaFile start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertProductRecordViaFile done")

	var failureResponse []*models.NewProductFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO products_table(product_code, product_name, weightage, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6)`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.ProductName, eachRecord.Weightage, time.Now(), time.Now(), "insert file")
		default:
			query := `UPDATE products_table SET product_name=$2, weightage=$3, modified=$4, comment=$5  WHERE product_code=$1`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.ProductName, eachRecord.Weightage, time.Now(), "update file")
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.NewProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.NewProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.NewProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.NewProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.NewProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertPlanRecord(dbReq *models.NewPlanRequest) *errors.Error {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertPlanRecord start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertPlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO plans_table(plan_code, plan_name, product_code, plan_price, sip_plan_rate, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.PlanCode, dbReq.PlanName, dbReq.ProductCode, dbReq.PlanPrice, dbReq.SIPPlanRate, time.Now(), time.Now(), "insert api")
	default:
		query := `UPDATE plans_table SET plan_name=$2, product_code=$3, plan_price=$4, sip_plan_rate=$5, modified=$6, comment=$7 WHERE plan_code=$1`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.PlanCode, dbReq.PlanName, dbReq.ProductCode, dbReq.PlanPrice, dbReq.SIPPlanRate, time.Now(), "update api")
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "UpsertPlansRecord : failed duplicate entry",
				Details:    "UpsertPlansRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertPlansRecord : failed",
			Details:    "UpsertPlansRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertPlansRecord : rowsAffected failed",
			Details:    "UpsertPlansRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertPlansRecord : rowsAffected not 1",
			Details:    "UpsertPlansRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) UpsertPlanRecordViaFile(dbReq []*models.NewPlanRequest) []*models.NewPlanFailedResponse {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertPlanRecordViaFile start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertPlanRecordViaFile done")

	var failureResponse []*models.NewPlanFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO plans_table(plan_code, plan_name, product_code, plan_price, sip_plan_rate, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.PlanName, eachRecord.ProductCode, eachRecord.PlanPrice, eachRecord.SIPPlanRate, time.Now(), time.Now(), "insert file")
		default:
			query := `UPDATE plans_table SET plan_name=$2, product_code=$3, plan_price=$4, sip_plan_rate=$5, modified=$6, comment=$7 WHERE plan_code=$1`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.PlanName, eachRecord.ProductCode, eachRecord.PlanPrice, eachRecord.SIPPlanRate, time.Now(), "update file")
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.NewPlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.NewPlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.NewPlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.NewPlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.NewPlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertTargetRecord(dbReq *models.NewTargetRequest) *errors.Error {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertTargetRecord start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var outletName string

	rows, err := dbConn.Operation.QueryContext(ctx, "select outlet_name from outlets_table where region_code=$1 AND outlet_code=$2", dbReq.RegionCode, dbReq.OutletCode)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertTargetRecord Select Query: failed",
			Details:    "UpsertTargetRecord Select Query: " + err.Error(),
		}
	}

	if err := rows.Scan(&outletName); err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertTargetRecord Select Scan: failed",
			Details:    "UpsertTargetRecord Select Scan: " + err.Error(),
		}
	}

	if outletName == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "OutletName Not Found",
			Details:    "UpsertTargetRecord: blank outletName found",
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result sql.Result

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO targets_table(target_code, region_code, outlet_code, year, month, baqati, hbb, hayyak, upgrade, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.TargetCode, dbReq.RegionCode, dbReq.OutletCode, dbReq.Year, dbReq.Month, dbReq.Baqati, dbReq.HBB, dbReq.Hayyak, dbReq.Upgrade, time.Now(), time.Now(), "insert api")
	default:
		query := `UPDATE targets_table SET baqati=$2, hbb=$3, hayyak=$4, upgrade=$5, modified=$6, comment=$7 WHERE target_code=$1`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.TargetCode, dbReq.Baqati, dbReq.HBB, dbReq.Hayyak, dbReq.Upgrade, time.Now(), "update api")
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "UpsertTargetRecord : failed duplicate entry",
				Details:    "UpsertTargetRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertTargetRecord : failed",
			Details:    "UpsertTargetRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertTargetRecord : rowsAffected failed",
			Details:    "UpsertTargetRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "UpsertTargetRecord : rowsAffected not 1",
			Details:    "UpsertTargetRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) UpsertTargetRecordViaFile(dbReq []*models.NewTargetRequest) []*models.NewTargetFailedResponse {

	logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertTargetRecordViaFile start")
	defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertTargetRecordViaFile done")

	var failureResponse []*models.NewTargetFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var outletName string

		rows, err := dbConn.Operation.QueryContext(ctx, "select outlet_name from outlets_table where region_code=$1 AND outlet_code=$2", eachRecord.RegionCode, eachRecord.OutletCode)
		if err != nil {
			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: fetch outletName query failed: " + err.Error()})
			continue
		}

		if err := rows.Scan(&outletName); err != nil {
			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: fetch outletName scan failed: " + err.Error()})
			continue
		}

		if outletName == "" {
			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: blank outletName found"})
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO targets_table(target_code, region_code, outlet_code, year, month, baqati, hbb, hayyak, upgrade, created, modified, comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, eachRecord.RegionCode, eachRecord.OutletCode, eachRecord.Year, eachRecord.Month, eachRecord.Baqati, eachRecord.HBB, eachRecord.Hayyak, eachRecord.Upgrade, time.Now(), time.Now(), "insert file")
		default:
			query := `UPDATE targets_table SET baqati=$2, hbb=$3, hayyak=$4, upgrade=$5, modified=$6, comment=$7 WHERE target_code=$1`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, eachRecord.Baqati, eachRecord.HBB, eachRecord.Hayyak, eachRecord.Upgrade, time.Now(), "update file")
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.NewTargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) DeleteRegionRecord(regionCodes []string) *errors.Error {

	logger.FileLogger.Info().Interface("record", regionCodes).Time("start", time.Now()).Msg("DeleteRegionRecord start")
	defer logger.FileLogger.Info().Interface("record", regionCodes).Time("end", time.Now()).Msg("DeleteRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM regions_table WHERE region_code IN ("

	for _, v := range regionCodes {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ")"

	result, err := dbConn.Operation.ExecContext(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "DeleteRegionRecord : failed duplicate entry",
				Details:    "DeleteRegionRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteRegionRecord : failed ",
			Details:    "DeleteRegionRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteRegionRecord : rowsAffected failed",
			Details:    "DeleteRegionRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteRegionRecord : rowsAffected not 1",
			Details:    "DeleteRegionRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteOutletRecord(outletCode []string) *errors.Error {

	logger.FileLogger.Info().Interface("record", outletCode).Time("start", time.Now()).Msg("DeleteOutletRecord start")
	defer logger.FileLogger.Info().Interface("record", outletCode).Time("end", time.Now()).Msg("DeleteOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM outlets_table WHERE outlet_code IN ("

	for _, v := range outletCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ")"

	result, err := dbConn.Operation.ExecContext(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "DeleteOutletRecord : failed duplicate entry",
				Details:    "DeleteOutletRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteOutletRecord : failed ",
			Details:    "DeleteOutletRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteOutletRecord : rowsAffected failed",
			Details:    "DeleteOutletRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteOutletRecord : rowsAffected not 1",
			Details:    "DeleteOutletRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteProductRecord(productCode []string) *errors.Error {

	logger.FileLogger.Info().Interface("record", productCode).Time("start", time.Now()).Msg("DeleteProductRecord start")
	defer logger.FileLogger.Info().Interface("record", productCode).Time("end", time.Now()).Msg("DeleteProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM products_table WHERE product_code IN ("

	for _, v := range productCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ")"

	result, err := dbConn.Operation.ExecContext(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "DeleteProductRecord : failed duplicate entry",
				Details:    "DeleteProductRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteProductRecord : failed ",
			Details:    "DeleteProductRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteProductRecord : rowsAffected failed",
			Details:    "DeleteProductRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteProductRecord : rowsAffected not 1",
			Details:    "DeleteProductRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeletePlanRecord(planCode []string) *errors.Error {

	logger.FileLogger.Info().Interface("record", planCode).Time("start", time.Now()).Msg("DeletePlanRecord start")
	defer logger.FileLogger.Info().Interface("record", planCode).Time("end", time.Now()).Msg("DeletePlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM plans_table WHERE plan_code IN ("

	for _, v := range planCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ")"

	result, err := dbConn.Operation.ExecContext(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "DeletePlanRecord : failed duplicate entry",
				Details:    "DeletePlanRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeletePlanRecord : failed ",
			Details:    "DeletePlanRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeletePlanRecord : rowsAffected failed",
			Details:    "DeletePlanRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeletePlanRecord : rowsAffected not 1",
			Details:    "DeletePlanRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteTargetRecord(targetCode []string) *errors.Error {

	logger.FileLogger.Info().Interface("record", targetCode).Time("start", time.Now()).Msg("DeleteTargetRecord start")
	defer logger.FileLogger.Info().Interface("record", targetCode).Time("end", time.Now()).Msg("DeleteTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM targets_table WHERE target_code IN ("

	for _, v := range targetCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ")"

	result, err := dbConn.Operation.ExecContext(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    "DeleteTargetRecord : failed duplicate entry",
				Details:    "DeleteTargetRecord failed: " + err.Error(),
			}
		}

		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteTargetRecord : failed ",
			Details:    "DeleteTargetRecord failed: " + err.Error(),
		}
	}

	insertedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteTargetRecord : rowsAffected failed",
			Details:    "DeleteTargetRecord rowsAffected failed: " + err.Error(),
		}
	}

	if insertedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteTargetRecord : rowsAffected not 1",
			Details:    "DeleteTargetRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		}
	}

	return nil
}

////////////////

// InsertEmployeeReport - delete and insert each outlet user report data daily in db
func (dbConn *DBImplementation) InsertEmployeeReport(usr *models.UserReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trx, err := dbConn.Operation.Begin()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "InsertEmployeeReport: begin trx failed - " + err.Error()}
	}

	if isQuarterReport {
		_, err = trx.ExecContext(ctx, "DELETE FROM user_reports_table where emp_id=$1 and isquarterly=true and quarter=$2 and year=$3", usr.UserId, quarter, year)
	} else {
		_, err = trx.ExecContext(ctx, "DELETE FROM user_reports_table where emp_id=$1 and isquarterly=false and month=$2 and year=$3", usr.UserId, month, year)
	}

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeReport: del failed",
			Details:    "InsertEmployeeReport del failed: " + err.Error(),
		}
	}

	rslt, err := trx.ExecContext(ctx, "INSERT INTO user_reports_table(emp_id, reporting_manager, isquarterly, month, quarter, year, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32)",
		usr.UserId,
		usr.ManagerId,
		usr.Isquarterly,
		usr.Month,
		usr.Quarter,
		usr.Year,
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
		usr.CreateAt,
	)

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeReport: insert failed",
			Details:    "InsertEmployeeReport insert failed: " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeReport: not inserted",
			Details:    "InsertEmployeeReport: rows affected " + strconv.Itoa(int(rows))}
	}

	err = trx.Commit()
	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeReport: commit failed",
			Details:    "InsertEmployeeReport: commit failed - " + err.Error()}
	}

	return nil
}

// InsertManagerReport - delete and insert each manager report data daily in db
func (dbConn *DBImplementation) InsertManagerReport(mgr *models.ManagerReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trx, err := dbConn.Operation.Begin()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "InsertManagerReport: begin trx failed - " + err.Error()}
	}

	if isQuarterReport {
		_, err = trx.ExecContext(ctx, "DELETE FROM manager_reports_table where mngr_id=$1 and isquarterly=true and quarter=$2 and year=$3", mgr.UserId, quarter, year)
	} else {
		_, err = trx.ExecContext(ctx, "DELETE FROM manager_reports_table where mngr_id=$1 and isquarterly=false and month=$2 and year=$3", mgr.UserId, month, year)
	}

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertManagerReport: del failed",
			Details:    "InsertManagerReport del failed: " + err.Error(),
		}
	}

	rslt, err := trx.ExecContext(ctx, "INSERT INTO manager_reports_table(mngr_id, position, reporting_manager, isquarterly, month, quarter, year, region, outlet, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)",
		mgr.UserId,
		mgr.Position,
		mgr.ManagerId,
		isQuarterReport,
		month,
		quarter,
		year,
		mgr.Region,
		mgr.Outlet,
		mgr.Latitude,
		mgr.Longitude,
		mgr.HBBSalesCount,
		mgr.HBBTargetCount,
		mgr.HayyakSalesCount,
		mgr.HayyakTargetCount,
		mgr.BaqatiSalesCount,
		mgr.BaqatiTargetCount,
		mgr.UpgradeSalesCount,
		mgr.UpgradeTargetCount,
		mgr.CreatedAt,
	)

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertManagerReport: insert failed",
			Details:    "InsertManagerReport insert failed: " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertManagerReport: not inserted",
			Details:    "InsertManagerReport: rows affected " + strconv.Itoa(int(rows))}
	}

	err = trx.Commit()
	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertManagerReport: commit failed",
			Details:    "InsertManagerReport: commit failed - " + err.Error()}
	}

	return nil
}

// InsertOutletReport - delete and insert each outlet report data daily in db
func (dbConn *DBImplementation) InsertOutletReport(outlet *models.OutletReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trx, err := dbConn.Operation.Begin()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "InsertOutletReport: begin trx failed - " + err.Error()}
	}

	if isQuarterReport {
		_, err = trx.ExecContext(ctx, "DELETE FROM outlet_reports_table where outlet=$1 and isquarterly=true and quarter=$2 and year=$3", outlet.Outlet, quarter, year)
	} else {
		_, err = trx.ExecContext(ctx, "DELETE FROM outlet_reports_table where outlet=$1 and isquarterly=false and month=$2 and year=$3", outlet.Outlet, month, year)
	}

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertOutletReport: del failed",
			Details:    "InsertOutletReport del failed: " + err.Error(),
		}
	}

	rslt, err := trx.ExecContext(ctx, "INSERT INTO outlet_reports_table(outlet, isquarterly, month, quarter, year, region, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)",
		outlet.Outlet,
		isQuarterReport,
		month,
		quarter,
		year,
		outlet.Region,
		outlet.Latitude,
		outlet.Longitude,
		outlet.HBBSalesCount,
		outlet.HBBTargetCount,
		outlet.HayyakSalesCount,
		outlet.HayyakTargetCount,
		outlet.BaqatiSalesCount,
		outlet.BaqatiTargetCount,
		outlet.UpgradeSalesCount,
		outlet.UpgradeTargetCount,
		outlet.CreatedAt)

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertOutletReport: insert failed",
			Details:    "InsertOutletReport insert failed: " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertOutletReport: not inserted",
			Details:    "InsertOutletReport: rows affected " + strconv.Itoa(int(rows))}
	}

	err = trx.Commit()
	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertOutletReport: commit failed",
			Details:    "InsertOutletReport: commit failed - " + err.Error()}
	}

	return nil
}

// InsertRegionReport - delete and insert each region report data daily in db
func (dbConn *DBImplementation) InsertRegionReport(region *models.RegionReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trx, err := dbConn.Operation.Begin()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "InsertRegionReport: begin trx failed - " + err.Error()}
	}

	if isQuarterReport {
		_, err = trx.ExecContext(ctx, "DELETE FROM region_reports_table where region=$1 and isquarterly=true and quarter=$2 and year=$3", region.Region, quarter, year)
	} else {
		_, err = trx.ExecContext(ctx, "DELETE FROM region_reports_table where region=$1 and isquarterly=false and month=$2 and year=$3", region.Region, month, year)
	}

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertRegionReport: del failed",
			Details:    "InsertRegionReport del failed: " + err.Error(),
		}
	}

	rslt, err := trx.ExecContext(ctx, "INSERT INTO region_reports_table(region, isquarterly, month, quarter, year, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)",
		region.Region,
		isQuarterReport,
		month,
		quarter,
		year,
		region.HBBSalesCount,
		region.HBBTargetCount,
		region.HayyakSalesCount,
		region.HayyakTargetCount,
		region.BaqatiSalesCount,
		region.BaqatiTargetCount,
		region.UpgradeSalesCount,
		region.UpgradeTargetCount,
		region.CreatedAt,
	)

	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertRegionReport: insert failed",
			Details:    "InsertRegionReport insert failed: " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertRegionReport: not inserted",
			Details:    "InsertRegionReport: rows affected " + strconv.Itoa(int(rows))}
	}

	err = trx.Commit()
	if err != nil {
		trx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertRegionReport: commit failed",
			Details:    "InsertRegionReport: commit failed - " + err.Error()}
	}

	return nil
}

// FetchUserReport - pre-computed user report
func (dbConn *DBImplementation) FetchUserReport(userId string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select emp_id, reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created from user_reports_table where emp_id=$1 AND isquarterly=true AND quarter=$2 AND year=$3", userId, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select emp_id, reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created from user_reports_table where emp_id=$1 AND isquarterly=false AND month=$2 AND year=$3", userId, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchUserReport: query failed",
			Details:    "FetchUserReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.UserReportDB{}

	for rows.Next() {

		err = rows.Scan(
			&rcrd.UserId,
			&rcrd.ManagerId,
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
			&rcrd.CreateAt,
		)

		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserReport: scan failed",
				Details:    "FetchUserReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

// FetchManagerReport - pre-computed manager report
func (dbConn *DBImplementation) FetchManagerReport(mgrId, position string, isQuarterReport bool, quarter, month, year int) (*models.ManagerReportDB, *errors.Error) {

	_, ok := constant.ManagerMap[position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select mngr_id, position, reporting_manager, region, outlet, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from reports_table where mngr_id=$1 AND isquarterly=true AND quarter=$2 AND year=$3", mgrId, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select mngr_id, position, reporting_manager, region, outlet, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from reports_table where mngr_id=$1 AND isquarterly=false AND month=$2 AND year=$3", mgrId, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchManagerReport: query failed",
			Details:    "FetchManagerReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.ManagerReportDB{}

	for rows.Next() {

		err = rows.Scan(
			&rcrd.UserId,
			&rcrd.Position,
			&rcrd.ManagerId,
			&rcrd.Region,
			&rcrd.Outlet,
			&rcrd.Latitude,
			&rcrd.Longitude,
			&rcrd.HBBSalesCount,
			&rcrd.HBBTargetCount,
			&rcrd.HayyakSalesCount,
			&rcrd.HayyakTargetCount,
			&rcrd.BaqatiSalesCount,
			&rcrd.BaqatiTargetCount,
			&rcrd.UpgradeSalesCount,
			&rcrd.UpgradeTargetCount,
			&rcrd.CreatedAt,
		)

		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchManagerReport: scan failed",
				Details:    "FetchManagerReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

// FetchOutletReport - pre-computed outlet report
func (dbConn *DBImplementation) FetchOutletReport(outlet, userPosition string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	_, ok := constant.ManagerMap[userPosition]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select outlet, region, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from outlet_reports_table where outlet=$1 AND isquarterly=true AND quarter=$2 AND year=$3", outlet, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select outlet, region, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from outlet_reports_table where outlet=$1 AND isquarterly=false AND month=$2 AND year=$3", outlet, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletReport: query failed",
			Details:    "FetchOutletReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.OutletReportDB{}

	for rows.Next() {

		err = rows.Scan(
			&rcrd.Outlet,
			&rcrd.Region,
			&rcrd.Latitude,
			&rcrd.Longitude,
			&rcrd.HBBSalesCount,
			&rcrd.HBBTargetCount,
			&rcrd.HayyakSalesCount,
			&rcrd.HayyakTargetCount,
			&rcrd.BaqatiSalesCount,
			&rcrd.BaqatiTargetCount,
			&rcrd.UpgradeSalesCount,
			&rcrd.UpgradeTargetCount,
			&rcrd.CreatedAt,
		)

		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletReport: scan failed",
				Details:    "FetchOutletReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

// FetchRegionReport - pre-computed region report
func (dbConn *DBImplementation) FetchRegionReport(region, userPosition string, isQuarterReport bool, quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	_, ok := constant.ManagerMap[userPosition]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select region, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from region_reports_table where region=$1 AND isquarterly=true AND quarter=$2 AND year=$3", region, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select region, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from region_reports_table where region=$1 AND isquarterly=false AND month=$2 AND year=$3", region, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchRegionReport: query failed",
			Details:    "FetchRegionReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.RegionReportDB{}

	for rows.Next() {

		err = rows.Scan(
			&rcrd.Region,
			&rcrd.HBBSalesCount,
			&rcrd.HBBTargetCount,
			&rcrd.HayyakSalesCount,
			&rcrd.HayyakTargetCount,
			&rcrd.BaqatiSalesCount,
			&rcrd.BaqatiTargetCount,
			&rcrd.UpgradeSalesCount,
			&rcrd.UpgradeTargetCount,
			&rcrd.CreatedAt,
		)

		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchRegionReport: scan failed",
				Details:    "FetchRegionReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

////////////////

// FetchUserSales - all sale transaction of a given user and time period
func (dbConn *DBImplementation) FetchUserSales(clientID string, month, year, quarter int, isQuarterReport bool) (*models.SaleTransaction, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", clientID).Int("month", month).Int("year", year).Int("quarter", quarter).Bool("isQuarterReport", isQuarterReport).Time("start", time.Now()).Msg("FetchUserSales start")
	defer logger.FileLogger.Info().Str("clientID", clientID).Int("month", month).Int("year", year).Int("quarter", quarter).Bool("isQuarterReport", isQuarterReport).Time("end", time.Now()).Msg("FetchUserSales done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rows *sql.Rows
		err  error
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category, b.subcategory, b.sales_date, b.inactive_date from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and b.quarter= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category, b.subcategory, b.sales_date, b.inactive_date from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and EXTRACT('month' from b.sales_date)= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchUserSales Query: failed",
			Details:    "FetchUserSales Query: " + err.Error(),
		}
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
		subcategory   interface{}
		salesDate     time.Time
		inactiveDate  time.Time
	)

	trx := &models.SaleTransaction{}
	trx.SaleRecords = make(map[string][]*models.SaleRecord)

	count := 0

	for rows.Next() {
		count++
		if err := rows.Scan(&position, &outlet, &region, &prev_product, &curr_product, &quantity, &product_group, &category, &subcategory, &salesDate, &inactiveDate); err != nil {
			return trx, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserSales Scan: failed",
				Details:    "FetchUserSales Scan: " + err.Error(),
			}
		}

		trx.EmpID = clientID
		trx.Position = position
		trx.Outlet = outlet
		trx.Region = region

		switch {
		case strings.Contains(product_group, "FIXED"):
			product_group = "HBB"
		case strings.Contains(product_group, "PREPAID"):
			product_group = "HAYYAK"
		case strings.Contains(product_group, "POSTPAID"):
			product_group = "BAQATI"
		case product_group == "HBB" || product_group == "HAYYAK" || product_group == "BAQATI":

		default:
			return trx, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "Product Group Not Found",
				Details:    "FetchUserSales: invalid productGrp-" + product_group,
			}
		}

		lapsedSecs := salesDate.Sub(inactiveDate).Seconds()

		if objecthelper.RoundTime(lapsedSecs/2600640) < 3 { //Months<3
			if !strings.Contains(subcategory.(string), "INVOLUNTARY DEACTIVATION") {
				//terminated logs
				continue
			}
		}

		rcrd := &models.SaleRecord{
			Quantity:        quantity,
			Category:        category,
			PreviousProduct: prev_product,
			CurrentProduct:  curr_product,
		}

		if trx.SaleRecords[product_group] == nil {
			trx.SaleRecords[product_group] = make([]*models.SaleRecord, 0)
		}

		trx.SaleRecords[product_group] = append(trx.SaleRecords[product_group], rcrd)
	}

	return trx, nil
}

// FetchUserSalesPerProduct - all sale transaction of a given user, product and time period
func (dbConn *DBImplementation) FetchUserSalesCountPerPlan(clientID, product string, month, year, quarter int, isQuarterReport bool) ([]*models.SaleCountPerPlan, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rows *sql.Rows
		err  error
	)

	response := make([]*models.SaleCountPerPlan, 0)

	if product == constant.UpgradeCategory {

		var (
			quantity     int
			curr_product string
			subcategory  interface{}
			salesDate    time.Time
			inactiveDate time.Time
		)

		if isQuarterReport {
			rows, err = dbConn.Operation.QueryContext(ctx, "select quantity, curr_product, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and category_type='UPGRADE' and quarter= $2 and EXTRACT('year' from sales_date)= $3", clientID, quarter, year)
		} else {
			rows, err = dbConn.Operation.QueryContext(ctx, "select quantity, curr_product, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and category_type='UPGRADE' and EXTRACT('month' from b.sales_date)= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, month, year)
		}

		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserSalesCountPerPlan Query: failed",
				Details:    "FetchUserSalesCountPerPlan Query: " + err.Error(),
			}
		}

		for rows.Next() {
			if err := rows.Scan(&quantity, &curr_product, &subcategory, &salesDate, &inactiveDate); err != nil {
				return nil, &errors.Error{
					StatusCode: http.StatusInternalServerError,
					Code:       errors.InternalServerError,
					Message:    "FetchUserSalesCountPerPlan Scan: failed",
					Details:    "FetchUserSalesCountPerPlan Scan: " + err.Error(),
				}
			}

			lapsedSecs := salesDate.Sub(inactiveDate).Seconds()

			if objecthelper.RoundTime(lapsedSecs/2600640) < 3 { //Months<3
				if !strings.Contains(subcategory.(string), "INVOLUNTARY DEACTIVATION") {
					//terminated logs
					continue
				}
			}

			response = append(response, &models.SaleCountPerPlan{
				CurrentPlan: curr_product,
				Quantity:    quantity,
			})
		}

		return response, nil
	}

	//for non-upgrade sales

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select quantity, curr_product, product_group, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and category_type='NEW' and quarter= $2 and EXTRACT('year' from sales_date)= $3", clientID, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select quantity, curr_product, product_group, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and category_type='NEW' and EXTRACT('month' from b.sales_date)= $2 and EXTRACT('year' from b.sales_date)= $3", clientID, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchUserSalesCountPerPlan Query: failed",
			Details:    "FetchUserSalesCountPerPlan Query: " + err.Error(),
		}
	}

	var (
		quantity      int
		curr_product  string
		product_group string
		subcategory   interface{}
		salesDate     time.Time
		inactiveDate  time.Time
	)

	for rows.Next() {
		if err := rows.Scan(&quantity, &curr_product, &product_group, &subcategory, &salesDate, &inactiveDate); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserSalesCountPerPlan Scan: failed",
				Details:    "FetchUserSalesCountPerPlan Scan: " + err.Error(),
			}
		}

		switch {
		case strings.Contains(product_group, "FIXED"):
			product_group = "HBB"
		case strings.Contains(product_group, "PREPAID"):
			product_group = "HAYYAK"
		case strings.Contains(product_group, "POSTPAID"):
			product_group = "BAQATI"
		default:
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "Product Group Not Found",
				Details:    "FetchUserSales: invalid productGrp-" + product_group,
			}
		}

		lapsedSecs := salesDate.Sub(inactiveDate).Seconds()

		if objecthelper.RoundTime(lapsedSecs/2600640) < 3 { //Months<3
			if !strings.Contains(subcategory.(string), "INVOLUNTARY DEACTIVATION") {
				//terminated logs
				continue
			}
		}

		if product_group == product {
			response = append(response, &models.SaleCountPerPlan{
				CurrentPlan: curr_product,
				Quantity:    quantity,
			})
		}
	}

	return response, nil
}

// FetchUserPosition - user position from user master
func (dbConn *DBImplementation) FetchUserPosition(id string) (string, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Time("start", time.Now()).Msg("FetchUserPosition start")
	defer logger.FileLogger.Info().Str("clientID", id).Time("end", time.Now()).Msg("FetchUserPosition done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select usr_position from user_master_table where emp_id=$1", id)
	if err != nil {
		return "", &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchUserPosition: query failed",
			Details:    "FetchUserPosition query failed: " + err.Error(),
		}
	}

	var (
		count    int
		position string
	)

	for rows.Next() {
		count++

		if err := rows.Scan(&position); err != nil {
			return "", &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserPosition: scan failed",
				Details:    "FetchUserPosition scan failed: " + err.Error(),
			}
		}

		if count > 1 {
			return "", &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserPosition: more than one record found",
				Details:    "FetchUserPosition: more than one record found",
			}
		}
	}

	return position, nil
}

func (dbConn *DBImplementation) FetchOutletDetailsViaName(regionName, outletName string) (*models.GetOutletResponse, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	outletResponse := &models.GetOutletResponse{}

	err := dbConn.Operation.QueryRowContext(ctx,
		"SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table where region_name=$1 and outlet_name=$2", regionName, outletName).
		Scan(&outletResponse.OutletCode, &outletResponse.OutletName, &outletResponse.RegionCode, &outletResponse.RegionName, &outletResponse.Lat, &outletResponse.Lang)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletDetails: query failed",
			Details:    "FetchOutletDetails query failed: " + err.Error(),
		}
	}

	if outletResponse.OutletCode == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "OUTLATE NAME NOT FOUND",
			Details:    "FetchOutletDetails value not found, " + fmt.Sprintf("region:%v, outlet:%v", regionName, outletName),
		}
	}

	return outletResponse, nil
}

func (dbConn *DBImplementation) FetchAllOutletNames(regionName string) ([]string, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	outlets := make([]string, 0)

	rows, err := dbConn.Operation.QueryContext(ctx, "SELECT outlet_name FROM outlets_table where region_name=$1", regionName)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllOutletNames: query failed",
			Details:    "FetchAllOutletNames query failed: " + err.Error(),
		}
	}

	var outlentName string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&outlentName); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchAllOutletNames: scan failed",
				Details:    "FetchAllOutletNames scan failed: " + err.Error(),
			}
		}

		outlets = append(outlets, outlentName)
	}

	return outlets, nil
}

////////////////

// FetchTotalNonManagerialEmployees - all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchTotalNonManagerialEmployees(region, outlet string) (int, *errors.Error) {

	// logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("start", time.Now()).Msg("FetchTotalNonManagerialEmployees start")
	// defer logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("end", time.Now()).Msg("FetchTotalNonManagerialEmployees done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	err := dbConn.Operation.QueryRowContext(ctx, "select count(*) from user_master_table where region= $1 and outlet= $2 and usr_position!='BACK OFFICE'",
		region, outlet).Scan(&count)
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchTotalNonManagerialEmployees Query: failed",
			Details:    "FetchTotalNonManagerialEmployees Query: " + err.Error(),
		}
	}

	return count, nil
}

// InsertUserMasterRecord reads data from user master input file and inserts into table
func (dbConn *DBImplementation) InsertUserMasterRecord(emp *models.Employee) *errors.Error {

	logger.FileLogger.Info().Interface("record", emp).Time("start", time.Now()).Msg("InsertUserMasterRecord start")
	defer logger.FileLogger.Info().Interface("record", emp).Time("end", time.Now()).Msg("InsertUserMasterRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into User_Master_Table(created, modified, emp_id, name, usr_position, region, outlet, mobile_number ,remark) values($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), emp.Id, emp.Name, emp.Position, emp.Region, emp.Outlet, emp.MobileNos, emp.Remark)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertUserMasterRecord: failed",
			Details:    "InsertUserMasterRecord: " + err.Error(),
		}
	}

	return nil
}

// InsertSaleTransaction reads data from input file and inserts into table
func (dbConn *DBImplementation) InsertSaleTransaction(trx *models.Transaction) *errors.Error {

	logger.FileLogger.Info().Interface("record", trx).Time("start", time.Now()).Msg("InsertSaleTransaction start")
	defer logger.FileLogger.Info().Interface("record", trx).Time("end", time.Now()).Msg("InsertSaleTransaction done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into Sale_Transaction_Table(created, modified, region, outlet, emp_id, category_type, prev_product, curr_product, quantity, category, product_group, sales_date, quarter) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), trx.Reqion, trx.Outlet, trx.Id, trx.CategoryType, trx.PreviousProduct,
		trx.CurrentProduct, trx.Quantity, trx.Category, trx.ProductGrp, trx.SaleDate, trx.Quarter)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertSaleTransaction: failed",
			Details:    "InsertSaleTransaction: " + err.Error(),
		}
	}

	return nil
}

// DeleteProductPlans removes product plan data from table
func (dbConn *DBImplementation) DeleteProductPlans() *errors.Error {

	logger.FileLogger.Info().Interface("record", "*").Time("start", time.Now()).Msg("DeleteProductPlans start")
	defer logger.FileLogger.Info().Interface("record", "*").Time("end", time.Now()).Msg("DeleteProductPlans done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := dbConn.Operation.ExecContext(ctx, "DELETE FROM Product_Plans_Table")
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteProductPlans: failed",
			Details:    "DeleteProductPlans: " + err.Error(),
		}
	}

	return nil
}

// InsertProductPlans reads data from plan input file and inserts into table
func (dbConn *DBImplementation) InsertProductPlans(plan *models.ProductPlans) *errors.Error {

	logger.FileLogger.Info().Interface("record", plan).Time("start", time.Now()).Msg("InsertProductPlans start")
	defer logger.FileLogger.Info().Interface("record", plan).Time("end", time.Now()).Msg("InsertProductPlans done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into Product_Plans_Table(created, modified, plan_name, plan_price, sip_plan_rate) values($1,$2,$3,$4,$5)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), plan.PlanName, plan.PlanPrice, plan.SIPRatePlan)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertProductPlans: failed",
			Details:    "InsertProductPlans: " + err.Error(),
		}
	}

	return nil
}

// DeleteStoreTargets removes store target data from table
func (dbConn *DBImplementation) DeleteStoreTargets() *errors.Error {

	logger.FileLogger.Info().Interface("record", "*").Time("start", time.Now()).Msg("DeleteStoreTargets start")
	defer logger.FileLogger.Info().Interface("record", "*").Time("end", time.Now()).Msg("DeleteStoreTargets done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := dbConn.Operation.ExecContext(ctx, "DELETE FROM Sales_Target_Table")
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteStoreTargets: failed",
			Details:    "DeleteStoreTargets: " + err.Error(),
		}
	}

	return nil
}

// DeleteHierarchy removes all hierarchy data from table
func (dbConn *DBImplementation) DeleteHierarchy() *errors.Error {

	logger.FileLogger.Info().Interface("record", "*").Time("start", time.Now()).Msg("DeleteHierarchy start")
	defer logger.FileLogger.Info().Interface("record", "*").Msg("DeleteHierarchy done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := dbConn.Operation.ExecContext(ctx, "DELETE FROM user_hierarchy_table")
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "DeleteHierarchy: failed",
			Details:    "DeleteHierarchy: " + err.Error(),
		}
	}

	return nil
}

// InsertStoreTargets reads data from target input file and inserts into table
func (dbConn *DBImplementation) InsertStoreTargets(region, outlet string, month int, trgt *models.ProductGrpTarget) *errors.Error {

	logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Int("month", month).Interface("record", trgt).Time("start", time.Now()).Msg("InsertStoreTargets start")
	defer logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Int("month", month).Interface("record", trgt).Time("end", time.Now()).Msg("InsertStoreTargets done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into Sales_Target_Table(created, modified, region, outlet, months, Baqati, HBB, Hayyak, Upgrade) values($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	_, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), time.Now(), region, outlet, month, trgt.Baqati, trgt.HBB, trgt.Hayyak, trgt.Upgrade)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertStoreTargets: failed",
			Details:    "InsertStoreTargets: " + err.Error(),
		}
	}

	return nil
}

// InsertValueOfPoints reads data from ScakeOrder input file and inserts into table
func (dbConn *DBImplementation) InsertValueOfPoints(category string, percent int, vop float64) (int64, *errors.Error) {

	logger.FileLogger.Info().Str("category", category).Int("percent", percent).Float64("vop", vop).Time("start", time.Now()).Msg("InsertValueOfPoints start")
	defer logger.FileLogger.Info().Str("category", category).Int("percent", percent).Float64("vop", vop).Time("end", time.Now()).Msg("InsertValueOfPoints done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into Scale_Order_Table(category, percentage, vop) values($1,$2,$3)"

	rslt, err := dbConn.Operation.ExecContext(ctx, query, category, percent, vop)
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertValueOfPoints: failed",
			Details:    "InsertValueOfPoints: " + err.Error(),
		}
	}

	id, err := rslt.LastInsertId()
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertValueOfPoints LastInsertId: failed",
			Details:    "InsertValueOfPoints LastInsertId: " + err.Error(),
		}
	}

	return id, nil
}

func (dbConn *DBImplementation) FetchPromotersPresentYearSales(clientID string) (map[int]*models.SaleTransaction, map[int]*models.SaleTransaction, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", clientID).Time("start", time.Now()).Msg("FetchPromotersPresentYearSales start")
	defer logger.FileLogger.Info().Str("clientID", clientID).Time("end", time.Now()).Msg("FetchPromotersPresentYearSales done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select a.usr_position, a.outlet, a.region, b.prev_product, b.curr_product, b.quantity, b.product_group, b.category, b.sales_date from User_Master_Table a INNER JOIN Sale_Transaction_Table b on a.emp_id=b.emp_id where a.emp_id= $1 and EXTRACT('year' from b.sales_date)= $2", clientID, 2021)
	if err != nil {
		return nil, nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchPromotersPresentYearSales Query: failed",
			Details:    "FetchPromotersPresentYearSales Query: " + err.Error(),
		}
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

	monthly := make(map[int]*models.SaleTransaction)
	quarterly := make(map[int]*models.SaleTransaction)

	for rows.Next() {
		if err := rows.Scan(&position, &outlet, &region, &prev_product, &curr_product, &quantity, &product_group, &category, &sale_date); err != nil {
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPromotersPresentYearSales Scan: failed",
				Details:    "FetchPromotersPresentYearSales Scan: " + err.Error(),
			}
		}

		rcrd := &models.SaleRecord{
			Quantity:        quantity,
			Category:        category,
			PreviousProduct: prev_product,
			CurrentProduct:  curr_product,
			SaleDate:        sale_date,
		}

		v, ok := monthly[int(sale_date.Month())]
		if !ok {
			trx := &models.SaleTransaction{}
			trx.EmpID = clientID
			trx.Position = position
			trx.Outlet = outlet
			trx.Region = region
			trx.SaleRecords = map[string][]*models.SaleRecord{
				product_group: {rcrd},
			}

			monthly[int(sale_date.Month())] = trx
		} else {

			if v.SaleRecords[product_group] == nil {
				v.SaleRecords[product_group] = make([]*models.SaleRecord, 0)
			}

			v.SaleRecords[product_group] = append(v.SaleRecords[product_group], rcrd)

			monthly[int(sale_date.Month())] = v
		}

		v, ok = quarterly[getQuarterFromMonth(int(sale_date.Month()))]
		if !ok {
			trx := &models.SaleTransaction{}
			trx.EmpID = clientID
			trx.Position = position
			trx.Outlet = outlet
			trx.Region = region
			trx.SaleRecords = map[string][]*models.SaleRecord{
				product_group: {rcrd},
			}

			quarterly[getQuarterFromMonth(int(sale_date.Month()))] = trx
		} else {

			if v.SaleRecords[product_group] == nil {
				v.SaleRecords[product_group] = make([]*models.SaleRecord, 0)
			}

			v.SaleRecords[product_group] = append(v.SaleRecords[product_group], rcrd)

			quarterly[getQuarterFromMonth(int(sale_date.Month()))] = v
		}
	}

	return monthly, quarterly, nil
}

// FetchTotalNonManagerialEmployees returns all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchEmployeeContactNumber(clientID string) (int, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", clientID).Time("start", time.Now()).Msg("FetchEmployeeContactNumber start")
	defer logger.FileLogger.Info().Str("clientID", clientID).Time("end", time.Now()).Msg("FetchEmployeeContactNumber done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select mobile_number from User_Master_Table where emp_id= $1", clientID)
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchEmployeeContactNumber Query: failed",
			Details:    "FetchEmployeeContactNumber Query: " + err.Error(),
		}
	}

	var (
		count         int
		contactNumber string
	)

	for rows.Next() {

		if count == 1 {
			return -1, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchEmployeeContactNumber: found more than on record",
				Details:    "FetchEmployeeContactNumber: found more than on record",
			}
		}

		if err := rows.Scan(&contactNumber); err != nil {
			return -1, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchEmployeeContactNumber Scan: failed",
				Details:    "FetchEmployeeContactNumber Scan: " + err.Error(),
			}
		}

		count++
	}

	if contactNumber == "" {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchEmployeeContactNumber: contact number not found",
			Details:    "FetchEmployeeContactNumber: contact number not found",
		}
	}

	mobNum, err := strconv.Atoi(contactNumber)
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchEmployeeContactNumber: contact number parsing failed",
			Details:    "FetchEmployeeContactNumber contact number parsing failed: " + err.Error(),
		}
	}

	return mobNum, nil
}

func (dbConn *DBImplementation) InsertStaffOTP(update bool, otpModel *models.OtpModel) *errors.Error {

	logger.FileLogger.Info().Interface("record", otpModel).Time("start", time.Now()).Msg("InsertStaffOTP start")
	defer logger.FileLogger.Info().Interface("record", otpModel).Time("end", time.Now()).Msg("InsertStaffOTP done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error

	if update {
		_, err = dbConn.Operation.QueryContext(ctx, "UPDATE user_authentication_table SET user_active=$1, otp=$2, otp_resend_count=$3, otp_verify_failed_count=$4, otp_used=$5, otp_create_time=$6, otp_update_time=$7, last_login_time=$8, modified=$9, modified_reason=$10",
			otpModel.IsActiveUser,
			otpModel.OTP,
			otpModel.OTPResendCount,
			otpModel.OTPVerifyFailedCount,
			otpModel.OTPUsed,
			otpModel.OtpCreateTime,
			otpModel.OtpLastUpdateTime,
			otpModel.LastLoginTime,
			otpModel.UpdatedTime,
			otpModel.UpdateReason,
		)

		if err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "UpdateStaffOTP: update failed",
				Details:    "UpdateStaffOTP update failed: " + err.Error(),
			}
		}

	} else {
		_, err = dbConn.Operation.QueryContext(ctx, "INSERT INTO user_authentication_table(emp_id, user_active, otp, otp_resend_count, otp_verify_failed_count, otp_used, otp_create_time, otp_update_time, last_login_time, created, modified, modified_reason) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);",
			otpModel.StaffID,
			otpModel.IsActiveUser,
			otpModel.OTP,
			otpModel.OTPResendCount,
			otpModel.OTPVerifyFailedCount,
			otpModel.OTPUsed,
			otpModel.OtpCreateTime,
			otpModel.OtpLastUpdateTime,
			otpModel.LastLoginTime,
			otpModel.CreatedTime,
			otpModel.UpdatedTime,
			otpModel.UpdateReason,
		)

		if err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "InsertStaffOTP: insert failed",
				Details:    "InsertStaffOTP insert failed: " + err.Error(),
			}
		}
	}

	return nil
}

func (dbConn *DBImplementation) FetchOTPValidity(clientID string) (string, *models.OtpModel, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", clientID).Time("start", time.Now()).Msg("FetchOTPValidity start")
	defer logger.FileLogger.Info().Str("clientID", clientID).Time("end", time.Now()).Msg("FetchOTPValidity done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select a.mobile_number, a.name, a.usr_position, b.emp_id, b.user_active, b.otp, b.otp_resend_count, b.otp_verify_failed_count, b.otp_used, b.otp_create_time, b.otp_update_time, b.last_login_time from User_Master_Table a inner join user_authentication_table b on a.emp_id=b.emp_id and a.emp_id= $1", clientID)
	if err != nil {
		return "", nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOTPValidity: query failed",
			Details:    "FetchOTPValidity query failed: " + err.Error(),
		}
	}

	var (
		count   int
		mobile  string
		otModel = &models.OtpModel{}
	)

	for rows.Next() {

		if count == 1 {
			return "", nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOTPValidity: found more than on record",
				Details:    "FetchOTPValidity: found more than on record",
			}
		}

		if err := rows.Scan(&mobile, &otModel.StaffName, &otModel.StaffPosition, &otModel.StaffID, &otModel.IsActiveUser, &otModel.OTP, &otModel.OTPResendCount, &otModel.OTPVerifyFailedCount, &otModel.OTPUsed, &otModel.OtpCreateTime, &otModel.OtpLastUpdateTime, &otModel.LastLoginTime); err != nil {
			return "", nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOTPValidity: scan failed",
				Details:    "FetchOTPValidity scan failed: " + err.Error(),
			}
		}

		count++
	}

	return mobile, otModel, nil
}

// FetchAllGivenPositionEmpIdAndTheirManager returns an array of all employees of the given position along with their managers, returns empty array if no records found
func (dbConn *DBImplementation) FetchAllGivenPositionEmpIdAndTheirManager(position string) ([]*models.UserLead, *errors.Error) {

	logger.FileLogger.Info().Str("position", position).Time("start", time.Now()).Msg("FetchAllGivenPositionEmpIdAndTheirManager start")
	defer logger.FileLogger.Info().Str("position", position).Time("end", time.Now()).Msg("FetchAllGivenPositionEmpIdAndTheirManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if !(position == constant.VicePresident || position == constant.GeneralManager || position == constant.SeniorManager || position == constant.Manager || position == constant.OutletManager || position == constant.Promoter) {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllGivenPositionEmpIdAndTheirManager: incorrect position",
			Details:    "FetchAllGivenPositionEmpIdAndTheirManager: incorrect position req- " + position,
		}
	}

	var query string

	switch position {
	case constant.VicePresident:
		query = "select distinct vice_president, 0 from user_hierarchy_table"
	case constant.GeneralManager:
		query = "select distinct general_manager, vice_president from user_hierarchy_table"
	case constant.SeniorManager:
		query = "select distinct senior_manager, general_manager from user_hierarchy_table"
	case constant.Manager:
		query = "select distinct manager, senior_manager from user_hierarchy_table"
	case constant.OutletManager:
		query = "select distinct outlet_manager, manager from user_hierarchy_table"
	case constant.Promoter:
		query = "select distinct emp_id, outlet_manager from user_hierarchy_table"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllGivenPositionEmpIdAndTheirManager: query failed",
			Details:    "FetchAllGivenPositionEmpIdAndTheirManager query failed: " + err.Error(),
		}
	}

	var array = make([]*models.UserLead, 0)

	for rows.Next() {

		usr := &models.UserLead{}

		if err := rows.Scan(&usr.EmpID, &usr.ReportTo); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchAllGivenPositionEmpIdAndTheirManager: scan failed failed",
				Details:    "FetchAllGivenPositionEmpIdAndTheirManager scan failed: " + err.Error(),
			}
		}

		array = append(array, usr)
	}

	return array, nil
}

// FetchAllEmployeesReportingToGivenManager returns all employees reporting to given manager, returns empty array if no records found
func (dbConn *DBImplementation) FetchAllEmployeesReportingToGivenManagerRole(managerPosition string, managerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllEmployeesReportingToGivenManagerRole start")
	defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllEmployeesReportingToGivenManagerRole done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if !(managerPosition == "vice_president" || managerPosition == "general_manager" || managerPosition == "senior_manager" || managerPosition == "manager" || managerPosition == "outlet_manager") {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllEmployeesReportingToGivenManagerRole: incorrect position",
			Details:    "FetchAllEmployeesReportingToGivenManagerRole: incorrect position req- " + managerPosition,
		}
	}

	var query string

	switch managerPosition {
	case "vice_president":
		query = "select distinct a.general_manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.general_manager=b.emp_id where a.vice_president= $1"
	case "general_manager":
		query = "select distinct a.senior_manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.senior_manager=b.emp_id where a.general_manager= $1"
	case "senior_manager":
		query = "select distinct a.manager,  b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.manager=b.emp_id where a.senior_manager= $1"
	case "manager":
		query = "select distinct a.outlet_manager,  b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager= $1"
	case "outlet_manager":
		query = "select distinct a.emp_id, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.emp_id=b.emp_id where a.outlet_manager= $1"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllEmployeesReportingToGivenManagerRole: query failed",
			Details:    "FetchAllEmployeesReportingToGivenManagerRole query failed: " + err.Error(),
		}
	}

	var (
		array = make([]*models.ReportingStaffIDAndName, 0)
	)

	for rows.Next() {

		usr := &models.ReportingStaffIDAndName{}

		if err := rows.Scan(&usr.ID, &usr.Position, &usr.Name); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchAllEmployeesReportingToGivenManagerRole: scan failed failed",
				Details:    "FetchAllEmployeesReportingToGivenManagerRole scan failed: " + err.Error(),
			}
		}

		array = append(array, usr)
	}

	return array, nil
}

func (dbConn *DBImplementation) FetchAllOutletManagersReportingToGivenManagerRole(managerPosition string, managerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllOutletManagersReportingToGivenManagerRole start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllOutletManagersReportingToGivenManagerRole done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if !(managerPosition == "vice_president" || managerPosition == "general_manager" || managerPosition == "senior_manager" || managerPosition == "manager" || managerPosition == "outlet_manager") {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllOutletManagersReportingToGivenManagerRole: incorrect position",
			Details:    "FetchAllOutletManagersReportingToGivenManagerRole: incorrect position req- " + managerPosition,
		}
	}

	var query string

	switch managerPosition {
	case "vice_president":
		query = "select distinct a.outlet_manager, b.outlet, b.region from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.vice_president= $1"
	case "general_manager":
		query = "select distinct a.outlet_manager, b.outlet, b.region from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.general_manager= $1"
	case "senior_manager":
		query = "select distinct a.outlet_manager, b.outlet, b.region from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.senior_manager= $1"
	case "manager":
		query = "select distinct a.outlet_manager, b.outlet, b.region from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager= $1"
	case "outlet_manager":
		query = "select emp_id, outlet, region from user_master_table where emp_id= $1"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllOutletManagersReportingToGivenManagerRole: query failed",
			Details:    "FetchAllOutletManagersReportingToGivenManagerRole query failed: " + err.Error(),
		}
	}

	var (
		array = make([]*models.ReportingStaffIDAndName, 0)
	)

	for rows.Next() {

		usr := &models.ReportingStaffIDAndName{}

		if err := rows.Scan(&usr.ID, &usr.Outlet, &usr.Region); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchAllOutletManagersReportingToGivenManagerRole: scan failed failed",
				Details:    "FetchAllOutletManagersReportingToGivenManagerRole scan failed: " + err.Error(),
			}
		}

		array = append(array, usr)
	}

	return array, nil
}

// InsertEmployeeMultipleRecords inserts each users multiple statistical data daily in db, and is updated everyday once
func (dbConn *DBImplementation) InsertEmployeeMultipleRecords(usr []*models.UserRecord) *errors.Error {

	logger.FileLogger.Info().Interface("record", usr).Time("start", time.Now()).Msg("InsertEmployeeMultipleRecords start")
	defer logger.FileLogger.Info().Interface("record", usr).Time("end", time.Now()).Msg("InsertEmployeeMultipleRecords done")

	if len(usr) == 0 {
		return nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}

	count := 0
	for _, w := range usr {
		valueStrings = append(valueStrings, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)", count+1, count+2, count+3, count+4, count+5, count+6, count+7, count+8, count+9, count+10, count+11, count+12, count+13, count+14, count+15, count+16, count+17, count+18, count+19, count+20, count+21, count+22, count+23, count+24, count+25, count+26, count+27, count+28, count+29, count+30, count+31))

		valueArgs = append(valueArgs, w.EmpId)
		valueArgs = append(valueArgs, w.CreateAt)
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

	smt := `INSERT INTO reports_table(emp_id, created, reporting_manager, isquarterly, month, quarter, year, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan)
    VALUES %s ON CONFLICT (emp_id, months, quarter) DO NOTHING`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx, err := dbConn.Operation.Begin()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeMultipleRecords: begin trx failed",
			Details:    "InsertEmployeeMultipleRecords begin trx failed: " + err.Error(),
		}
	}

	// _, err = tx.Exec("Delete FROM reports_table where emp_id=$1", usr[0].EmpId)
	// if err != nil {
	// 	tx.Rollback()
	// 	return &commons.GenericError{Code: 1001, ErrorMessage: "InsertEmployeeMultipleRecords db delete error " + err.Error()}
	// }

	if usr[0].EmpId == "34296" {
		for _, v := range usr {
			fmt.Printf("%+v\n", v)
		}
	}

	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeMultipleRecords: exec failed",
			Details:    "InsertEmployeeMultipleRecords exec failed: " + err.Error(),
		}
	}

	if err := tx.Commit(); err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertEmployeeMultipleRecords: commit failed",
			Details:    "InsertEmployeeMultipleRecords commit failed: " + err.Error(),
		}
	}

	return nil
}

func (dbConn *DBImplementation) InsertHierarchy(usr ...*models.UserHierarchy) *errors.Error {

	logger.FileLogger.Info().Interface("record", usr).Time("start", time.Now()).Msg("InsertHierarchy start")
	defer logger.FileLogger.Info().Interface("record", usr).Time("end", time.Now()).Msg("InsertHierarchy done")

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
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertHierarchy: begin trx failed",
			Details:    "InsertHierarchy begin trx failed: " + err.Error(),
		}
	}

	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertHierarchy: exec failed",
			Details:    "InsertHierarchy exec failed: " + err.Error(),
		}
	}

	if err := tx.Commit(); err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "InsertHierarchy: commit failed",
			Details:    "InsertHierarchy commit failed: " + err.Error(),
		}
	}

	return nil
}

func (dbConn *DBImplementation) FetchRegionRecord(regionCodes []string) ([]*models.GetRegionResponse, *errors.Error) {

	logger.FileLogger.Info().Interface("record", regionCodes).Time("start", time.Now()).Msg("FetchRegionRecord start")
	defer logger.FileLogger.Info().Interface("record", regionCodes).Time("end", time.Now()).Msg("FetchRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT region_code, region_name FROM regions_table"

	if len(regionCodes) != 0 {

		query = "SELECT region_code, region_name FROM regions_table WHERE region_code IN ("

		for _, v := range regionCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchRegionRecord: query failed",
			Details:    "FetchRegionRecord query failed: " + err.Error(),
		}
	}

	var (
		regionCode string
		regionName string
		response   []*models.GetRegionResponse
	)

	for rows.Next() {

		if err := rows.Scan(&regionCode, &regionName); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchRegionRecord: scan failed",
				Details:    "FetchRegionRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.GetRegionResponse{RegionCode: regionCode, RegionName: regionName})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchOutletRecord(outletCodes []string) ([]*models.GetOutletResponse, *errors.Error) {

	logger.FileLogger.Info().Interface("record", outletCodes).Time("start", time.Now()).Msg("FetchOutletRecord start")
	defer logger.FileLogger.Info().Interface("record", outletCodes).Time("end", time.Now()).Msg("FetchOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT outlet_code, outlet_name, region_code, latitude, longitude FROM outlets_table"

	if len(outletCodes) != 0 {

		query = "SELECT outlet_code, outlet_name, region_code, latitude, longitude FROM outlets_table WHERE outlet_code IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletRecord: query failed",
			Details:    "FetchOutletRecord query failed: " + err.Error(),
		}
	}

	var (
		outletCode string
		outletName string
		regionCode string
		latitude   float64
		longitude  float64
		response   []*models.GetOutletResponse
	)

	for rows.Next() {

		if err := rows.Scan(&outletCode, &outletName, &regionCode, &latitude, &longitude); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletRecord: scan failed",
				Details:    "FetchOutletRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.GetOutletResponse{OutletCode: outletCode, OutletName: outletName, RegionCode: regionCode, Lat: latitude, Lang: longitude})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchProductRecord(productCodes []string) ([]*models.GetProductResponse, *errors.Error) {

	logger.FileLogger.Info().Interface("record", productCodes).Time("start", time.Now()).Msg("FetchProductRecord start")
	defer logger.FileLogger.Info().Interface("record", productCodes).Time("end", time.Now()).Msg("FetchProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT product_code, product_name, weightage FROM products_table"

	if len(productCodes) != 0 {

		query = "SELECT product_code, product_name, weightage FROM products_table WHERE product_code IN ("

		for _, v := range productCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchProductRecord: query failed",
			Details:    "FetchProductRecord query failed: " + err.Error(),
		}
	}

	var (
		productCode string
		productName string
		weightage   float64
		response    []*models.GetProductResponse
	)

	for rows.Next() {

		if err := rows.Scan(&productCode, &productName, &weightage); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchProductRecord: scan failed",
				Details:    "FetchProductRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.GetProductResponse{ProductCode: productCode, ProductName: productName, Weightage: weightage})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchPlanRecord(planCodes []string) ([]*models.GetPlanResponse, *errors.Error) {

	logger.FileLogger.Info().Interface("record", planCodes).Time("start", time.Now()).Msg("FetchPlanRecord start")
	defer logger.FileLogger.Info().Interface("record", planCodes).Time("end", time.Now()).Msg("FetchPlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT plan_code, plan_name, product_code, plan_price, sip_plan_rate FROM plans_table"

	if len(planCodes) != 0 {

		query = "SELECT plan_code, plan_name, product_code, plan_price, sip_plan_rate FROM plans_table WHERE plan_code IN ("

		for _, v := range planCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchPlanRecord: query failed",
			Details:    "FetchPlanRecord query failed: " + err.Error(),
		}
	}

	var (
		planCode    string
		planName    string
		productCode string
		planPrice   float64
		sipPlanRate float64
		response    []*models.GetPlanResponse
	)

	for rows.Next() {

		if err := rows.Scan(&planCode, &planName, &productCode, &planPrice, &sipPlanRate); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPlanRecord: scan failed",
				Details:    "FetchPlanRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.GetPlanResponse{PlanCode: planCode, PlanName: planName, ProductCode: productCode, PlanPrice: planPrice, SIPPlanRate: sipPlanRate})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchTargetRecord(targetCodes []string, outletCodes []string, year, month int) (map[string][]*models.GetTargetResponse, *errors.Error) {

	logger.FileLogger.Info().Interface("record", targetCodes).Time("start", time.Now()).Msg("FetchTargetRecord start")
	defer logger.FileLogger.Info().Interface("record", targetCodes).Time("end", time.Now()).Msg("FetchTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table"

	switch {
	case len(targetCodes) != 0:
		query = "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE target_code IN ("

		for _, v := range targetCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	case len(outletCodes) != 0:
		query = "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE outlet_code IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ")"
	case year != 0:
		query = "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE year=" + fmt.Sprintf("%d", year)
	case month != 0:
		query = "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE month=" + fmt.Sprintf("%d", month)
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchTargetRecord: query failed",
			Details:    "FetchTargetRecord query failed: " + err.Error(),
		}
	}

	var (
		targetCode string
		regionCode string
		outletCode string
		months     int
		years      int
		baqati     float64
		hbb        float64
		hayyak     float64
		upgrade    float64
		response   = make(map[string][]*models.GetTargetResponse)
	)

	for rows.Next() {

		if err := rows.Scan(&targetCode, &regionCode, &outletCode, &months, &years, &baqati, &hbb, &hayyak, &upgrade); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchTargetRecord: scan failed",
				Details:    "FetchTargetRecord scan failed: " + err.Error(),
			}
		}

		response[outletCode] = append(response[outletCode], &models.GetTargetResponse{
			TargetCode: targetCode,
			RegionCode: regionCode,
			Month:      months,
			Year:       years,
			Baqati:     baqati,
			HBB:        hbb,
			Hayyak:     hayyak,
			Upgrade:    upgrade,
		})
	}

	for k, v := range response {

		sort.Slice(v, func(i, j int) bool {
			return (v[i].Year > v[j].Year) && (v[i].Month > v[j].Month)
		})

		response[k] = v
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchIndividualUserRecord(id string) ([]*models.UserRecord, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Time("start", time.Now()).Msg("FetchIndividualUserRecord start")
	defer logger.FileLogger.Info().Str("clientID", id).Time("end", time.Now()).Msg("FetchIndividualUserRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select reporting_manager, isquarterly, months, quarter, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from reports_table where emp_id= $1", id)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchIndividualUserRecord: query failed",
			Details:    "FetchIndividualUserRecord query failed: " + err.Error(),
		}
	}

	rcrds := make([]*models.UserRecord, 0)

	for rows.Next() {

		rcrd := &models.UserRecord{
			EmpId: id,
		}

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
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchIndividualUserRecord: scan failed",
				Details:    "FetchIndividualUserRecord scan failed: " + err.Error(),
			}
		}

		rcrds = append(rcrds, rcrd)
	}

	return rcrds, nil
}

func (dbConn *DBImplementation) FetchIndividualUserPayoutReport(id string, month, quarter int, isQuarterly bool) (*models.UserRecord, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("start", time.Now()).Msg("FetchIndividualUserPayoutReport start")
	defer logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("end", time.Now()).Msg("FetchIndividualUserPayoutReport done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_payout, hayyak_extra_payout, hbb_payout, hbb_extra_payout, baqati_payout, baqati_extra_payout, upgrade_payout, upgrade_extra_payout, total_payout from reports_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_payout, hayyak_extra_payout, hbb_payout, hbb_extra_payout, baqati_payout, baqati_extra_payout, upgrade_payout, upgrade_extra_payout, total_payout from reports_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchIndividualUserPayoutReport: query failed",
			Details:    "FetchIndividualUserPayoutReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.UserRecord{
		EmpId: id,
	}

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
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchIndividualUserPayoutReport: scan failed",
				Details:    "FetchIndividualUserPayoutReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserTargetReport(id string, month, quarter int, isQuarterly bool) (*models.UserRecord, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("start", time.Now()).Msg("FetchIndividualUserTargetReport start")
	defer logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("end", time.Now()).Msg("FetchIndividualUserTargetReport done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hbb_sales_count, hbb_target, hbb_achieved_percentage, baqati_sales_count, baqati_target, baqati_achieved_percentage, upgrade_sales_count, upgrade_target, uprade_achieved_percentage from reports_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hbb_sales_count, hbb_target, hbb_achieved_percentage, baqati_sales_count, baqati_target, baqati_achieved_percentage, upgrade_sales_count, upgrade_target, uprade_achieved_percentage from reports_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchIndividualUserTargetReport: query failed",
			Details:    "FetchIndividualUserTargetReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.UserRecord{
		EmpId: id,
	}

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
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchIndividualUserTargetReport: scan failed",
				Details:    "FetchIndividualUserTargetReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserSaleCountPerPlanReport(id string, month, quarter int, isQuarterly bool) (*models.UserRecord, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("start", time.Now()).Msg("FetchIndividualUserSaleCountPerPlanReport start")
	defer logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("end", time.Now()).Msg("FetchIndividualUserSaleCountPerPlanReport done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from reports_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan from reports_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchIndividualUserSaleCountPerPlanReport: query failed",
			Details:    "FetchIndividualUserSaleCountPerPlanReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.UserRecord{EmpId: id}

	for rows.Next() {

		if err := rows.Scan(
			&rcrd.HayyakSaleCountPerPlanArray,
			&rcrd.HBBSaleCountPerPlanArray,
			&rcrd.BaqatiSaleCountPerPlanArray,
			&rcrd.UpgradeSaleCountPerPlanArray,
		); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchIndividualUserSaleCountPerPlanReport: scan failed",
				Details:    "FetchIndividualUserSaleCountPerPlanReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchIndividualUserOverAllReport(id string, month, quarter int, isQuarterly bool) (*models.UserRecord, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("start", time.Now()).Msg("FetchIndividualUserOverAllReport start")
	defer logger.FileLogger.Info().Str("clientID", id).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("end", time.Now()).Msg("FetchIndividualUserOverAllReport done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.QueryContext(ctx, "select reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout from reports_table where emp_id= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select reporting_manager, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout from reports_table where emp_id= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchIndividualUserOverAllReport: query failed",
			Details:    "FetchIndividualUserOverAllReport query failed: " + err.Error(),
		}
	}

	rcrd := &models.UserRecord{
		EmpId: id,
	}

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
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchIndividualUserOverAllReport: scan failed",
				Details:    "FetchIndividualUserOverAllReport scan failed: " + err.Error(),
			}
		}
	}

	return rcrd, nil
}

func (dbConn *DBImplementation) FetchProductSalesPerOutletManager(id string, month, year, quarter int, isQuarterly bool) (*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Int("year", year).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("start", time.Now()).Msg("FetchProductSalesPerOutletManager start")
	defer logger.FileLogger.Info().Str("clientID", id).Int("year", year).Int("month", month).Int("quarter", quarter).Bool("isQuarterReport", isQuarterly).Time("end", time.Now()).Msg("FetchProductSalesPerOutletManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	if !isQuarterly {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count, hbb_sales_count, baqati_sales_count, upgrade_sales_count from reports_table where reporting_manager= $1 and months= $2", id, month)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select hayyak_sales_count, hbb_sales_count, baqati_sales_count, upgrade_sales_count from reports_table where reporting_manager= $1 and quarter= $2", id, quarter)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchProductSalesPerOutletManager: query failed",
			Details:    "FetchProductSalesPerOutletManager query failed: " + err.Error(),
		}
	}

	var (
		hayyak  int
		hbb     int
		baqati  int
		upgrade int
	)

	rcrds := &models.ProductWiseSalesCountPerOutlet{}

	for rows.Next() {

		if err := rows.Scan(&hayyak, &hbb, &baqati, &upgrade); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchProductSalesPerOutletManager: scan failed",
				Details:    "FetchProductSalesPerOutletManager scan failed: " + err.Error(),
			}
		}

		rcrds.Hayyak += hayyak
		rcrds.HBB += hbb
		rcrds.Baqati += baqati
		rcrds.Upgrade += upgrade
	}

	rcrds.Total = rcrds.Hayyak + rcrds.HBB + rcrds.Baqati + rcrds.Upgrade

	return rcrds, nil
}

func (dbConn *DBImplementation) FetchUserName(id string) (string, string, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", id).Time("start", time.Now()).Msg("FetchUserName start")
	defer logger.FileLogger.Info().Str("clientID", id).Time("end", time.Now()).Msg("FetchUserName done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select name, mobile_number from user_master_table where emp_id=$1", id)
	if err != nil {
		return "", "", &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchUserName: query failed",
			Details:    "FetchUserName query failed: " + err.Error(),
		}
	}

	var (
		count  int
		name   string
		mobile string
	)

	for rows.Next() {
		count++

		if err := rows.Scan(&name, &mobile); err != nil {
			return "", "", &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserName: scan failed",
				Details:    "FetchUserName scan failed: " + err.Error(),
			}
		}

		if count > 1 {
			return "", "", &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchUserName: more than one record found",
				Details:    "FetchUserName: more than one record found",
			}
		}
	}

	return name, mobile, nil
}

func (dbConn *DBImplementation) FetchOutletDetails(outletCode, regionCode string) (float64, float64, *errors.Error) {

	logger.FileLogger.Info().Str("outlet", outletCode).Str("region", regionCode).Time("start", time.Now()).Msg("FetchOutletDetails start")
	defer logger.FileLogger.Info().Str("outlet", outletCode).Str("region", regionCode).Time("end", time.Now()).Msg("FetchOutletDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select latitude, longitude from outlets_table where outlet_code=$1 and region_code=$2", outletCode, regionCode)
	if err != nil {
		return 0, 0, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletDetails: query failed",
			Details:    "FetchOutletDetails query failed: " + err.Error(),
		}
	}

	var (
		count     int
		latitude  float64
		longitude float64
	)

	for rows.Next() {
		count++

		if err := rows.Scan(&latitude, &longitude); err != nil {
			return 0, 0, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletDetails: scan failed",
				Details:    "FetchOutletDetails scan failed: " + err.Error(),
			}
		}

		if count > 1 {
			return 0, 0, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletDetails: more than one record found",
				Details:    "FetchOutletDetails: more than one record found",
			}
		}
	}

	return latitude, longitude, nil
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

func (dbConn *DBImplementation) FetchUserBasicDetails(userID string) (*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchUserBasicDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchUserBasicDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.UserBasicDetailsDB{}

	err := dbConn.Operation.QueryRowContext(ctx,
		"select a.emp_id, a.name, a.usr_position, a.mobile_number, b.region_name, b.outlet_name, c.outlet_manager from user_master_table a INNER JOIN outlets_table b on a.emp_id=$1 and a.outlet=b.outlet_name INNER JOIN user_hierarchy_table c on a.emp_id=c.emp_id",
		userID).Scan(&user.UserId, &user.Name, &user.Position, &user.MobileNumber, &user.Region, &user.Outlet, &user.OutletManageId)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "FetchUserBasicDetails query failed: " + err.Error(),
		}
	}

	if user.UserId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    "User Not Found",
			Details:    "FetchUserBasicDetails user not found:" + userID,
		}
	}

	return user, nil
}

func (dbConn *DBImplementation) FetchManagerBasicDetails(userID string) (*models.ManagerBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchManagerBasicDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchManagerBasicDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	manager := &models.ManagerBasicDetailsDB{}

	err := dbConn.Operation.QueryRowContext(ctx,
		"select emp_id, name, usr_position, mobile_number, region, outlet from user_master_table where emp_id=$1",
		userID).Scan(&manager.UserId, &manager.Name, &manager.Position, &manager.MobileNumber, &manager.Region, &manager.Outlet)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "FetchManagerBasicDetails query failed: " + err.Error(),
		}
	}

	if manager.UserId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    "User Not Found",
			Details:    "FetchManagerBasicDetails user not found:" + userID,
		}
	}

	mangerPosition, ok := constant.ManagerMap[manager.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	var query string

	switch mangerPosition {
	case "vice_president":
		return manager, nil
	case "general_manager":
		query = "select vice_president from user_hierarchy_table where general_manager= $1"
		manager.ManagerPosition = constant.VicePresident
	case "senior_manager":
		query = "select general_manager from user_hierarchy_table where senior_manager= $1"
		manager.ManagerPosition = constant.GeneralManager
	case "manager":
		query = "select senior_manager from user_hierarchy_table where manager= $1"
		manager.ManagerPosition = constant.SeniorManager
	case "outlet_manager":
		query = "select manager from user_hierarchy_table where outlet_manager= $1"
		manager.ManagerPosition = constant.Manager
	}

	err = dbConn.Operation.QueryRowContext(ctx, query, userID).Scan(&manager.ManagerId)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    errors.InternalServerError,
			Details:    "FetchManagerBasicDetails manager query failed: " + err.Error(),
		}
	}

	return manager, nil
}

func (dbConn *DBImplementation) FetchOutletSalesCount(region, outlet string, month, year, quarter int, isQuarterReport bool) (hbbSales, hayyakSales, baqatiSales, upgradeSales int, err *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rows   *sql.Rows
		sqlErr error
	)

	if isQuarterReport {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select product_group, SUM(quantity) AS total from sale_transaction_table  where region=$1 and outlet=$2 and quarter=$3 and EXTRACT('year' from sales_date)=$4 GROUP BY product_group",
			region, outlet, quarter, year)
	} else {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select product_group, SUM(quantity) AS total from sale_transaction_table  where region=$1 and outlet=$2 and EXTRACT('month' from b.sales_date)=$3 and EXTRACT('year' from b.sales_date)=$4 GROUP BY product_group",
			region, outlet, month, year)
	}

	if sqlErr != nil {
		err = &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletSalesCount: query failed",
			Details:    "FetchOutletSalesCount query failed: " + sqlErr.Error(),
		}
		return
	}

	var (
		productGroup  string
		quantityTotal int
	)

	for rows.Next() {

		if scanErr := rows.Scan(&productGroup, &quantityTotal); scanErr != nil {
			err = &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletSalesCount: scan failed",
				Details:    "FetchOutletSalesCount scan failed: " + err.Error(),
			}
			return
		}

		switch productGroup {
		case constant.HBBCategory:
			hbbSales = quantityTotal
		case constant.HayyakCategory:
			hayyakSales = quantityTotal
		case constant.BaqatiCategory:
			baqatiSales = quantityTotal
		case constant.UpgradeCategory:
			upgradeSales = quantityTotal
		}
	}

	return
}
