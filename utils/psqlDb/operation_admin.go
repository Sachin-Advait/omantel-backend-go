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
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (dbConn *DBImplementation) UpsertRegionRecord(userId string, dbReq []*models.CreateRegionRequest) ([]*models.RegionFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if len(dbReq) == 0 {
		return nil, err
	}

	var (
		failureResponse []*models.RegionFailedResponse
		errorCollection = make(map[string]string)
	)

	for _, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[""] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: "", Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			dbErr  error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO regions_table(region_code, region_name, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7);`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.RegionName, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE regions_table SET region_name=$4, modified=$5, comment=$6 WHERE region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.RegionCode, eachRecord.Quarter, eachRecord.Year, eachRecord.RegionName, time.Now(), "update file - "+userId)
		}

		if dbErr != nil {
			errorCollection[eachRecord.RegionCode] = dbErr.Error()
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "REGION RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[eachRecord.RegionCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[eachRecord.RegionCode] = "REGION OPERATION FAILED"
			failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "REGION OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.RegionFailedResponse{RegionCode: eachRecord.RegionCode, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertOutletRecord(userId string, dbReq []*models.CreateOutletRequest) ([]*models.OutletFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if len(dbReq) == 0 {
		return nil, err
	}

	var (
		failureResponse []*models.OutletFailedResponse
		errorCollection = make(map[string]string)
	)

	for _, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[""] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: "", Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var regionName string

		dbErr := dbConn.Operation.QueryRowContext(ctx, "select region_name from regions_table where region_code=$1 and quarter=$2 and year=$3 and deleted IS NULL",
			eachRecord.RegionCode, eachRecord.Quarter, eachRecord.Year).Scan(&regionName)
		if dbErr != nil && dbErr != sql.ErrNoRows {
			errorCollection[eachRecord.OutletCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: errors.InternalServerError})
			continue
		}
		if regionName == "" || eachRecord.RegionName != regionName {
			errorCollection[eachRecord.OutletCode] = "INVALID REGION DETAILS"
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "INVALID REGION DETAILS"})
			continue
		}

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO outlets_table(outlet_code, outlet_name, region_code, region_name, subregion_name, latitude, longitude, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.OutletName, eachRecord.RegionCode, regionName, eachRecord.SubRegionName, eachRecord.Lat,
				eachRecord.Lang, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE outlets_table SET outlet_name=$4, region_code=$5, region_name=$6, subregion_name=$7, latitude=$8, longitude=$9, modified=$10, comment=$11 WHERE outlet_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.OutletCode, eachRecord.Quarter, eachRecord.Year, eachRecord.OutletName, eachRecord.RegionCode,
				eachRecord.RegionName, eachRecord.SubRegionName, eachRecord.Lat, eachRecord.Lang, time.Now(), "update api - "+userId)
		}

		if dbErr != nil {
			errorCollection[eachRecord.OutletCode] = dbErr.Error()
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "OUTLET RECORD EXISTS"})
				continue
			}
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[eachRecord.OutletCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[eachRecord.OutletCode] = "OUTLET OPERATION FAILED"
			failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "OUTLET OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.OutletFailedResponse{OutletCode: eachRecord.OutletCode, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertProductRecord(userId string, dbReq []*models.CreateProductRequest) ([]*models.ProductFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if len(dbReq) == 0 {
		return nil, err
	}

	var (
		failureResponse []*models.ProductFailedResponse
		errorCollection = make(map[string]string)
	)

	for _, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[""] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: "", Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			dbErr  error
			result sql.Result
		)

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO products_table(product_code, product_name, weightage, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.ProductName, eachRecord.Weightage, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE products_table SET product_name=$4, weightage=$5, modified=$6, comment=$7 WHERE product_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.ProductCode, eachRecord.Quarter, eachRecord.Year, eachRecord.ProductName, eachRecord.Weightage, time.Now(), "update file - "+userId)
		}

		if dbErr != nil {
			errorCollection[eachRecord.ProductCode] = dbErr.Error()
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "PRODUCT RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[eachRecord.ProductCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[eachRecord.ProductCode] = "PRODUCT OPERATION FAILED"
			failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "PRODUCT OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.ProductFailedResponse{ProductCode: eachRecord.ProductCode, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertPlanRecord(userId string, dbReq []*models.CreatePlanRequest) ([]*models.PlanFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if len(dbReq) == 0 {
		return nil, err
	}

	var (
		failureResponse []*models.PlanFailedResponse
		errorCollection = make(map[string]string)
	)

	for _, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[""] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: "", Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var productName string

		dbErr := dbConn.Operation.QueryRowContext(ctx, "select product_name from products_table where product_code=$1 and quarter=$2 and year=$3 and deleted IS NULL",
			eachRecord.ProductCode, eachRecord.Quarter, eachRecord.Year).Scan(&productName)
		if dbErr != nil && dbErr != sql.ErrNoRows {
			errorCollection[eachRecord.PlanCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: errors.InternalServerError})
			continue
		}
		if productName == "" {
			errorCollection[eachRecord.PlanCode] = "INVALID PRODUCT DETAILS"
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "INVALID PRODUCT DETAILS"})
			continue
		}

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := `INSERT INTO plans_table(plan_code, plan_name, product_code, plan_price, sip_plan_rate, created, modified, comment, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.PlanName, eachRecord.ProductCode, eachRecord.PlanPrice,
				eachRecord.SIPPlanRate, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE plans_table SET plan_name=$4, product_code=$5, plan_price=$6, sip_plan_rate=$7, modified=$8, comment=$9 WHERE plan_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.PlanCode, eachRecord.Quarter, eachRecord.Year, eachRecord.PlanName, eachRecord.ProductCode,
				eachRecord.PlanPrice, eachRecord.SIPPlanRate, time.Now(), "update file - "+userId)
		}

		if dbErr != nil {
			errorCollection[eachRecord.PlanCode] = dbErr.Error()
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "PLAN RECORD EXISTS"})
				continue
			}
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[eachRecord.PlanCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[eachRecord.PlanCode] = "PLAN OPERATION FAILED"
			failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "PLAN OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.PlanFailedResponse{PlanCode: eachRecord.PlanCode, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertTargetRecord(userId string, dbReq []*models.CreateTargetRequest) ([]*models.TargetFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if len(dbReq) == 0 {
		return nil, err
	}

	var (
		failureResponse []*models.TargetFailedResponse
		errorCollection = make(map[string]string)
	)

	for _, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[""] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: "", Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			regionName string
			outletName string
		)

		dbErr := dbConn.Operation.QueryRowContext(ctx, "select region_name, outlet_name from outlets_table where region_code=$1 and outlet_code=$2 and quarter=$3 and year=$4 and deleted IS NULL",
			eachRecord.RegionName, eachRecord.OutletName, eachRecord.Quarter, eachRecord.Year).Scan(&regionName, &outletName)
		if dbErr != nil && dbErr != sql.ErrNoRows {
			errorCollection[eachRecord.TargetCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: errors.InternalServerError})
			continue
		}
		if regionName == "" || outletName == "" {
			errorCollection[eachRecord.TargetCode] = "INVALID REGION OUTLET DETAILS"
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "INVALID REGION OUTLET DETAILS"})
			continue
		}

		var result sql.Result

		switch eachRecord.IsUpdate {
		case false:
			query := "insert into targets_table(target_code, region, outlet, month, baqati, hbb, hayyak, upgrade, created, modified, comment, quarter, year) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)"
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, regionName, outletName, eachRecord.Month, eachRecord.Baqati, eachRecord.HBB, eachRecord.Hayyak,
				eachRecord.Upgrade, time.Now(), time.Now(), "insert file - "+userId, eachRecord.Quarter, eachRecord.Year)
		default:
			query := `UPDATE targets_table SET region=$4, outlet=$5, baqati=$6, hbb=$7, hayyak=$8, upgrade=$9, modified=$10, comment=$11 WHERE target_code=$1 and quarter=$2 and year=$3 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.TargetCode, eachRecord.Quarter, eachRecord.Year, regionName, outletName, eachRecord.Baqati,
				eachRecord.HBB, eachRecord.Hayyak, eachRecord.Upgrade, time.Now(), "update api - "+userId)
		}

		if dbErr != nil {
			errorCollection[eachRecord.TargetCode] = dbErr.Error()
			if strings.Contains(err.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "TARGET RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[eachRecord.TargetCode] = dbErr.Error()
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[eachRecord.TargetCode] = "TARGET OPERATION FAILED"
			failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "TARGET OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.TargetFailedResponse{TargetCode: eachRecord.TargetCode, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) FetchRegionRecord(regionCodes []string, quarter, year int) ([]*models.RegionResponse, *errors.Error) {

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
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchRegionRecord query failed: " + err.Error(),
		}
	}

	response := make([]*models.RegionResponse, 0)

	for rows.Next() {

		var (
			regionCode string
			regionName string
		)

		if err := rows.Scan(&regionCode, &regionName); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchRegionRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.RegionResponse{RegionCode: regionCode, RegionName: regionName})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchOutletRecord(outletCodes []string, quarter, year int) ([]*models.OutletResponse, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT outlet_code, outlet_name, region_code, region_name, subregion_name, latitude, longitude FROM outlets_table where deleted IS NULL and quarter=$1 and year=$2"

	if len(outletCodes) != 0 {

		query = "SELECT outlet_code, outlet_name, region_code, region_name, subregion_name, latitude, longitude FROM outlets_table WHERE outlet_code IN ("

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
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchOutletRecord query failed: " + err.Error(),
		}
	}

	response := make([]*models.OutletResponse, 0)

	for rows.Next() {

		var (
			outletCode    string
			outletName    string
			regionCode    string
			regionName    string
			subRegionName string
			latitude      float64
			longitude     float64
		)

		if err := rows.Scan(&outletCode, &outletName, &regionCode, &regionName, &subRegionName, &latitude, &longitude); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchOutletRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.OutletResponse{
			OutletCode: outletCode,
			OutletName: outletName,
			RegionCode: regionCode,
			RegionName: regionName,
			SubRegion:  subRegionName,
			Lat:        latitude,
			Lang:       longitude,
		})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchProductRecord(productCodes []string, quarter, year int) ([]*models.ProductResponse, *errors.Error) {

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
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchProductRecord query failed: " + err.Error(),
		}
	}

	response := make([]*models.ProductResponse, 0)

	for rows.Next() {

		var (
			productCode string
			productName string
			weightage   float64
		)

		if err := rows.Scan(&productCode, &productName, &weightage); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
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
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchPlanRecord query failed: " + err.Error(),
		}
	}

	response := make([]*models.PlanResponse, 0)

	for rows.Next() {

		var (
			planCode    string
			planName    string
			productCode string
			planPrice   float64
			sipPlanRate float64
		)

		if err := rows.Scan(&planCode, &planName, &productCode, &planPrice, &sipPlanRate); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchPlanRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &models.PlanResponse{PlanCode: planCode, PlanName: planName, ProductCode: productCode, PlanPrice: planPrice, SIPPlanRate: sipPlanRate})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchTargetRecord(targetCodes []string, outletCodes []string, month, quarter, year int) (map[string][]*models.TargetResponse, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err  error
		rows *sql.Rows
	)

	switch {
	case len(targetCodes) != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE target_code IN ("

		for _, v := range targetCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	case len(outletCodes) != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE outlet IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	case month != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE deleted IS NULL and month=$1 and quarter=$2 and year=$3"
		rows, err = dbConn.Operation.QueryContext(ctx, query, month, quarter, year)

	case quarter != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and quarter=$1 and year=$2"
		rows, err = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	default:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and year=$1"
		rows, err = dbConn.Operation.QueryContext(ctx, query, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchTargetRecord query failed: " + err.Error(),
		}
	}

	var response = make(map[string][]*models.TargetResponse)

	for rows.Next() {

		var (
			targetCode string
			regionName string
			outletName string
			month      int
			year       int
			baqati     float64
			hbb        float64
			hayyak     float64
			upgrade    float64
		)

		if err := rows.Scan(&targetCode, &regionName, &outletName, &month, &year, &baqati, &hbb, &hayyak, &upgrade); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchTargetRecord scan failed: " + err.Error(),
			}
		}

		response[outletName] = append(response[outletName], &models.TargetResponse{
			TargetCode: targetCode,
			RegionCode: regionName,
			Month:      month,
			Year:       year,
			Baqati:     baqati,
			HBB:        hbb,
			Hayyak:     hayyak,
			Upgrade:    upgrade,
		})
	}

	return response, nil
}

func (dbConn *DBImplementation) DeleteRegionRecord(userId string, regionCodes []string, quarter, year int) *errors.Error {

	if len(regionCodes) == 0 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteRegionRecord: invalid db arguments",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE regions_table SET deleted=$1, comment=$2 WHERE region_code IN ("

	for _, v := range regionCodes {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteRegionRecord failed: " + err.Error(),
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteRegionRecord rowsAffected failed: " + err.Error(),
		}
	}

	if updatedRow == 0 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.InvalidReqParameters,
			Details:    "DeleteRegionRecord rowsAffected failed: " + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteOutletRecord(userId string, outletCode []string, quarter, year int) *errors.Error {

	if len(outletCode) == 0 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteOutletRecord: invalid db arguments",
		}
	}

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
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteOutletRecord failed: " + err.Error(),
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteOutletRecord rowsAffected failed: " + err.Error(),
		}
	}

	if updatedRow == 0 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.InvalidReqParameters,
			Details:    "DeleteOutletRecord rowsAffected failed: " + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteProductRecord(userId string, productCode []string, quarter, year int) *errors.Error {

	if len(productCode) == 0 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteProductRecord: invalid db arguments",
		}
	}

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
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteProductRecord failed: " + err.Error(),
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteProductRecord rowsAffected failed: " + err.Error(),
		}
	}

	if updatedRow == 0 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.InvalidReqParameters,
			Details:    "DeleteProductRecord rowsAffected failed: " + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeletePlanRecord(userId string, planCode []string, quarter, year int) *errors.Error {

	if len(planCode) == 0 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeletePlanRecord: invalid db arguments",
		}
	}

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
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeletePlanRecord failed: " + err.Error(),
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeletePlanRecord rowsAffected failed: " + err.Error(),
		}
	}

	if updatedRow == 0 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.InvalidReqParameters,
			Details:    "DeletePlanRecord rowsAffected failed: " + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteTargetRecord(userId string, targetCodes []string, outletCodes []string, month, quarter, year int) *errors.Error {

	if len(targetCodes) == 0 && len(outletCodes) == 0 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteTargetRecord: invalid db arguments",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var query string

	if len(targetCodes) != 0 {
		query = "UPDATE targets_table SET deleted=$1, comment=$2 WHERE target_code IN ("

		for _, v := range targetCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and month=$3 and quarter=$4 and year=$5"
	}

	if len(outletCodes) != 0 {
		query = "UPDATE targets_table SET deleted=$1, comment=$2 WHERE outlet IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and month=$3 and quarter=$4 and year=$5"
	}

	result, err := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, month, quarter, year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteTargetRecord failed: " + err.Error(),
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "DeleteTargetRecord rowsAffected failed: " + err.Error(),
		}
	}

	if updatedRow == 0 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.InvalidReqParameters,
			Details:    "DeleteTargetRecord rowsAffected failed: " + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}
