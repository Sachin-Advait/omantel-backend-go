package psqldb

import (
	"context"
	"database/sql"
	"fmt"
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
	}

	if len(dbReq) == 0 {
		err.Details = "UpsertRegionRecord: invalid db arguments"
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
	}

	if len(dbReq) == 0 {
		err.Details = "UpsertOutletRecord: invalid db arguments"
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
	}

	if len(dbReq) == 0 {
		err.Details = "UpsertProductRecord: invalid db arguments"
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
	}

	if len(dbReq) == 0 {
		err.Details = "UpsertPlanRecord: invalid db arguments"
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
	}

	if len(dbReq) == 0 {
		err.Details = "UpsertTargetRecord: invalid db arguments"
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

func (dbConn *DBImplementation) UpsertVOCCriteriaRecord(userId string, dbReq []*models.CreateVOCCriteriaRequest) ([]*models.VOCCriteriaFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(dbReq) == 0 {
		err.Details = "invalid db arguments"
		return nil, err
	}

	var (
		failureResponse []*models.VOCCriteriaFailedResponse
		errorCollection = make(map[int]string)
	)

	for index, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[index] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			result sql.Result
			dbErr  error
		)

		switch eachRecord.IsUpdate {
		case false:
			query := "insert into voc_criteria_table(unsatisfactory, almost_achieved, achieved, well_achieved, outstanding, year, created, modified, comment) values($1,$2,$3,$4,$5,$6,$7,$8,$9)"
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.Unsatisfactory, eachRecord.AlmostAchieved, eachRecord.Achieved, eachRecord.WellAchieved,
				eachRecord.Outstanding, eachRecord.Year, time.Now(), time.Now(), "insert file - "+userId)
		default:
			query := `UPDATE voc_criteria_table SET unsatisfactory=$2, almost_achieved=$3, achieved=$4, well_achieved=$5, outstanding=$6, modified=$7, comment=$8 WHERE year=$1 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.Year, eachRecord.Unsatisfactory, eachRecord.AlmostAchieved, eachRecord.Achieved,
				eachRecord.WellAchieved, eachRecord.Outstanding, time.Now(), "update api - "+userId)
		}

		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("query err - %v", dbErr.Error())
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Year: eachRecord.Year, Message: "VOC-CRITERIA RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Year: eachRecord.Year, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("rowsAffected err - %v", dbErr.Error())
			failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Year: eachRecord.Year, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[index] = "VOC-CRITERIA OPERATION FAILED"
			failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Year: eachRecord.Year, Message: "VOC-CRITERIA OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.VOCCriteriaFailedResponse{Year: eachRecord.Year, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertVOCMaxPayRecord(userId string, dbReq []*models.CreateVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(dbReq) == 0 {
		err.Details = "invalid db arguments"
		return nil, err
	}

	var (
		failureResponse []*models.VOCMaxPayFailedResponse
		errorCollection = make(map[int]string)
	)

	for index, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[index] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			result sql.Result
			dbErr  error
		)

		switch eachRecord.IsUpdate {
		case false:
			query := "insert into voc_max_pay_table(position, max_pay, year, created, modified, comment) values($1,$2,$3,$4,$5,$6)"
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.Position, eachRecord.MaxPay, eachRecord.Year, time.Now(), time.Now(), "insert file - "+userId)
		default:
			query := `UPDATE voc_max_pay_table SET max_pay=$3, modified=$4, comment=$5 WHERE position=$1 and year=$2 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.Position, eachRecord.Year, eachRecord.MaxPay, time.Now(), "update api - "+userId)
		}

		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("query err - %v", dbErr.Error())
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: eachRecord.Year, Position: eachRecord.Position, Message: "VOC-MaxPay RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: eachRecord.Year, Position: eachRecord.Position, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("rowsAffected err - %v", dbErr.Error())
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: eachRecord.Year, Position: eachRecord.Position, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[index] = "VOC-MaxPay OPERATION FAILED"
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: eachRecord.Year, Position: eachRecord.Position, Message: "VOC-MaxPay OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: eachRecord.Year, Position: eachRecord.Position, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) UpsertVOCUserRecord(userId string, dbReq []*models.CreateVOCUserRequest) ([]*models.VOCCUserFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(dbReq) == 0 {
		err.Details = "invalid db arguments"
		return nil, err
	}

	var (
		failureResponse []*models.VOCCUserFailedResponse
		errorCollection = make(map[int]string)
	)

	for index, eachRecord := range dbReq {

		if eachRecord == nil {
			errorCollection[index] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var (
			result sql.Result
			dbErr  error
		)

		switch eachRecord.IsUpdate {
		case false:
			query := "insert into voc_user_table(emp_id, achieved_percent, year, created, modified, comment) values($1,$2,$3,$4,$5,$6)"
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.UserId, eachRecord.AchievedPercent, eachRecord.Year, time.Now(), time.Now(), "insert file - "+userId)
		default:
			query := `UPDATE voc_user_table SET achieved_percent=$3, modified=$4, comment=$5 WHERE emp_id=$1 and year=$2 and deleted IS NULL`
			result, dbErr = dbConn.Operation.ExecContext(ctx, query, eachRecord.UserId, eachRecord.Year, eachRecord.AchievedPercent, time.Now(), "update api - "+userId)
		}

		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("query err - %v", dbErr.Error())
			if strings.Contains(dbErr.Error(), "duplicate") {
				failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{UserId: eachRecord.UserId, Year: eachRecord.Year, Message: "VOC-USER RECORD EXISTS"})
				continue
			}

			failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{UserId: eachRecord.UserId, Year: eachRecord.Year, Message: errors.InternalServerError})
			continue
		}

		insertedRow, dbErr := result.RowsAffected()
		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("rowsAffected err - %v", dbErr.Error())
			failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{UserId: eachRecord.UserId, Year: eachRecord.Year, Message: errors.InternalServerError})
			continue
		}

		if insertedRow != 1 {
			errorCollection[index] = "VOC-USER OPERATION FAILED"
			failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{UserId: eachRecord.UserId, Year: eachRecord.Year, Message: "VOC-USER OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.VOCCUserFailedResponse{UserId: eachRecord.UserId, Year: eachRecord.Year, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) FetchRegionRecord(regionCodes []string, quarter, year int) ([]*models.RegionResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

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

	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, quarter, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchRegionRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchRegionRecord: query failed-" + dbErr.Error()
		return nil, err
	}

	response := make([]*models.RegionResponse, 0)

	for rows.Next() {

		var (
			regionCode string
			regionName string
		)

		if dbErr := rows.Scan(&regionCode, &regionName); dbErr != nil {
			err.Details = "FetchRegionRecord: scan failed-" + dbErr.Error()
			return nil, err
		}

		response = append(response, &models.RegionResponse{RegionCode: regionCode, RegionName: regionName})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchOutletRecord(outletCodes []string, quarter, year int) ([]*models.OutletResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

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

	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, quarter, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchOutletRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchOutletRecord: query failed-" + dbErr.Error()
		return nil, err
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

		if dbErr := rows.Scan(&outletCode, &outletName, &regionCode, &regionName, &subRegionName, &latitude, &longitude); dbErr != nil {
			err.Details = "FetchOutletRecord: scan failed-" + dbErr.Error()
			return nil, err
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

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

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

	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, quarter, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchProductRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchProductRecord: query failed-" + dbErr.Error()
		return nil, err
	}

	response := make([]*models.ProductResponse, 0)

	for rows.Next() {

		var (
			productCode string
			productName string
			weightage   float64
		)

		if dbErr := rows.Scan(&productCode, &productName, &weightage); dbErr != nil {
			err.Details = "FetchProductRecord: scan failed-" + dbErr.Error()
			return nil, err
		}

		response = append(response, &models.ProductResponse{ProductCode: productCode, ProductName: productName, Weightage: weightage})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchPlanRecord(planCodes []string, quarter, year int) ([]*models.PlanResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

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

	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, quarter, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchPlanRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchPlanRecord: query failed-" + dbErr.Error()
		return nil, err
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

		if dbErr := rows.Scan(&planCode, &planName, &productCode, &planPrice, &sipPlanRate); dbErr != nil {
			err.Details = "FetchPlanRecord: scan failed-" + dbErr.Error()
			return nil, err
		}

		response = append(response, &models.PlanResponse{PlanCode: planCode, PlanName: planName, ProductCode: productCode, PlanPrice: planPrice, SIPPlanRate: sipPlanRate})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchTargetRecord(targetCodes []string, outletCodes []string, month, quarter, year int) (map[string][]*models.TargetResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		dbErr error
		rows  *sql.Rows
	)

	switch {
	case len(targetCodes) != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE target_code IN ("

		for _, v := range targetCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
		rows, dbErr = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	case len(outletCodes) != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE outlet IN ("

		for _, v := range outletCodes {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and quarter=$1 and year=$2"
		rows, dbErr = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	case month != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table WHERE deleted IS NULL and month=$1 and quarter=$2 and year=$3"
		rows, dbErr = dbConn.Operation.QueryContext(ctx, query, month, quarter, year)

	case quarter != 0:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and quarter=$1 and year=$2"
		rows, dbErr = dbConn.Operation.QueryContext(ctx, query, quarter, year)

	default:
		query := "SELECT target_code, region, outlet, month, year, baqati, hbb, hayyak, upgrade FROM targets_table where deleted IS NULL and year=$1"
		rows, dbErr = dbConn.Operation.QueryContext(ctx, query, year)
	}

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchTargetRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchTargetRecord: query failed-" + dbErr.Error()
		return nil, err
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

		if dbErr := rows.Scan(&targetCode, &regionName, &outletName, &month, &year, &baqati, &hbb, &hayyak, &upgrade); dbErr != nil {
			err.Details = "FetchTargetRecord: scan failed-" + dbErr.Error()
			return nil, err
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

func (dbConn *DBImplementation) FetchVOCCriteriaRecord(year int) (*models.VOCCriteriaResponse, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var response models.VOCCriteriaResponse

	query := "SELECT unsatisfactory, almost_achieved, achieved, well_achieved, outstanding FROM voc_criteria_table where deleted IS NULL and year=$1"
	err := dbConn.Operation.QueryRowContext(ctx, query, year).Scan(&response.Unsatisfactory, &response.AlmostAchieved, &response.Achieved,
		&response.WellAchieved, &response.Outstanding)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchVOCCriteriaRecord: query failed-" + err.Error(),
			}
		}

		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Err:        errors.InternalServerError,
			Details:    "FetchVOCCriteriaRecord: query failed-" + err.Error(),
		}
	}

	return &response, nil
}

func (dbConn *DBImplementation) FetchVOCMaxPayRecord(year int) ([]*models.VOCMaxPayResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT position, max_pay FROM voc_max_pay_table where deleted IS NULL and year=$1"
	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchVOCMaxPayRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchVOCMaxPayRecord: query failed-" + dbErr.Error()
		return nil, err
	}

	response := make([]*models.VOCMaxPayResponse, 0)

	for rows.Next() {

		var (
			position string
			maxPay   float64
		)

		if dbErr := rows.Scan(&position, &maxPay); dbErr != nil {
			err.Details = "FetchVOCMaxPayRecord: scan failed-" + dbErr.Error()
			return nil, err
		}

		response = append(response, &models.VOCMaxPayResponse{Position: position, MaxPay: maxPay})
	}

	return response, nil
}

func (dbConn *DBImplementation) FetchVOCUserRecord(userIds []string, year int) ([]*models.VOCUserResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT emp_id, achieved_percent FROM voc_user_table where deleted IS NULL and year=$1"

	if len(userIds) != 0 {

		query = "SELECT emp_id, achieved_percent FROM voc_user_table WHERE emp_id IN ("

		for _, v := range userIds {
			query += "'" + v + "', "
		}

		query = strings.TrimSuffix(query, ", ")
		query += ") and deleted IS NULL and year=$1"
	}

	rows, dbErr := dbConn.Operation.QueryContext(ctx, query, year)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchVOCUserRecord: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchVOCUserRecord: query failed-" + dbErr.Error()
		return nil, err
	}

	response := make([]*models.VOCUserResponse, 0)

	for rows.Next() {

		var (
			empId           string
			achievedPercent float64
		)

		if dbErr := rows.Scan(&empId, &achievedPercent); dbErr != nil {
			err.Details = "FetchVOCUserRecord: scan failed-" + dbErr.Error()
			return nil, err
		}

		response = append(response, &models.VOCUserResponse{UserId: empId, AchievedPercent: achievedPercent})
	}

	return response, nil
}

func (dbConn *DBImplementation) DeleteRegionRecord(userId string, regionCodes []string, quarter, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(regionCodes) == 0 {
		err.Details = "DeleteRegionRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE regions_table SET deleted=$1, comment=$2 WHERE region_code IN ("

	for _, v := range regionCodes {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
	if dbErr != nil {
		err.Details = "DeleteRegionRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteRegionRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow == 0 {
		err.Details = "DeleteRegionRecord: row count-" + strconv.Itoa(int(updatedRow))
		return err
	}

	return nil
}

func (dbConn *DBImplementation) DeleteOutletRecord(userId string, outletCode []string, quarter, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if len(outletCode) == 0 {
		err.Details = "DeleteOutletRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE outlets_table SET deleted=$1, comment=$2 WHERE outlet_code IN ("

	for _, v := range outletCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
	if dbErr != nil {
		err.Details = "DeleteOutletRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteOutletRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow == 0 {
		err.Details = "DeleteOutletRecord: row count-" + strconv.Itoa(int(updatedRow))
		return err
	}

	return nil
}

func (dbConn *DBImplementation) DeleteProductRecord(userId string, productCode []string, quarter, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if len(productCode) == 0 {
		err.Details = "DeleteProductRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE products_table SET deleted=$1, comment=$2 WHERE product_code IN ("

	for _, v := range productCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
	if dbErr != nil {
		err.Details = "DeleteProductRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteProductRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow == 0 {
		err.Details = "DeleteProductRecord: row count-" + strconv.Itoa(int(updatedRow))
		return err
	}

	return nil
}

func (dbConn *DBImplementation) DeletePlanRecord(userId string, planCode []string, quarter, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if len(planCode) == 0 {
		err.Details = "DeletePlanRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE plans_table SET deleted=$1, comment=$2 WHERE plan_code IN ("

	for _, v := range planCode {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and quarter=$3 and year=$4"

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, quarter, year)
	if dbErr != nil {
		err.Details = "DeletePlanRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeletePlanRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow == 0 {
		err.Details = "DeletePlanRecord: row count-" + strconv.Itoa(int(updatedRow))
		return err
	}

	return nil
}

func (dbConn *DBImplementation) DeleteTargetRecord(userId string, targetCodes []string, outletCodes []string, month, quarter, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if len(targetCodes) == 0 && len(outletCodes) == 0 {
		err.Details = "DeleteTargetRecord: invalid db arguments"
		return err
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

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, time.Now(), "del - "+userId, month, quarter, year)
	if dbErr != nil {
		err.Details = "DeleteTargetRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteTargetRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow == 0 {
		err.Details = "DeleteTargetRecord: row count-" + strconv.Itoa(int(updatedRow))
		return err
	}

	return nil
}

func (dbConn *DBImplementation) DeleteVOCCriteriaRecord(userId string, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if year == 0 {
		err.Details = "DeleteVOCCriteriaRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE voc_criteria_table SET deleted=$2, comment=$3 WHERE year=$1 and deleted IS NULL"
	result, dbErr := dbConn.Operation.ExecContext(ctx, query, year, time.Now(), "del - "+userId)
	if dbErr != nil {
		err.Details = "DeleteVOCCriteriaRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteVOCCriteriaRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if updatedRow != 1 {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Err:        errors.InvalidReqParameters,
			Message:    "NO RECORDS FOUND",
			Details:    "DeleteVOCCriteriaRecord: row count-" + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}

func (dbConn *DBImplementation) DeleteVOCMaxPayRecord(userId string, positions []string, year []int) ([]*models.VOCMaxPayFailedResponse, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(positions) == 0 || len(positions) != len(year) {
		err.Details = "DeleteVOCMaxPayRecord: invalid db arguments"
		return nil, err
	}

	var (
		failureResponse []*models.VOCMaxPayFailedResponse
		errorCollection = make(map[int]string)
	)

	for index, pos := range positions {

		if pos == "" || year[index] == 0 {
			errorCollection[index] = errors.InternalServerError
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: year[index], Position: pos, Message: errors.InternalServerError})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		query := "UPDATE voc_max_pay_table SET deleted=$3, comment=$4 WHERE position=$1 and year=$2 and deleted IS NULL"

		result, dbErr := dbConn.Operation.ExecContext(ctx, query, pos, year[index], time.Now(), "del - "+userId)
		if dbErr != nil {
			errorCollection[index] = fmt.Sprintf("query err - %v", dbErr.Error())
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: year[index], Position: pos, Message: errors.InternalServerError})
			continue
		}

		updatedRow, err := result.RowsAffected()
		if err != nil {
			errorCollection[index] = fmt.Sprintf("rowsAffected err - %v", dbErr.Error())
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: year[index], Position: pos, Message: errors.InternalServerError})
			continue
		}

		if updatedRow != 1 {
			errorCollection[index] = "VOC-MaxPay OPERATION FAILED"
			failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: year[index], Position: pos, Message: "VOC-MaxPay OPERATION FAILED"})
			continue
		}

		failureResponse = append(failureResponse, &models.VOCMaxPayFailedResponse{Year: year[index], Position: pos, Message: "success"})
	}

	if len(errorCollection) > 0 {
		errVal, _ := objecthelper.MarshalJSON(&errorCollection)
		err.Details = string(errVal)
		return failureResponse, err
	}

	return failureResponse, nil
}

func (dbConn *DBImplementation) DeleteVOCUserRecord(userId string, agentIds []string, year int) *errors.Error {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Err:        errors.InternalServerError,
	}

	if len(agentIds) == 0 {
		err.Details = "DeleteVOCUserRecord: invalid db arguments"
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "UPDATE voc_user_table SET deleted=$2, comment=$3 WHERE emp_id IN ("

	for _, v := range agentIds {
		query += "'" + v + "', "
	}

	query = strings.TrimSuffix(query, ", ")
	query += ") and year=$1 and deleted IS NULL"

	result, dbErr := dbConn.Operation.ExecContext(ctx, query, year, time.Now(), "del - "+userId)
	if dbErr != nil {
		err.Details = "DeleteVOCUserRecord: query failed-" + dbErr.Error()
		return err
	}

	updatedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = "DeleteVOCUserRecord: rowsAffected failed-" + dbErr.Error()
		return err
	}

	if int(updatedRow) != len(agentIds) {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Err:        errors.InvalidReqParameters,
			Message:    "FEW RECORDS NOT FOUND",
			Details:    "DeleteVOCUserRecord: row count-" + strconv.Itoa(int(updatedRow)),
		}
	}

	return nil
}
