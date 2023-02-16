package psqldb

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

/////////////////////// OUTLET TABLE ///////////////////////////

func (dbConn *DBImplementation) FetchOutletDetailsViaName(regionName, outletName string, quarter, year int) (*models.OutletResponse, *errors.Error) {

	// logger.FileLogger.Info().Str("region", regionName).Str("outlet", outletName).Time("start", time.Now()).Msg("FetchOutletDetailsViaName start")
	// defer logger.FileLogger.Info().Str("region", regionName).Str("outlet", outletName).Time("end", time.Now()).Msg("FetchOutletDetailsViaName done")

	if regionName == "" || outletName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "outlet or region name missing",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	outlet := &models.OutletResponse{}

	err := dbConn.Operation.QueryRowContext(ctx,
		"SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table where deleted IS NULL and region_name=$1 and outlet_name=$2 and quarter=$3 and year=$4",
		regionName, outletName, quarter, year).Scan(&outlet.OutletCode, &outlet.OutletName, &outlet.RegionCode, &outlet.RegionName, &outlet.Lat, &outlet.Lang)
	if err == nil {
		return outlet, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    "OUTLATE NAME NOT FOUND",
			Details:    "FetchOutletDetailsViaName: region-" + regionName + ", outlet-" + outletName + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchOutletDetailsViaName: region-" + regionName + ", outlet-" + outletName + ", err - " + err.Error(),
	}
}

func (dbConn *DBImplementation) FetchAllOutletNames(regionName string, quarter, year int) ([]string, *errors.Error) {

	// logger.FileLogger.Info().Str("region", regionName).Time("start", time.Now()).Msg("FetchAllOutletNames start")
	// defer logger.FileLogger.Info().Str("region", regionName).Time("end", time.Now()).Msg("FetchAllOutletNames done")

	if regionName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "outlet name missing",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	outlets := make([]string, 0)

	rows, err := dbConn.Operation.QueryContext(ctx, "SELECT outlet_name FROM outlets_table where deleted IS NULL and region_name=$1 and quarter=$2 and year=$3", regionName, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchAllOutletNames: query failed",
			Details:    "FetchAllOutletNames: query failed - " + err.Error(),
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
				Details:    "FetchAllOutletNames: scan failed - " + err.Error(),
			}
		}

		outlets = append(outlets, outlentName)
	}

	return outlets, nil
}

/////////////////////// SALES TRANSACTION TABLE ///////////////////////////

func (dbConn *DBImplementation) FetchOutletSalesCount(region, outlet string, isQuarterReport bool, quarter, month, year int) (hbbSales, hayyakSales, baqatiSales, upgradeSales int, err *errors.Error) {

	// logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("start", time.Now()).Msg("FetchOutletSalesCount start")
	// defer logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("end", time.Now()).Msg("FetchOutletSalesCount done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		rows   *sql.Rows
		sqlErr error
	)

	if isQuarterReport {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select transaction_type, product_group, SUM(quantity) AS total from sale_transaction_table where region=$1 and outlet=$2 and quarter=$3 and EXTRACT('year' from sales_date)=$4 GROUP BY transaction_type, product_group",
			region, outlet, quarter, year)
	} else {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select transaction_type, product_group, SUM(quantity) AS total from sale_transaction_table where region=$1 and outlet=$2 and EXTRACT('month' from sales_date)=$3 and EXTRACT('year' from sales_date)=$4 GROUP BY transaction_type, product_group",
			region, outlet, month, year)
	}

	if sqlErr != nil {
		err = &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchOutletSalesCount: query failed",
			Details:    "FetchOutletSalesCount: region-" + region + ", outlet-" + outlet + ", query failed - " + sqlErr.Error(),
		}
		return
	}

	for rows.Next() {

		var (
			transactionType string
			productGroup    string
			quantityTotal   int
		)

		if scanErr := rows.Scan(&transactionType, &productGroup, &quantityTotal); scanErr != nil {
			err = &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletSalesCount: scan failed",
				Details:    "FetchOutletSalesCount: region-" + region + ", outlet-" + outlet + ", scan failed - " + err.Error(),
			}
			return
		}

		switch {
		case transactionType == constant.UpgradeCategory:
			upgradeSales += quantityTotal
		case productGroup == constant.HBBCategory || strings.Contains(productGroup, "FIXED"):
			hbbSales += quantityTotal
		case productGroup == constant.HayyakCategory || strings.Contains(productGroup, "PREPAID"):
			hayyakSales += quantityTotal
		case productGroup == constant.BaqatiCategory || strings.Contains(productGroup, "POSTPAID"):
			baqatiSales += quantityTotal
		}
	}

	return
}

// FetchActiveUserSales - all active sale transaction of a given user and time period
func (dbConn *DBImplementation) FetchActiveUserSales(userID string, region string, outlet string, isQuarterReport bool, quarter, month, year int) (*models.SaleTransaction, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Int("month", month).Int("year", year).Int("quarter", quarter).Bool("isQuarterReport", isQuarterReport).Time("start", time.Now()).Msg("FetchUserSales start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Int("month", month).Int("year", year).Int("quarter", quarter).Bool("isQuarterReport", isQuarterReport).Time("end", time.Now()).Msg("FetchUserSales done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		rows *sql.Rows
		err  error
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select id, prev_plan, curr_plan, quantity, product_group, transaction_type, subcategory, sales_date, terminate_date from sale_transaction_table where emp_id=$1 and quarter=$2 and EXTRACT('year' from sales_date)=$3",
			userID, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select id, prev_plan, curr_plan, quantity, product_group, transaction_type, subcategory, sales_date, terminate_date from sale_transaction_table where emp_id=$1 and EXTRACT('month' from sales_date)=$2 and EXTRACT('year' from sales_date)=$3",
			userID, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchUserSales: userID-" + userID + ", query failed - " + err.Error(),
		}
	}

	trx := &models.SaleTransaction{
		EmpID:  userID,
		Region: region,
		Outlet: outlet,
	}

	trx.SaleRecords = make(map[string][]*models.SaleRecord)

	count := 0

	for rows.Next() {

		var (
			id              int
			prevPlan        string
			currPlan        string
			quantity        int
			productGroup    string
			transactionType string
			subCategory     interface{}
			salesDate       time.Time
			terminateDate   time.Time
		)

		count++

		if err := rows.Scan(&id, &prevPlan, &currPlan, &quantity, &productGroup, &transactionType, &subCategory, &salesDate, &terminateDate); err != nil {
			return trx, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchUserSales: " + userID + ", scan failed - " + err.Error(),
			}
		}

		switch {
		case strings.Contains(productGroup, "FIXED") || productGroup == constant.HBBCategory:
			productGroup = constant.HBBCategory

		case strings.Contains(productGroup, "PREPAID") || productGroup == constant.HayyakCategory:
			productGroup = constant.HayyakCategory

		case strings.Contains(productGroup, "POSTPAID") || productGroup == constant.BaqatiCategory:
			productGroup = constant.BaqatiCategory

		default:
			logger.FileLogger.CtxWarn(ctx).Str("userID", userID).Str("region", region).Str("outlet", outlet).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", id).Msg("FetchActiveUserSales: incorrect productGroup")
			// return trx, &errors.Error{
			// 	StatusCode: http.StatusInternalServerError,
			// 	Code:       errors.InternalServerErrorCode,
			// 	Message:    errors.InternalServerError,
			// 	Details:    "FetchUserSales: invalid productGrp - " + productGroup,
			// }
		}

		if !(productGroup == constant.HBBCategory || productGroup == constant.HayyakCategory || productGroup == constant.BaqatiCategory) {
			continue
		}

		if terminateDate.Before(salesDate) {
			logger.FileLogger.CtxError(ctx).Str("userID", userID).Str("region", region).Str("outlet", outlet).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", id).Msg("FetchActiveUserSales: saleDate-terminateDate incorrect")
		}

		expiryDate := salesDate.Add(90 * 24 * time.Hour)

		if terminateDate.Sub(expiryDate) <= 0 { //Months<3
			if !strings.Contains(subCategory.(string), "INVOLUNTARY DEACTIVATION") {
				logger.FileLogger.CtxError(ctx).Str("userID", userID).Str("region", region).Str("outlet", outlet).Bool("isQuarter", isQuarterReport).
					Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", id).Msg("FetchActiveUserSales: sale terminated")
				continue
			}
		}

		rcrd := &models.SaleRecord{
			TrxId:           id,
			Quantity:        quantity,
			Category:        transactionType,
			PreviousProduct: prevPlan,
			CurrentProduct:  currPlan,
		}

		if trx.SaleRecords[productGroup] == nil {
			trx.SaleRecords[productGroup] = make([]*models.SaleRecord, 0)
		}

		trx.SaleRecords[productGroup] = append(trx.SaleRecords[productGroup], rcrd)
	}

	return trx, nil
}

func (dbConn *DBImplementation) FetchAgentUserHierarchy(userID string, quarter, year int) (*models.UserHierarchy, *errors.Error) {

	// logger.FileLogger.Info().Str("userID", userID).Time("start", time.Now()).Msg("FetchAgentUserHierarchy start")
	// defer logger.FileLogger.Info().Str("userID", userID).Time("end", time.Now()).Msg("FetchAgentUserHierarchy done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rcrd := &models.UserHierarchy{}

	err := dbConn.Operation.QueryRowContext(ctx, "select b.outlet_manager, b.manager, b.senior_manager, b.general_manager, b.vice_president from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.emp_id AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND a.emp_id=$1 AND a.quarter=$2 AND a.year=$3 AND a.isactive=true",
		userID, quarter, year).Scan(&rcrd.OutletManager, &rcrd.Manager, &rcrd.SeniorManager, &rcrd.GeneralManager, &rcrd.VicePresident)
	if err == nil {
		return rcrd, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    "OUTLATE NAME NOT FOUND",
			Details:    "FetchAgentUserHierarchy: userId-" + userID + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchAgentUserHierarchy: userId-" + userID + ", err - " + err.Error(),
	}
}

// FetchAllAgentsUserOfOM - all outlet employees reporting to given outlet manager in argument
func (dbConn *DBImplementation) FetchAllAgentsUserOfOM(outletManagerID string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("userID", userID).Time("start", time.Now()).Msg("FetchAgentUserHierarchy start")
	// defer logger.FileLogger.Info().Str("userID", userID).Time("end", time.Now()).Msg("FetchAgentUserHierarchy done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	//rows, err := dbConn.Operation.QueryContext(ctx, "select distinct emp_id from user_hierarchy_table where deleted IS NULL and outlet_manager=$1 and team_lead=$2 and quarter=$3 and year=$4", outletManagerID, outletManagerID, quarter, year)
	rows, err := dbConn.Operation.QueryContext(ctx, "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.emp_id AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.outlet_manager=$1 AND a.quarter=$2 AND a.year=$3", // AND a.isactive=true
		outletManagerID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAllAgentsUserOfOM: query failed - " + err.Error(),
		}
	}

	var (
		userId  string
		userIds = make([]string, 0)
	)

	for rows.Next() {

		if err := rows.Scan(&userId); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllAgentsUserOfOM: scan failed - " + err.Error(),
			}
		}

		if userId != "" && userId != outletManagerID {
			userIds = append(userIds, userId)
		}

		userId = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachUser := range userIds {
		userRecord, err := dbConn.FetchUserBasicDetails(eachUser, quarter, year)
		if err != nil {
			return nil, err
		}

		userRecord.ManagerId = outletManagerID
		userRecord.ManagerPosition = constant.OutletManager

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllOMReportingToGivenManager - all outlet managers reporting to given manager in argument
func (dbConn *DBImplementation) FetchAllOMReportingToGivenManager(managerPosition string, managerID string, quarter, year int) (map[string]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllOMReportingToGivenManager start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllOMReportingToGivenManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.vice_president=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.GeneralManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.general_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.SeniorManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.senior_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.Manager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "NOT A MANAGER",
			Details:    "FetchAllOMReportingToGivenManager: incorrect position - " + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAllOMReportingToGivenManager: query failed - " + err.Error(),
		}
	}

	var (
		outletManager string
		response      = make(map[string]*models.UserBasicDetailsDB, 0)
	)

	for rows.Next() {

		if err := rows.Scan(&outletManager); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllOMReportingToGivenManager: scan failed - " + err.Error(),
			}
		}

		if outletManager == managerID {
			continue
		}

		if response[outletManager] == nil {
			response[outletManager] = &models.UserBasicDetailsDB{}
		}

		// if response[teamLead] == nil {
		// 	response[teamLead] = &models.UserBasicDetailsDB{}
		// }
	}

	for eachUser := range response {
		userRecord, err := dbConn.FetchUserDetails(eachUser, quarter, year)
		if err != nil {
			return nil, err
		}

		response[eachUser] = userRecord
	}

	return response, nil
}

// FetchAllRMReportingToGivenManager - all regional managers reporting to given manager in argument
func (dbConn *DBImplementation) FetchAllRMReportingToGivenManager(managerPosition string, managerID string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllManagerReportingToGivenManager start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllManagerReportingToGivenManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.vice_president=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.GeneralManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.general_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.SeniorManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.senior_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "NOT A MANAGER",
			Details:    "FetchAllManagerReportingToGivenManager: incorrect position - " + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAllManagerReportingToGivenManager: query failed - " + err.Error(),
		}
	}

	var (
		manager      string
		managerArray = make([]string, 0)
	)

	for rows.Next() {

		if err := rows.Scan(&manager); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllManagerReportingToGivenManager: scan failed - " + err.Error(),
			}
		}

		if manager == managerID {
			continue
		}

		if manager != "" {
			managerArray = append(managerArray, manager)
		}
		manager = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachManager := range managerArray {
		userRecord, err := dbConn.FetchUserDetails(eachManager, quarter, year)
		if err != nil {
			return nil, err
		}

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllEmpReportingToGivenManagerRole -
func (dbConn *DBImplementation) FetchAllMngrReportingToManagerRole(managerID, managerPosition string, quarter, year int) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllMngrReportingToManagerRole start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllMngrReportingToManagerRole done")

	_, ok := constant.ManagerMap[managerPosition]

	if managerID == "" || !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + " DB CALL WITHOOUT DETAILS",
			Details:    "FetchAllMngrReportingToManagerRole: invalid db call, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.general_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.vice_president=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.GeneralManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.senior_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.general_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.SeniorManager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.senior_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.Manager:
		query = "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + " MANAGER NOT SUPPORTED",
			Details:    "FetchAllMngrReportingToManagerRole: not supported, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAllMngrReportingToManagerRole: query failed: " + err.Error(),
		}
	}

	var (
		manager      string
		managerArray = make([]string, 0)
	)

	for rows.Next() {

		if err := rows.Scan(&manager); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllMngrReportingToManagerRole: scan failed - " + err.Error(),
			}
		}

		if manager == managerID {
			continue
		}

		if manager != "" {
			managerArray = append(managerArray, manager)
		}
		manager = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachManager := range managerArray {
		userRecord, err := dbConn.FetchUserDetails(eachManager, quarter, year)
		if err != nil {
			return nil, err
		}

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllEmpReportingToGivenManagerRole
func (dbConn *DBImplementation) FetchAllEmpReportingToGivenManagerRole(managerPosition string, managerID string, quarter, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllEmpReportingToGivenManagerRole start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllEmpReportingToGivenManagerRole done")

	_, ok := constant.ManagerMap[managerPosition]
	if managerID == "" || !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + ": INVALID MANAGER POSITION - " + managerPosition,
			Details:    "FetchAllEmpReportingToGivenManagerRole: invalid position, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct a.emp_id, a.name, a.usr_position, a.region, a.outlet  from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.general_manager AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.vice_president=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.GeneralManager:
		query = "select distinct a.emp_id, a.name, a.usr_position, a.region, a.outlet  from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.senior_manager AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.general_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.SeniorManager:
		query = "select distinct a.emp_id, a.name, a.usr_position, a.region, a.outlet  from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.manager AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.senior_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.Manager:
		query = "select distinct a.emp_id, a.name, a.usr_position, a.region, a.outlet  from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.outlet_manager AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	case constant.OutletManager:
		query = "select distinct a.emp_id, a.name, a.usr_position, a.region, a.outlet from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.emp_id AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.outlet_manager=$1 AND a.quarter=$2 AND a.year=$3" // AND a.isactive=true
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAllEmpReportingToGivenManagerRole: query failed - " + err.Error(),
		}
	}

	var (
		array = make([]*models.ReportingStaffIDAndName, 0)
	)

	for rows.Next() {

		usr := &models.ReportingStaffIDAndName{}

		if err := rows.Scan(&usr.ID, &usr.Name, &usr.Position, &usr.Region, &usr.Outlet); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllEmpReportingToGivenManagerRole: scan failed - " + err.Error(),
			}
		}

		if usr.ID == managerID {
			continue
		}

		array = append(array, usr)
	}

	// if managerPosition == constant.Manager {

	// 	rows, err := dbConn.Operation.QueryContext(ctx,
	// 		"select distinct a.team_lead, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager=$1 and a.quarter=$2 and a.year=$3 and b.quarter=$4 and b.year=$5",
	// 		managerID, quarter, year, quarter, year)
	// 	if err != nil {
	// 		return nil, &errors.Error{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Code:       errors.InternalServerErrorCode,
	// 			Message:    errors.InternalServerError,
	// 			Details:    "FetchAllEmpReportingToGivenManagerRole: tl query failed - " + err.Error(),
	// 		}
	// 	}

	// 	for rows.Next() {

	// 		usr := &models.ReportingStaffIDAndName{}

	// 		if err := rows.Scan(&usr.ID, &usr.Position, &usr.Name); err != nil {
	// 			return nil, &errors.Error{
	// 				StatusCode: http.StatusInternalServerError,
	// 				Code:       errors.InternalServerErrorCode,
	// 				Message:    errors.InternalServerError,
	// 				Details:    "FetchAllEmpReportingToGivenManagerRole: tl scan failed - " + err.Error(),
	// 			}
	// 		}

	// 		array = append(array, usr)
	// 	}
	// }

	return array, nil
}

// FetchVP -
func (dbConn *DBImplementation) FetchVP(quarter, year int) ([]string, *errors.Error) {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchVP start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchVP done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.vice_president AND a.quarter=b.quarter AND a.year=b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND a.quarter=$1 AND a.year=$2 AND a.isactive=true"
	rows, err := dbConn.Operation.QueryContext(ctx, query, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchVP: query failed: " + err.Error(),
		}
	}

	var (
		vp      string
		vpArray = make([]string, 0)
	)

	for rows.Next() {
		vp = ""
		if err := rows.Scan(&vp); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchVP: scan failed - " + err.Error(),
			}
		}

		if vp != "" {
			vpArray = append(vpArray, vp)
		}
	}

	return vpArray, nil
}
