package impl

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

// FetchUserBasicDetails - userID, name, position, region, outlet, mobile, isAdmin, isActive
func (dbConn *DbServiceImpl) FetchUserBasicDetails(userID string, quarter, year int,
	args ...bool) (*models.UserBasicDetailsDB, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err  error
		user = &models.UserBasicDetailsDB{}
	)

	if args != nil {
		err = dbConn.DBClient.QueryRowContext(ctx, "select emp_id, name, usr_position, region, outlet, mobile_number, "+
			"isAdmin, isActive from user_master_table where deleted IS NULL and emp_id=$1 and isActive=true",
			userID).Scan(&user.UserId, &user.Name, &user.Position, &user.Region, &user.Outlet, &user.MobileNumber,
			&user.IsAdmin, &user.IsActive)
	} else {
		err = dbConn.DBClient.QueryRowContext(ctx, "select emp_id, name, usr_position, region, outlet, "+"mobile_number, "+
			"isAdmin, isActive from user_master_table where deleted IS NULL and emp_id=$1 and quarter=$2 and year=$3",
			userID, quarter, year).Scan(&user.UserId, &user.Name, &user.Position, &user.Region, &user.Outlet,
			&user.MobileNumber, &user.IsAdmin, &user.IsActive)
	}

	if err == nil {
		position, ok := constant.EmployeePositionMap[user.Position]
		if !ok {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
				Details: "FetchUserBasicDetails: incorrect user position: userId-" + user.UserId +
					" position-" + user.Position,
			}
		}
		user.Position = position
		return user, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.UserNotFound,
			Details:    "FetchUserBasicDetails: no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchUserBasicDetails: err - " + err.Error(),
	}
}

// TODO: FetchUserManagerDetails - for manger of that user
func (dbConn *DbServiceImpl) FetchUserManagerDetails(userID string, userPosition string, quarter,
	year int) (*models.UserManagerDetailsDB, *errors.Error) {

	var (
		query string
		user  = &models.UserManagerDetailsDB{
			UserId:   userID,
			Position: userPosition,
		}
	)

	switch {
	case userPosition == constant.VicePresident || userPosition == constant.VicePresidentAbb ||
		userPosition == constant.VicePresidentCC:
		return user, nil
	case userPosition == constant.GeneralManager || userPosition == constant.GeneralManagerAbb ||
		userPosition == constant.GeneralManagerCC:
		user.ManagerPosition = constant.VicePresidentAbb
		query = "select vice_president from user_hierarchy_table where deleted IS NULL and general_manager=$1 " +
			"and quarter=$2 and year=$3 limit 1"

	case userPosition == constant.SeniorManager || userPosition == constant.SeniorManagerAbb ||
		userPosition == constant.SeniorManagerCC:
		user.ManagerPosition = constant.GeneralManagerAbb
		query = "select general_manager from user_hierarchy_table where deleted IS NULL and senior_manager=$1 " +
			"and quarter=$2 and year=$3 limit 1"

	case userPosition == constant.Manager || userPosition == constant.ManagerAbb ||
		userPosition == constant.ManagerCC:
		user.ManagerPosition = constant.SeniorManagerAbb
		query = "select senior_manager from user_hierarchy_table where deleted IS NULL and manager=$1 and " +
			"quarter=$2 and year=$3 limit 1"

	case userPosition == constant.OutletManager || userPosition == constant.OutletManagerAbb ||
		userPosition == constant.OutletManagerCC:
		user.ManagerPosition = constant.ManagerAbb
		query = "select manager from user_hierarchy_table where deleted IS NULL and outlet_manager=$1 and " +
			"quarter=$2 and year=$3 limit 1"

	case userPosition == constant.Promoter || userPosition == constant.PromoterAbb ||
		userPosition == constant.PromoterCC || userPosition == constant.CSR ||
		userPosition == constant.CSRAbb || userPosition == constant.CSRCC:
		user.ManagerPosition = constant.OutletManagerAbb
		query = "select outlet_manager from user_hierarchy_table where deleted IS NULL and emp_id=$1 and " +
			"quarter=$2 and year=$3 limit 1"

	default:
		return user, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "INVALID USER POSITION",
			Details:    "FetchUserManagerDetails: id-" + userID + " incorrect position - " + userPosition,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	scanErr := dbConn.DBClient.QueryRowContext(ctx, query, user.UserId, quarter, year).Scan(&user.ManagerId)
	if scanErr != nil {
		if strings.Contains(scanErr.Error(), "no rows") {
			return user, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    errors.UserNotFound,
				Details:    "FetchUserManagerDetails: id-" + userID + " no rows mngr - " + scanErr.Error(),
			}
		}

		return user, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchUserManagerDetails: id-" + userID + " err mngr - " + scanErr.Error(),
		}
	}

	return user, nil
}

// FetchTotalNonManagerialEmployees - all non-managerial employees count of a given region and outlet
func (dbConn *DbServiceImpl) FetchTotalNonManagerialEmployees(region, outlet, channel string, quarter,
	year int) (int, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var count int

	channelFilter := ""
	switch channel {
	case "all":
		channelFilter = " AND a.chanel in ('DIRECT', 'SELCO', 'INFOLINE')"
	case "all-franchise":
		channelFilter = " AND a.chanel in ('INFOLINE', 'SELCO')"
	case "selco":
		channelFilter = " AND a.chanel = 'SELCO'"
	case "franchiseinfoline":
		channelFilter = " AND a.chanel = 'INFOLINE'"
	case "direct":
		channelFilter = " AND a.chanel = 'DIRECT'"
	}

	//usr_position=='BACK OFFICE'
	err := dbConn.DBClient.QueryRowContext(ctx, "select count(*) from user_master_table where "+
		channelFilter+
		"deleted IS NULL and isactive=true and region=$1 and outlet=$2 and "+
		"(usr_position='PROMOTER' OR usr_position='CSR') and quarter=$3 and year=$4",
		region, outlet, quarter, year).Scan(&count)
	if err != nil {
		return -1, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchTotalNonManagerialEmployees: query - " + err.Error(),
		}
	}

	return count, nil
}

// FetchActiveUserSales - all active sale transaction of a given user and time period
func (dbConn *DbServiceImpl) FetchActiveUserSales(userID string, channel string, isQuarterReport bool,
	quarter, month, year int) (*models.SaleTransaction, *errors.Error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resp := &models.SaleTransaction{
		EmpID:       userID,
		SaleRecords: make(map[string][]*models.SaleRecord),
	}

	channelFilter := ""
	switch channel {
	case "all":
		channelFilter = " AND a.chanel in ('DIRECT', 'SELCO', 'INFOLINE')"
	case "all-franchise":
		channelFilter = " AND a.chanel in ('INFOLINE', 'SELCO')"
	case "selco":
		channelFilter = " AND a.chanel = 'SELCO'"
	case "franchiseinfoline":
		channelFilter = " AND a.chanel = 'INFOLINE'"
	case "direct":
		channelFilter = " AND a.chanel = 'DIRECT'"
	}

	var (
		rows *sql.Rows
		err  error
	)

	if isQuarterReport {
		rows, err = dbConn.DBClient.QueryContext(ctx, "select id, prev_plan, curr_plan, quantity, product_group, "+
			"transaction_type, subcategory, mnp_donor, sales_date, terminate_date from sale_transaction_table where emp_id=$1 and "+
			"quarter=$2 and EXTRACT('year' from sales_date)=$3"+channelFilter,
			userID, quarter, year)
	} else {
		rows, err = dbConn.DBClient.QueryContext(ctx, "select id, prev_plan, curr_plan, quantity, product_group, "+
			"transaction_type, subcategory, mnp_donor, sales_date, terminate_date from sale_transaction_table where emp_id=$1 and "+
			"EXTRACT('month' from sales_date)=$2 and EXTRACT('year' from sales_date)=$3"+channelFilter,
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
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var (
			productGroup  string
			subCategory   interface{}
			salesDate     time.Time
			terminateDate time.Time
			rcrd          models.SaleRecord
		)

		if err := rows.Scan(&rcrd.TrxId, &rcrd.PreviousProduct, &rcrd.CurrentProduct, &rcrd.Quantity, &productGroup,
			&rcrd.Category, &subCategory, &rcrd.MNPDonor, &salesDate, &terminateDate); err != nil {
			return nil, &errors.Error{
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
			logger.FileLogger.CtxWarn(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", rcrd.TrxId).
				Msg("FetchActiveUserSales: incorrect productGroup")
		}

		if !(productGroup == constant.HBBCategory || productGroup == constant.HayyakCategory ||
			productGroup == constant.BaqatiCategory) {
			continue
		}

		if !terminateDate.IsZero() && terminateDate.Before(salesDate) {
			logger.FileLogger.CtxError(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", rcrd.TrxId).
				Msg("Db-FetchActiveUserSales: saleDate-terminateDate incorrect")
			continue
		}

		// expiryDate := salesDate.Add(90 * 24 * time.Hour)
		// if !terminateDate.IsZero() && terminateDate.Sub(expiryDate) <= 0 { //Months<3
		// 	if !strings.Contains(subCategory.(string), "INVOLUNTARY DEACTIVATION") {
		// 		logger.FileLogger.CtxError(ctx).Str("userID", userID).Str("region", region).Str("outlet", outlet).
		// 			Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
		// 			Int("trx", id).Msg("FetchActiveUserSales: sale terminated")
		// 		continue
		// 	}
		// }

		//TODO: new logic is ignore records if it is terminated on the same month or same quarter depending on
		var isTerminated bool
		if !terminateDate.IsZero() {
			terminateMonth := int(terminateDate.Month())
			terminateQuarter := objecthelper.GetQuarterFromMonth(int(terminateDate.Month()))
			if isQuarterReport && terminateQuarter == quarter {
				isTerminated = true
			}
			if !isQuarterReport && terminateMonth == month {
				isTerminated = true
			}
		}

		if isTerminated {
			logger.FileLogger.CtxWarn(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Int("trx", rcrd.TrxId).
				Interface("terminateDate", terminateDate).Msg("Db-FetchActiveUserSales: sale terminated")
			continue
		}

		//ignore HBB Upgrade records
		if productGroup == constant.HBBCategory && rcrd.Category == constant.UpgradeCategory {
			continue
		}

		if resp.SaleRecords[productGroup] == nil {
			resp.SaleRecords[productGroup] = make([]*models.SaleRecord, 0)
		}

		resp.SaleRecords[productGroup] = append(resp.SaleRecords[productGroup], &rcrd)
	}

	return resp, nil
}

// FetchAdjustedActiveUserSales - for quarter only fetch all active sale transaction of current quarter
// ignoring those expired in that quarter as well as in this next quarter
func (dbConn *DbServiceImpl) FetchAdjustedActiveUserSales(userID string, channel string, isQuarterReport bool,
	quarter, month, year int) (*models.SaleTransaction, *errors.Error) {
	if !isQuarterReport {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAdjustedActiveUserSales: userID-" + userID + ", valid for quarter only",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resp := &models.SaleTransaction{
		EmpID:       userID,
		SaleRecords: make(map[string][]*models.SaleRecord),
	}

	channelFilter := ""
	switch channel {
	case "all":
		channelFilter = " AND a.chanel in ('DIRECT', 'SELCO', 'INFOLINE')"
	case "all-franchise":
		channelFilter = " AND a.chanel in ('INFOLINE', 'SELCO')"
	case "selco":
		channelFilter = " AND a.chanel = 'SELCO'"
	case "franchiseinfoline":
		channelFilter = " AND a.chanel = 'INFOLINE'"
	case "direct":
		channelFilter = " AND a.chanel = 'DIRECT'"
	}

	var (
		rows *sql.Rows
		err  error
	)

	rows, err = dbConn.DBClient.QueryContext(ctx,
		"select id, prev_plan, curr_plan, quantity, product_group, transaction_type, subcategory, mnp_donor, "+
			"sales_date, terminate_date from sale_transaction_table where "+
			"emp_id=$1 and quarter=$2 and EXTRACT('year' from sales_date)=$3"+channelFilter,
		userID, quarter, year)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAdjustedActiveUserSales: userID-" + userID + ", query failed - " + err.Error(),
		}
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var (
			productGroup  string
			subCategory   interface{}
			salesDate     time.Time
			terminateDate time.Time
			rcrd          models.SaleRecord
		)

		if err := rows.Scan(&rcrd.TrxId, &rcrd.PreviousProduct, &rcrd.CurrentProduct, &rcrd.Quantity, &productGroup,
			&rcrd.Category, &subCategory, &rcrd.MNPDonor, &salesDate, &terminateDate); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchAdjustedActiveUserSales: " + userID + ", scan failed - " + err.Error(),
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
			logger.FileLogger.CtxWarn(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("year", year).Int("trx", rcrd.TrxId).
				Msg("Db-FetchAdjustedActiveUserSales: incorrect productGroup")
		}

		if !(productGroup == constant.HBBCategory || productGroup == constant.HayyakCategory ||
			productGroup == constant.BaqatiCategory) {
			continue
		}

		if !terminateDate.IsZero() && terminateDate.Before(salesDate) {
			logger.FileLogger.CtxError(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("year", year).Int("trx", rcrd.TrxId).
				Msg("Db-FetchAdjustedActiveUserSales: saleDate-terminateDate incorrect")
			continue
		}

		nextQuarter, nextYear := objecthelper.GetNextQuarterAndMonth(quarter, year)

		// this logic ignores terminated orders in current quarter and next quarter
		var isTerminated bool
		if !terminateDate.IsZero() {
			terminateQuarter := objecthelper.GetQuarterFromMonth(int(terminateDate.Month()))
			if isQuarterReport &&
				((terminateQuarter == quarter && terminateDate.Year() == year) ||
					(terminateQuarter == nextQuarter && terminateDate.Year() == nextYear)) {
				isTerminated = true
			}
		}

		if isTerminated {
			logger.FileLogger.CtxWarn(ctx).Str("userID", userID).Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("year", year).Int("trx", rcrd.TrxId).
				Interface("terminateDate", terminateDate).Msg("Db-FetchActiveUserSales: sale terminated")
			continue
		}

		//ignore HBB Upgrade records
		if productGroup == constant.HBBCategory && rcrd.Category == constant.UpgradeCategory {
			continue
		}

		if resp.SaleRecords[productGroup] == nil {
			resp.SaleRecords[productGroup] = make([]*models.SaleRecord, 0)
		}

		resp.SaleRecords[productGroup] = append(resp.SaleRecords[productGroup], &rcrd)
	}

	return resp, nil
}
