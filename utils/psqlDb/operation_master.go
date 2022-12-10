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
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

/////////////////////// OUTLET TABLE ///////////////////////////

func (dbConn *DBImplementation) FetchOutletDetailsViaName(regionName, outletName string) (*models.GetOutletResponse, *errors.Error) {

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	outlet := &models.GetOutletResponse{}

	err := dbConn.Operation.QueryRowContext(ctx,
		"SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table where region_name=$1 and outlet_name=$2", regionName, outletName).
		Scan(&outlet.OutletCode, &outlet.OutletName, &outlet.RegionCode, &outlet.RegionName, &outlet.Lat, &outlet.Lang)
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

func (dbConn *DBImplementation) FetchAllOutletNames(regionName string) ([]string, *errors.Error) {

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	outlets := make([]string, 0)

	rows, err := dbConn.Operation.QueryContext(ctx, "SELECT outlet_name FROM outlets_table where region_name=$1", regionName)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rows   *sql.Rows
		sqlErr error
	)

	if isQuarterReport {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select category, product_group, SUM(quantity) AS total from sale_transaction_table where region=$1 and outlet=$2 and quarter=$3 and EXTRACT('year' from sales_date)=$4 GROUP BY category, product_group",
			region, outlet, quarter, year)
	} else {
		rows, sqlErr = dbConn.Operation.QueryContext(ctx,
			"select category, product_group, SUM(quantity) AS total from sale_transaction_table where region=$1 and outlet=$2 and EXTRACT('month' from b.sales_date)=$3 and EXTRACT('year' from b.sales_date)=$4 GROUP BY category, product_group",
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

	var (
		category      string
		productGroup  string
		quantityTotal int
	)

	for rows.Next() {

		if scanErr := rows.Scan(&category, &productGroup, &quantityTotal); scanErr != nil {
			err = &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchOutletSalesCount: scan failed",
				Details:    "FetchOutletSalesCount: region-" + region + ", outlet-" + outlet + ", scan failed - " + err.Error(),
			}
			return
		}

		switch {
		case category == constant.UpgradeCategory:
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rows *sql.Rows
		err  error
	)

	if isQuarterReport {
		rows, err = dbConn.Operation.QueryContext(ctx, "select prev_product, curr_product, quantity, product_group, category, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and region=$2 and outlet=$3 and quarter=$4 and EXTRACT('year' from sales_date)=$5",
			userID, region, outlet, quarter, year)
	} else {
		rows, err = dbConn.Operation.QueryContext(ctx, "select prev_product, curr_product, quantity, product_group, category, subcategory, sales_date, inactive_date from sale_transaction_table where emp_id=$1 and region=$2 and outlet=$3 and EXTRACT('month' from sales_date)=$4 and EXTRACT('year' from b.sales_date)=$5",
			userID, region, outlet, month, year)
	}

	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchUserSales: userID-" + userID + ", query failed - " + err.Error(),
		}
	}

	var (
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
		if err := rows.Scan(&prev_product, &curr_product, &quantity, &product_group, &category, &subcategory, &salesDate, &inactiveDate); err != nil {
			return trx, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchUserSales: " + userID + ", scan failed - " + err.Error(),
			}
		}

		trx.EmpID = userID
		trx.Outlet = outlet
		trx.Region = region

		switch {
		case strings.Contains(product_group, "FIXED") || product_group == constant.HBBCategory:
			product_group = constant.HBBCategory

		case strings.Contains(product_group, "PREPAID") || product_group == constant.HayyakCategory:
			product_group = constant.HayyakCategory

		case strings.Contains(product_group, "POSTPAID") || product_group == constant.BaqatiCategory:
			product_group = constant.BaqatiCategory

		default:
			return trx, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchUserSales: invalid productGrp - " + product_group,
			}
		}

		lapsedSecs := salesDate.Sub(inactiveDate).Seconds()

		if objecthelper.RoundTime(lapsedSecs/2600640) < 3 { //Months<3
			if !strings.Contains(subcategory.(string), "INVOLUNTARY DEACTIVATION") {
				//TODO: terminated logs
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

func (dbConn *DBImplementation) FetchAgentUserHierarchy(userID string) (*models.UserHierarchy, *errors.Error) {

	// logger.FileLogger.Info().Str("userID", userID).Time("start", time.Now()).Msg("FetchAgentUserHierarchy start")
	// defer logger.FileLogger.Info().Str("userID", userID).Time("end", time.Now()).Msg("FetchAgentUserHierarchy done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rcrd := &models.UserHierarchy{}

	err := dbConn.Operation.QueryRowContext(ctx, "select outlet_manager, manager, senior_manager, general_manager, vice_president from user_hierarchy_table where emp_id", userID).
		Scan(&rcrd.OutletManager, &rcrd.Manager, &rcrd.SeniorManager, &rcrd.GeneralManager, &rcrd.VicePresident)
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
func (dbConn *DBImplementation) FetchAllAgentsUserOfOM(outletManagerID string) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("userID", userID).Time("start", time.Now()).Msg("FetchAgentUserHierarchy start")
	// defer logger.FileLogger.Info().Str("userID", userID).Time("end", time.Now()).Msg("FetchAgentUserHierarchy done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select distinct emp_id from user_hierarchy_table where outlet_manager=$1 and team_lead=$2", outletManagerID, outletManagerID)
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

		if userId != "" {
			userIds = append(userIds, userId)
		}
		userId = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachUser := range userIds {
		userRecord, err := dbConn.FetchUserDetails(eachUser)
		if err != nil {
			return nil, err
		}

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllOMReportingToGivenManager - all outlet managers reporting to given manager in argument
func (dbConn *DBImplementation) FetchAllOMReportingToGivenManager(managerPosition string, managerID string) (map[string]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllOMReportingToGivenManager start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllOMReportingToGivenManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select outlet_manager, team_lead from user_hierarchy_table where vice_president= $1"
	case constant.GeneralManager:
		query = "select outlet_manager, team_lead from user_hierarchy_table where general_manager= $1"
	case constant.SeniorManager:
		query = "select outlet_manager, team_lead from user_hierarchy_table where senior_manager= $1"
	case constant.Manager:
		query = "select outlet_manager, team_lead from user_hierarchy_table where manager= $1"
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "NOT A MANAGER",
			Details:    "FetchAllOMReportingToGivenManager: incorrect position - " + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
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
		teamLead      string
		response      = make(map[string]*models.UserBasicDetailsDB, 0)
	)

	for rows.Next() {

		if err := rows.Scan(&outletManager, &teamLead); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllOMReportingToGivenManager: scan failed - " + err.Error(),
			}
		}

		if response[outletManager] == nil {
			response[outletManager] = &models.UserBasicDetailsDB{}
		}

		if response[teamLead] == nil {
			response[teamLead] = &models.UserBasicDetailsDB{}
		}
	}

	for eachUser, _ := range response {
		userRecord, err := dbConn.FetchUserDetails(eachUser)
		if err != nil {
			return nil, err
		}

		response[eachUser] = userRecord
	}

	return response, nil
}

// FetchAllRMReportingToGivenManager - all regional managers reporting to given manager in argument
func (dbConn *DBImplementation) FetchAllRMReportingToGivenManager(managerPosition string, managerID string) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllManagerReportingToGivenManager start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllManagerReportingToGivenManager done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct manager from user_hierarchy_table where vice_president= $1"
	case constant.GeneralManager:
		query = "select distinct manager, team_lead from user_hierarchy_table where general_manager= $1"
	case constant.SeniorManager:
		query = "select distinct manager, team_lead from user_hierarchy_table where senior_manager= $1"
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "NOT A MANAGER",
			Details:    "FetchAllManagerReportingToGivenManager: incorrect position - " + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
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

		if manager != "" {
			managerArray = append(managerArray, manager)
		}
		manager = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachManager := range managerArray {
		userRecord, err := dbConn.FetchUserDetails(eachManager)
		if err != nil {
			return nil, err
		}

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllEmpReportingToGivenManagerRole -
func (dbConn *DBImplementation) FetchAllMngrReportingToManagerRole(managerID, managerPosition string) ([]*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllMngrReportingToManagerRole start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllMngrReportingToManagerRole done")

	_, ok := constant.ManagerMap[managerPosition]

	if managerID == "" || !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + "DB CALL WITHOOUT DETAILS",
			Details:    "FetchAllMngrReportingToManagerRole: invalid db call, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct general_manager from user_hierarchy_table where vice_president=$1"
	case constant.GeneralManager:
		query = "select distinct senior_manager from user_hierarchy_table where general_manager=$1"
	case constant.SeniorManager:
		query = "select distinct manager from user_hierarchy_table where senior_manager=$1"
	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + "MANAGER NOT SUPPORTED",
			Details:    "FetchAllMngrReportingToManagerRole: not supported, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
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

		if manager != "" {
			managerArray = append(managerArray, manager)
		}
		manager = ""
	}

	response := make([]*models.UserBasicDetailsDB, 0)

	for _, eachManager := range managerArray {
		userRecord, err := dbConn.FetchUserDetails(eachManager)
		if err != nil {
			return nil, err
		}

		response = append(response, userRecord)
	}

	return response, nil
}

// FetchAllEmpReportingToGivenManagerRole
func (dbConn *DBImplementation) FetchAllEmpReportingToGivenManagerRole(managerPosition string, managerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	// logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("start", time.Now()).Msg("FetchAllEmpReportingToGivenManagerRole start")
	// defer logger.FileLogger.Info().Str("managerID", managerID).Str("managerPosition", managerPosition).Time("end", time.Now()).Msg("FetchAllEmpReportingToGivenManagerRole done")

	_, ok := constant.ManagerMap[managerPosition]
	if managerID == "" || !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    managerID + "DB CALL WITHOOUT DETAILS",
			Details:    "FetchAllEmpReportingToGivenManagerRole: invalid db call, mngrID-" + managerID + ", mngrPos-" + managerPosition,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string

	switch managerPosition {
	case constant.VicePresident:
		query = "select distinct a.general_manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.general_manager=b.emp_id where a.vice_president=$1"
	case constant.GeneralManager:
		query = "select distinct a.senior_manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.senior_manager=b.emp_id where a.general_manager=$1"
	case constant.SeniorManager:
		query = "select distinct a.manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.manager=b.emp_id where a.senior_manager=$1"
	case constant.Manager:
		query = "select distinct a.outlet_manager, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager=$1"
	case constant.OutletManager:
		query = "select distinct a.emp_id, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.emp_id=b.emp_id where a.outlet_manager=$1"
	case constant.TeamLeader:
		query = "select distinct a.emp_id, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.emp_id=b.emp_id where a.outlet_manager=$1"
	}

	rows, err := dbConn.Operation.QueryContext(ctx, query, managerID)
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

		if err := rows.Scan(&usr.ID, &usr.Position, &usr.Name); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllEmpReportingToGivenManagerRole: scan failed - " + err.Error(),
			}
		}

		array = append(array, usr)
	}

	if managerPosition == constant.Manager {

		rows, err := dbConn.Operation.QueryContext(ctx,
			"select distinct a.team_lead, b.usr_position, b.name from user_hierarchy_table a INNER JOIN user_master_table b on a.outlet_manager=b.emp_id where a.manager=$1",
			managerID)
		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAllEmpReportingToGivenManagerRole: query failed - " + err.Error(),
			}
		}

		for rows.Next() {

			usr := &models.ReportingStaffIDAndName{}

			if err := rows.Scan(&usr.ID, &usr.Position, &usr.Name); err != nil {
				return nil, &errors.Error{
					StatusCode: http.StatusInternalServerError,
					Code:       errors.InternalServerErrorCode,
					Message:    errors.InternalServerError,
					Details:    "FetchAllEmpReportingToGivenManagerRole: scan failed - " + err.Error(),
				}
			}

			array = append(array, usr)
		}
	}

	return array, nil
}
