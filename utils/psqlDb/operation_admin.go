package psqldb

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (dbConn *DBImplementation) FetchRegionRecord(regionCodes []string, quarter, year int) ([]*models.RegionResponse, *errors.Error) {

	// logger.FileLogger.Info().Interface("record", regionCodes).Time("start", time.Now()).Msg("FetchRegionRecord start")
	// defer logger.FileLogger.Info().Interface("record", regionCodes).Time("end", time.Now()).Msg("FetchRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT region_code, region_name FROM regions_table where deleted IS NULL and quarter=$1 and year=$2"

	if len(regionCodes) != 0 {

		query = "SELECT region_code, region_name FROM regions_table WHERE region_code IN ("

		for _, v := range regionCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, quarter, year)
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
		response   []*models.RegionResponse
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

		response = append(response, &models.RegionResponse{RegionCode: regionCode, RegionName: regionName})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchOutletRecord(outletCodes []string, quarter, year int) ([]*models.OutletResponse, *errors.Error) {

	// logger.FileLogger.Info().Interface("record", outletCodes).Time("start", time.Now()).Msg("FetchOutletRecord start")
	// defer logger.FileLogger.Info().Interface("record", outletCodes).Time("end", time.Now()).Msg("FetchOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT outlet_code, outlet_name, region_code, latitude, longitude FROM outlets_table where deleted IS NULL and quarter=$1 and year=$2"

	if len(outletCodes) != 0 {

		query = "SELECT outlet_code, outlet_name, region_code, latitude, longitude FROM outlets_table WHERE outlet_code IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2 "
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, quarter, year)
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
		response   []*models.OutletResponse
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

		response = append(response, &models.OutletResponse{OutletCode: outletCode, OutletName: outletName, RegionCode: regionCode, Lat: latitude, Lang: longitude})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchProductRecord(productCodes []string, quarter, year int) ([]*models.ProductResponse, *errors.Error) {

	// logger.FileLogger.Info().Interface("record", productCodes).Time("start", time.Now()).Msg("FetchProductRecord start")
	// defer logger.FileLogger.Info().Interface("record", productCodes).Time("end", time.Now()).Msg("FetchProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT product_code, product_name, weightage FROM products_table where deleted IS NULL and quarter=$1 and year=$2"

	if len(productCodes) != 0 {

		query = "SELECT product_code, product_name, weightage FROM products_table WHERE product_code IN ("

		for _, v := range productCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, quarter, year)
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
		response    []*models.ProductResponse
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

		response = append(response, &models.ProductResponse{ProductCode: productCode, ProductName: productName, Weightage: weightage})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchPlanRecord(planCodes []string, quarter, year int) ([]*models.PlanResponse, *errors.Error) {

	// logger.FileLogger.Info().Interface("record", planCodes).Time("start", time.Now()).Msg("FetchPlanRecord start")
	// defer logger.FileLogger.Info().Interface("record", planCodes).Time("end", time.Now()).Msg("FetchPlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT plan_code, plan_name, product_code, plan_price, sip_plan_rate FROM plans_table where deleted IS NULL and quarter=$1 and year=$2"

	if len(planCodes) != 0 {

		query = "SELECT plan_code, plan_name, product_code, plan_price, sip_plan_rate FROM plans_table WHERE plan_code IN ("

		for _, v := range planCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, quarter, year)
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
		response    []*models.PlanResponse
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

		response = append(response, &models.PlanResponse{PlanCode: planCode, PlanName: planName, ProductCode: productCode, PlanPrice: planPrice, SIPPlanRate: sipPlanRate})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchTargetRecord(targetCodes []string, outletCodes []string, month, quarter, year int) (map[string][]*models.TargetResponse, *errors.Error) {

	// logger.FileLogger.Info().Interface("record", targetCodes).Time("start", time.Now()).Msg("FetchTargetRecord start")
	// defer logger.FileLogger.Info().Interface("record", targetCodes).Time("end", time.Now()).Msg("FetchTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	switch {
	case len(targetCodes) != 0:
		query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE target_code IN ("

		for _, v := range targetCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL"
		rows, err = dbConn.Operation.QueryContext(ctx, query)

	case len(outletCodes) != 0:
		query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE outlet_code IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	case month != 0:
		query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE deleted IS NULL and month=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, month, year)

	case quarter != 0:
		query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and quarter=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	default:
		query := "SELECT target_code, region_code, outlet_code, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and year=$1"
		rows, err = dbConn.Operation.QueryContext(ctx, query, year)
	}

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
		response   = make(map[string][]*models.TargetResponse)
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

		response[outletCode] = append(response[outletCode], &models.TargetResponse{
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

	return response, nil
}

func (dbConn *DBImplementation) UpsertRegionRecord(userId string, dbReq *models.CreateRegionRequest) *errors.Error {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertRegionRecord start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO regions_table(region_code, region_name, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (region_code, quarter, year) DO NOTHING;`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.RegionCode, dbReq.RegionName, time.Now(), time.Now(), "insert api - "+userId, dbReq.Quarter, dbReq.Year)
	default:
		query := `UPDATE regions_table SET region_name=$4, modified=$5, comment=$6 WHERE region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.RegionCode, dbReq.Quarter, dbReq.Year, dbReq.RegionName, time.Now(), "update api - "+userId)
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

func (dbConn *DBImplementation) UpsertRegionRecordViaFile(userId string, dbReq []*models.CreateRegionRequest) []*models.RegionFailedResponse {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertRegionRecordViaFile start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertRegionRecordViaFile done")

	var failureResponse []*models.RegionFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO regions_table(region_code, region_name, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (region_code, quarter, year) DO NOTHING;`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.RegionName, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE regions_table SET region_name=$4, modified=$5, comment=$6 WHERE region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.Quarter, eachRecord.Year, eachRecord.RegionName, time.Now(), "update file - "+userId)
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "UpsertRegionRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertOutletRecord(userId string, dbReq *models.CreateOutletRequest) *errors.Error {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertOutletRecord start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var regionName string

	rows, err := dbConn.Operation.QueryContext(ctx, "select region_name from regions_table where region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL", dbReq.RegionCode, dbReq.Quarter, dbReq.Year)
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

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var result sql.Result

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO outlets_table(outlet_code, outlet_name, region_code, region_name, latitude, longitude, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (outlet_code, quarter, year) DO NOTHING;`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.OutletCode, dbReq.OutletName, dbReq.RegionCode, regionName, dbReq.Lat, dbReq.Lang, time.Now(), time.Now(), "insert api - "+userId, dbReq.Quarter, dbReq.Year)
	default:
		query := `UPDATE outlets_table SET outlet_name=$4, region_code=$5, latitude=$6, longitude=$7, modified=$8, comment=$9 WHERE outlet_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.OutletCode, dbReq.Quarter, dbReq.Year, dbReq.OutletName, dbReq.RegionCode, dbReq.Lat, dbReq.Lang, time.Now(), "update api - "+userId)
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

func (dbConn *DBImplementation) UpsertOutletRecordViaFile(userId string, dbReq []*models.CreateOutletRequest) []*models.OutletFailedResponse {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertOutletRecordViaFile start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertOutletRecordViaFile done")

	var failureResponse []*models.OutletFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var regionName string

		rows, err := dbConn.Operation.QueryContext(ctx, "select region_name from regions_table where region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL", eachRecord.RegionCode, eachRecord.Quarter, eachRecord.Year)
		if err != nil {
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: fetch regionName query failed: " + err.Error()})
			continue
		}

		if err := rows.Scan(&regionName); err != nil {
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: fetch regionName scan failed: " + err.Error()})
			continue
		}

		if regionName == "" {
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: blank regionName found"})
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO outlets_table(outlet_code, outlet_name, region_code, region_name, latitude, longitude, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (outlet_code, quarter, year) DO NOTHING;`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.OutletName, eachRecord.RegionCode, regionName, eachRecord.Lat, eachRecord.Lang, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE outlets_table SET outlet_name=$4, region_code=$5, latitude=$6, longitude=$7, modified=$8, comment=$9 WHERE outlet_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.Quarter, eachRecord.Year, eachRecord.OutletName, eachRecord.RegionCode, eachRecord.Lat, eachRecord.Lang, time.Now(), "update api - "+userId)
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "UpsertOutletRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertProductRecord(userId string, dbReq *models.CreateProductRequest) *errors.Error {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertProductRecord start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO products_table(product_code, product_name, weightage, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (product_code, quarter, year) DO NOTHING;`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.ProductCode, dbReq.ProductName, dbReq.Weightage, time.Now(), time.Now(), "insert api - "+userId, dbReq.Quarter, dbReq.Year)
	default:
		query := `UPDATE products_table SET product_name=$4, weightage=$5, modified=$6, comment=$7 WHERE product_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.ProductCode, dbReq.Quarter, dbReq.Year, dbReq.ProductName, dbReq.Weightage, time.Now(), "update api - "+userId)
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

func (dbConn *DBImplementation) UpsertProductRecordViaFile(userId string, dbReq []*models.CreateProductRequest) []*models.ProductFailedResponse {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertProductRecordViaFile start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertProductRecordViaFile done")

	var failureResponse []*models.ProductFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO products_table(product_code, product_name, weightage, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (product_code, quarter, year) DO NOTHING;`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.ProductName, eachRecord.Weightage, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE products_table SET product_name=$4, weightage=$5, modified=$6, comment=$7 WHERE product_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.Quarter, eachRecord.Year, eachRecord.ProductName, eachRecord.Weightage, time.Now(), "update file - "+userId)
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "UpsertProductRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertPlanRecord(userId string, dbReq *models.CreatePlanRequest) *errors.Error {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertPlanRecord start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertPlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err    error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO plans_table(plan_code, plan_name, product_code, plan_price, sip_plan_rate, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (plan_code, quarter, year) DO NOTHING;`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.PlanCode, dbReq.PlanName, dbReq.ProductCode, dbReq.PlanPrice, dbReq.SIPPlanRate, time.Now(), time.Now(), "insert api - "+userId, dbReq.Quarter, dbReq.Year)
	default:
		query := `UPDATE plans_table SET plan_name=$4, product_code=$5, plan_price=$6, sip_plan_rate=$7, modified=$8, comment=$9 WHERE plan_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.PlanCode, dbReq.Quarter, dbReq.Year, dbReq.PlanName, dbReq.ProductCode, dbReq.PlanPrice, dbReq.SIPPlanRate, time.Now(), "update api - "+userId)
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

func (dbConn *DBImplementation) UpsertPlanRecordViaFile(userId string, dbReq []*models.CreatePlanRequest) []*models.PlanFailedResponse {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertPlanRecordViaFile start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertPlanRecordViaFile done")

	var failureResponse []*models.PlanFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			err    error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO plans_table(plan_code, plan_name, product_code, plan_price, sip_plan_rate, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (plan_code, quarter, year) DO NOTHING;`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.PlanName, eachRecord.ProductCode, eachRecord.PlanPrice, eachRecord.SIPPlanRate, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE plans_table SET plan_name=$4, product_code=$5, plan_price=$6, sip_plan_rate=$7, modified=$8, comment=$9 WHERE plan_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.Quarter, eachRecord.Year, eachRecord.PlanName, eachRecord.ProductCode, eachRecord.PlanPrice, eachRecord.SIPPlanRate, time.Now(), "update file - "+userId)
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "UpsertPlanRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) UpsertTargetRecord(userId string, dbReq *models.CreateTargetRequest) *errors.Error {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertTargetRecord start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var outletName string

	rows, err := dbConn.Operation.QueryContext(ctx, "select outlet_name from outlets_table where region_code=$1 and outlet_code=$2 quarter=$3 and year=$4 and deleted IS NULL", dbReq.RegionCode, dbReq.OutletCode, dbReq.Quarter, dbReq.Year)
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

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var result sql.Result

	switch dbReq.IsUpdate {
	case false:
		query := "insert into targets_table(target_code, region, outlet, month, baqati, hbb, hayyak, upgrade, created, modified, comment, quarter, year) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT (target_code, quarter, year) DO NOTHING;"
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.TargetCode, dbReq.RegionCode, dbReq.OutletCode, dbReq.Month, dbReq.Baqati, dbReq.HBB, dbReq.Hayyak, dbReq.Upgrade, time.Now(), time.Now(), "insert api - "+userId, dbReq.Quarter, dbReq.Year)
	default:
		query := `UPDATE targets_table SET baqati=$4, hbb=$5, hayyak=$6, upgrade=$7, modified=$8, comment=$9 WHERE target_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
		result, err = dbConn.Operation.ExecContext(ctx, query, dbReq.TargetCode, dbReq.Quarter, dbReq.Year, dbReq.Baqati, dbReq.HBB, dbReq.Hayyak, dbReq.Upgrade, time.Now(), "update api - "+userId)
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

func (dbConn *DBImplementation) UpsertTargetRecordViaFile(userId string, dbReq []*models.CreateTargetRequest) []*models.TargetFailedResponse {

	// logger.FileLogger.Info().Interface("record", dbReq).Time("start", time.Now()).Msg("UpsertTargetRecordViaFile start")
	// defer logger.FileLogger.Info().Interface("record", dbReq).Time("end", time.Now()).Msg("UpsertTargetRecordViaFile done")

	var failureResponse []*models.TargetFailedResponse

	for _, eachRecord := range dbReq {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var outletName string

		rows, err := dbConn.Operation.QueryContext(ctx, "select outlet_name from outlets_table where region_code=$1 and outlet_code=$2 quarter=$3 and year=$4 and deleted IS NULL", eachRecord.RegionCode, eachRecord.OutletCode, eachRecord.Quarter, eachRecord.Year)
		if err != nil {
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: fetch outletName query failed: " + err.Error()})
			continue
		}

		if err := rows.Scan(&outletName); err != nil {
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: fetch outletName scan failed: " + err.Error()})
			continue
		}

		if outletName == "" {
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: blank outletName found"})
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := "insert into targets_table(target_code, region, outlet, month, baqati, hbb, hayyak, upgrade, created, modified, comment, quarter, year) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT (target_code, quarter, year) DO NOTHING;"
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, eachRecord.RegionCode, eachRecord.OutletCode, eachRecord.Month, eachRecord.Baqati, eachRecord.HBB, eachRecord.Hayyak, eachRecord.Upgrade, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE targets_table SET baqati=$4, hbb=$5, hayyak=$6, upgrade=$7, modified=$8, comment=$9 WHERE target_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, err = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, eachRecord.Quarter, eachRecord.Year, eachRecord.Baqati, eachRecord.HBB, eachRecord.Hayyak, eachRecord.Upgrade, time.Now(), "update api - "+userId)
		}

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: duplicate entry"})
				continue
			}

			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: operation failure: " + err.Error()})
			continue
		}

		insertedRow, err := result.RowsAffected()
		if err != nil {
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: rowsAffected failed: " + err.Error()})
			continue
		}

		if insertedRow != 1 {
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "UpsertTargetRecordViaFile: rowsAffected count: " + strconv.Itoa(int(insertedRow))})
			continue
		}

		failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "success"})
	}

	return failureResponse
}

func (dbConn *DBImplementation) DeleteRegionRecord(userId string, regionCodes []string, quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Interface("record", regionCodes).Time("start", time.Now()).Msg("DeleteRegionRecord start")
	// defer logger.FileLogger.Info().Interface("record", regionCodes).Time("end", time.Now()).Msg("DeleteRegionRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE regions_table SET deleted=$, comment=$2 WHERE region_code IN ("

	for _, v := range regionCodes {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
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

func (dbConn *DBImplementation) DeleteOutletRecord(userId string, outletCode []string, quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Interface("record", outletCode).Time("start", time.Now()).Msg("DeleteOutletRecord start")
	// defer logger.FileLogger.Info().Interface("record", outletCode).Time("end", time.Now()).Msg("DeleteOutletRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE outlets_table SET deleted=$1, comment=$2 WHERE outlet_code IN ("

	for _, v := range outletCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
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

func (dbConn *DBImplementation) DeleteProductRecord(userId string, productCode []string, quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Interface("record", productCode).Time("start", time.Now()).Msg("DeleteProductRecord start")
	// defer logger.FileLogger.Info().Interface("record", productCode).Time("end", time.Now()).Msg("DeleteProductRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE products_table SET deleted=$1, comment=$2 WHERE product_code IN ("

	for _, v := range productCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
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

func (dbConn *DBImplementation) DeletePlanRecord(userId string, planCode []string, quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Interface("record", planCode).Time("start", time.Now()).Msg("DeletePlanRecord start")
	// defer logger.FileLogger.Info().Interface("record", planCode).Time("end", time.Now()).Msg("DeletePlanRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE plans_table SET deleted=$1, comment=$2 WHERE plan_code IN ("

	for _, v := range planCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
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

func (dbConn *DBImplementation) DeleteTargetRecord(userId string, targetCode []string, quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Interface("record", targetCode).Time("start", time.Now()).Msg("DeleteTargetRecord start")
	// defer logger.FileLogger.Info().Interface("record", targetCode).Time("end", time.Now()).Msg("DeleteTargetRecord done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE targets_table SET deleted=$1, comment=$2 WHERE target_code IN ("

	for _, v := range targetCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
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
		// return &errors.Error{
		// 	StatusCode: http.StatusInternalServerError,
		// 	Code:       errors.InternalServerError,
		// 	Message:    "DeleteTargetRecord : rowsAffected not 1",
		// 	Details:    "DeleteTargetRecord rowsAffected failed: " + strconv.Itoa(int(insertedRow)),
		// }
	}

	return nil
}
