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

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

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
