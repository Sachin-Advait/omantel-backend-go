package psqldb

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// FetchOutletReport - get pre-computed outlet report
func (dbConn *DBImplementation) FetchOutletReport(outlet string, isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	// logger.FileLogger.Info().Str("outlet", outlet).Time("start", time.Now()).Msg("FetchOutletReport start")
	// defer logger.FileLogger.Info().Str("outlet", outlet).Time("end", time.Now()).Msg("FetchOutletReport done")

	if outlet == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "blank outlet",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	var (
		hbbTarget     float64
		hayyakTarget  float64
		baqatiTarget  float64
		upgradeTarget float64
	)

	rcrd := &models.OutletReportDB{}

	err := dbConn.Operation.QueryRowContext(ctx, "select outlet, region, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from outlet_reports_table where outlet=$1 and month=$2 and quarter=$3 and year=$4",
		outlet, month, quarter, year).Scan(
		&rcrd.Outlet,
		&rcrd.Region,
		&rcrd.Latitude,
		&rcrd.Longitude,
		&rcrd.HBBSalesCount,
		&hbbTarget,
		&rcrd.HayyakSalesCount,
		&hayyakTarget,
		&rcrd.BaqatiSalesCount,
		&baqatiTarget,
		&rcrd.UpgradeSalesCount,
		&upgradeTarget,
		&rcrd.CreatedAt,
	)

	rcrd.HBBTargetCount = int(hbbTarget)
	rcrd.HayyakTargetCount = int(hayyakTarget)
	rcrd.BaqatiTargetCount = int(baqatiTarget)
	rcrd.UpgradeTargetCount = int(upgradeTarget)

	if err == nil {
		return rcrd, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORD FOUND",
			Details:    "FetchOutletReport: " + outlet + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchOutletReport: err - " + err.Error(),
	}
}

// InsertOutletReport - upsert outlet report
func (dbConn *DBImplementation) InsertOutletReport(outlet *models.OutletReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	// logger.FileLogger.Info().Str("outlet", outlet.Outlet).Time("start", time.Now()).Msg("InsertOutletReport start")
	// defer logger.FileLogger.Info().Str("outlet", outlet.Outlet).Time("end", time.Now()).Msg("InsertOutletReport done")

	if outlet == nil || outlet.Outlet == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid outlet",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rslt, err := dbConn.Operation.ExecContext(ctx, "INSERT INTO outlet_reports_table(outlet, region, latitude, longitude, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created, modified, isquarterly, month, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) "+
		"ON CONFLICT(outlet, month, quarter, year) DO UPDATE SET hbb_sales_count=$19, hbb_target=$20, hayyak_sales_count=$21, hayyak_target=$22, baqati_sales_count=$23, baqati_target=$24, upgrade_sales_count=$25, upgrade_target=$26, modified=$27",
		outlet.Outlet,
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
		outlet.CreatedAt,
		outlet.ModifiedAt,
		isQuarterReport,
		month,
		quarter,
		year,
		outlet.HBBSalesCount,
		outlet.HBBTargetCount,
		outlet.HayyakSalesCount,
		outlet.HayyakTargetCount,
		outlet.BaqatiSalesCount,
		outlet.BaqatiTargetCount,
		outlet.UpgradeSalesCount,
		outlet.UpgradeTargetCount,
		outlet.ModifiedAt,
	)

	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertOutletReport: insert failed - " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertOutletReport: rows affected - " + strconv.Itoa(int(rows))}
	}

	return nil
}

// FetchRegionReport - get pre-computed region report
func (dbConn *DBImplementation) FetchRegionReport(region string, isQuarterReport bool, quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	// logger.FileLogger.Info().Str("region", region).Time("start", time.Now()).Msg("FetchRegionReport start")
	// defer logger.FileLogger.Info().Str("region", region).Time("end", time.Now()).Msg("FetchRegionReport done")

	if region == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid region",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	var (
		err           error
		hbbTarget     float64
		hayyakTarget  float64
		baqatiTarget  float64
		upgradeTarget float64
	)

	rcrd := &models.RegionReportDB{}

	err = dbConn.Operation.QueryRowContext(ctx, "select region, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created from outlet_reports_table where region=$1 and month=$2 and quarter=$3 and year=$4",
		region, month, quarter, year).Scan(
		&rcrd.Region,
		&rcrd.HBBSalesCount,
		&hbbTarget,
		&rcrd.HayyakSalesCount,
		&hayyakTarget,
		&rcrd.BaqatiSalesCount,
		&baqatiTarget,
		&rcrd.UpgradeSalesCount,
		&upgradeTarget,
		&rcrd.CreatedAt,
	)

	rcrd.HBBTargetCount = int(hbbTarget)
	rcrd.HayyakTargetCount = int(hayyakTarget)
	rcrd.BaqatiTargetCount = int(baqatiTarget)
	rcrd.UpgradeTargetCount = int(upgradeTarget)

	if err == nil {
		return rcrd, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORD FOUND",
			Details:    "FetchRegionReport: " + region + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchRegionReport: err - " + err.Error(),
	}
}

// InsertRegionReport - upsert region report
func (dbConn *DBImplementation) InsertRegionReport(region *models.RegionReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	// logger.FileLogger.Info().Str("region", region.Region).Time("start", time.Now()).Msg("InsertRegionReport start")
	// defer logger.FileLogger.Info().Str("region", region.Region).Time("end", time.Now()).Msg("InsertRegionReport done")

	if region == nil || region.Region == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid region",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rslt, err := dbConn.Operation.ExecContext(ctx, "INSERT INTO region_reports_table(region, hbb_sales_count, hbb_target, hayyak_sales_count, hayyak_target, baqati_sales_count, baqati_target, upgrade_sales_count, upgrade_target, created, modified, isquarterly, month, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) "+
		"ON CONFLICT(region, month, quarter, year) DO UPDATE SET hbb_sales_count=$16, hbb_target=$17, hayyak_sales_count=$18, hayyak_target=$19, baqati_sales_count=$20, baqati_target=$21, upgrade_sales_count=$22, upgrade_target=$23, modified=$24",
		region.Region,
		region.HBBSalesCount,
		region.HBBTargetCount,
		region.HayyakSalesCount,
		region.HayyakTargetCount,
		region.BaqatiSalesCount,
		region.BaqatiTargetCount,
		region.UpgradeSalesCount,
		region.UpgradeTargetCount,
		region.CreatedAt,
		region.ModifiedAt,
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
		region.ModifiedAt,
	)

	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertRegionReport: insert failed - " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertRegionReport: rows affected - " + strconv.Itoa(int(rows))}
	}

	return nil
}

// FetchAgentUserReport - get pre-computed agent user report
func (dbConn *DBImplementation) FetchAgentUserReport(userId string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if userId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid userId",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rcrd := &models.UserReportDB{}

	err := dbConn.Operation.QueryRowContext(ctx, "select emp_id, position, reporting_manager, reporting_manager_position, region, outlet, latitude, longitude, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created from user_reports_table where emp_id=$1 AND month=$2 AND quarter=$3 AND year=$4",
		userId, month, quarter, year).Scan(
		&rcrd.UserId,
		&rcrd.Position,
		&rcrd.ManagerId,
		&rcrd.ManagerPosition,
		&rcrd.Region,
		&rcrd.Outlet,
		&rcrd.Latitude,
		&rcrd.Longitude,
		&rcrd.HayyakSalesCount,
		&rcrd.HayyakTargetCount,
		&rcrd.HayyakAchievedPercentage,
		&rcrd.HayyakPayout,
		&rcrd.HayyakExtraPayout,
		&rcrd.HBBSalesCount,
		&rcrd.HBBTargetCount,
		&rcrd.HBBAchievedPercentage,
		&rcrd.HBBPayout,
		&rcrd.HBBExtraPayout,
		&rcrd.BaqatiSalesCount,
		&rcrd.BaqatiTargetCount,
		&rcrd.BaqatiAchievedPercentage,
		&rcrd.BaqatiPayout,
		&rcrd.BaqatiExtraPayout,
		&rcrd.UpgradeSalesCount,
		&rcrd.UpgradeTargetCount,
		&rcrd.UpradeAchievedPercentage,
		&rcrd.UpgradePayout,
		&rcrd.UpgradeExtraPayout,
		&rcrd.TotalPayout,
		&rcrd.HayyakSaleCountPerPlanArray,
		&rcrd.HBBSaleCountPerPlanArray,
		&rcrd.BaqatiSaleCountPerPlanArray,
		&rcrd.UpgradeSaleCountPerPlanArray,
		&rcrd.CreatedAt,
	)

	if err == nil {
		return rcrd, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORD FOUND",
			Details:    "FetchAgentUserReport: " + userId + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchAgentUserReport: " + userId + " err - " + err.Error(),
	}
}

// InsertAgentUserReport - upsert agent user report
func (dbConn *DBImplementation) InsertAgentUserReport(usr *models.UserReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	// logger.FileLogger.Info().Str("clientID", usr.UserId).Time("start", time.Now()).Msg("InsertAgentUserReport start")
	// defer logger.FileLogger.Info().Str("clientID", usr.UserId).Time("end", time.Now()).Msg("InsertAgentUserReport done")

	if usr == nil || usr.UserId == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid user agent",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rslt, err := dbConn.Operation.ExecContext(ctx, "INSERT INTO user_reports_table(select emp_id, position, reporting_manager, reporting_manager_position, region, outlet, latitude, longitude, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created, modified, isquarterly, month, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39) "+
		"ON CONFLICT(emp_id, month, quarter, year) DO UPDATE SET position=$40, reporting_manager=$41, reporting_manager_position=$42, region=$43, outlet=$44, latitude=$45, longitude=$46, hayyak_sales_count=$47, hayyak_target=$48, hayyak_achieved_percentage=$49, hayyak_payout=$50, hayyak_extra_payout=$51, hbb_sales_count=$52, hbb_target=$53, hbb_achieved_percentage=$54, hbb_payout=$55, hbb_extra_payout=$56, baqati_sales_count=$57, baqati_target=$58, baqati_achieved_percentage=$59, baqati_payout=$60, baqati_extra_payout=$61, upgrade_sales_count=$62, upgrade_target=$63, uprade_achieved_percentage=$64, upgrade_payout=$65, upgrade_extra_payout=$66, total_payout=$67, hayyak_sales_count_per_plan=$68, hbb_sales_count_per_plan=$69, baqati_sales_count_per_plan=$70, upgrade_sales_count_per_plan=$71, modified=$72",
		&usr.UserId,
		&usr.Position,
		&usr.ManagerId,
		&usr.ManagerPosition,
		&usr.Region,
		&usr.Outlet,
		&usr.Latitude,
		&usr.Longitude,
		&usr.HayyakSalesCount,
		&usr.HayyakTargetCount,
		&usr.HayyakAchievedPercentage,
		&usr.HayyakPayout,
		&usr.HayyakExtraPayout,
		&usr.HBBSalesCount,
		&usr.HBBTargetCount,
		&usr.HBBAchievedPercentage,
		&usr.HBBPayout,
		&usr.HBBExtraPayout,
		&usr.BaqatiSalesCount,
		&usr.BaqatiTargetCount,
		&usr.BaqatiAchievedPercentage,
		&usr.BaqatiPayout,
		&usr.BaqatiExtraPayout,
		&usr.UpgradeSalesCount,
		&usr.UpgradeTargetCount,
		&usr.UpradeAchievedPercentage,
		&usr.UpgradePayout,
		&usr.UpgradeExtraPayout,
		&usr.TotalPayout,
		&usr.HayyakSaleCountPerPlanArray,
		&usr.HBBSaleCountPerPlanArray,
		&usr.BaqatiSaleCountPerPlanArray,
		&usr.UpgradeSaleCountPerPlanArray,
		&usr.CreatedAt,
		&usr.ModifiedAt,
		isQuarterReport,
		month,
		quarter,
		year,
		&usr.Position,
		&usr.ManagerId,
		&usr.ManagerPosition,
		&usr.Region,
		&usr.Outlet,
		&usr.Latitude,
		&usr.Longitude,
		&usr.HayyakSalesCount,
		&usr.HayyakTargetCount,
		&usr.HayyakAchievedPercentage,
		&usr.HayyakPayout,
		&usr.HayyakExtraPayout,
		&usr.HBBSalesCount,
		&usr.HBBTargetCount,
		&usr.HBBAchievedPercentage,
		&usr.HBBPayout,
		&usr.HBBExtraPayout,
		&usr.BaqatiSalesCount,
		&usr.BaqatiTargetCount,
		&usr.BaqatiAchievedPercentage,
		&usr.BaqatiPayout,
		&usr.BaqatiExtraPayout,
		&usr.UpgradeSalesCount,
		&usr.UpgradeTargetCount,
		&usr.UpradeAchievedPercentage,
		&usr.UpgradePayout,
		&usr.UpgradeExtraPayout,
		&usr.TotalPayout,
		&usr.HayyakSaleCountPerPlanArray,
		&usr.HBBSaleCountPerPlanArray,
		&usr.BaqatiSaleCountPerPlanArray,
		&usr.UpgradeSaleCountPerPlanArray,
		&usr.ModifiedAt,
	)

	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertAgentUserReport: insert failed - " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertAgentUserReport: rows affected - " + strconv.Itoa(int(rows))}
	}

	return nil
}

// FetchManagerReport - get pre-computed nmanager report
func (dbConn *DBImplementation) FetchManagerReport(userId string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	if userId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid userId",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rcrd := &models.UserReportDB{}

	err := dbConn.Operation.QueryRowContext(ctx, "select emp_id, position, reporting_manager, reporting_manager_position, region, outlet, latitude, longitude, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created from manager_reports_table where emp_id=$1 AND month=$2 AND quarter=$3 AND year=$4",
		userId, month, quarter, year).Scan(
		&rcrd.UserId,
		&rcrd.Position,
		&rcrd.ManagerId,
		&rcrd.ManagerPosition,
		&rcrd.Region,
		&rcrd.Outlet,
		&rcrd.Latitude,
		&rcrd.Longitude,
		&rcrd.HayyakSalesCount,
		&rcrd.HayyakTargetCount,
		&rcrd.HayyakAchievedPercentage,
		&rcrd.HayyakPayout,
		&rcrd.HayyakExtraPayout,
		&rcrd.HBBSalesCount,
		&rcrd.HBBTargetCount,
		&rcrd.HBBAchievedPercentage,
		&rcrd.HBBPayout,
		&rcrd.HBBExtraPayout,
		&rcrd.BaqatiSalesCount,
		&rcrd.BaqatiTargetCount,
		&rcrd.BaqatiAchievedPercentage,
		&rcrd.BaqatiPayout,
		&rcrd.BaqatiExtraPayout,
		&rcrd.UpgradeSalesCount,
		&rcrd.UpgradeTargetCount,
		&rcrd.UpradeAchievedPercentage,
		&rcrd.UpgradePayout,
		&rcrd.UpgradeExtraPayout,
		&rcrd.TotalPayout,
		&rcrd.HayyakSaleCountPerPlanArray,
		&rcrd.HBBSaleCountPerPlanArray,
		&rcrd.BaqatiSaleCountPerPlanArray,
		&rcrd.UpgradeSaleCountPerPlanArray,
		&rcrd.CreatedAt,
	)

	if err == nil {
		return rcrd, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORD FOUND",
			Details:    "FetchManagerReport: " + userId + ", no rows - " + err.Error(),
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchManagerReport: " + userId + " err - " + err.Error(),
	}
}

// InsertManagerReport - upsert manager report
func (dbConn *DBImplementation) InsertManagerReport(usr *models.UserReportDB, isQuarterReport bool, quarter, month, year int) *errors.Error {

	// logger.FileLogger.Info().Str("clientID", user.UserId).Time("start", time.Now()).Msg("InsertManagerReport start")
	// defer logger.FileLogger.Info().Str("clientID", clientID.UserId).Time("end", time.Now()).Msg("InsertManagerReport done")

	if usr == nil || usr.UserId == "" {
		return &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    errors.InvalidDetails,
			Details:    "invalid manager",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if isQuarterReport {
		month = 0
	} else {
		quarter = 0
	}

	rslt, err := dbConn.Operation.ExecContext(ctx, "INSERT INTO manager_reports_table(select emp_id, position, reporting_manager, reporting_manager_position, region, outlet, latitude, longitude, hayyak_sales_count, hayyak_target, hayyak_achieved_percentage, hayyak_payout, hayyak_extra_payout, hbb_sales_count, hbb_target, hbb_achieved_percentage, hbb_payout, hbb_extra_payout, baqati_sales_count, baqati_target, baqati_achieved_percentage, baqati_payout, baqati_extra_payout, upgrade_sales_count, upgrade_target, uprade_achieved_percentage, upgrade_payout, upgrade_extra_payout, total_payout, hayyak_sales_count_per_plan, hbb_sales_count_per_plan, baqati_sales_count_per_plan, upgrade_sales_count_per_plan, created, modified, isquarterly, month, quarter, year) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39) "+
		"ON CONFLICT(emp_id, month, quarter, year) DO UPDATE SET position=$40, reporting_manager=$41, reporting_manager_position=$42, region=$43, outlet=$44, latitude=$45, longitude=$46, hayyak_sales_count=$47, hayyak_target=$48, hayyak_achieved_percentage=$49, hayyak_payout=$50, hayyak_extra_payout=$51, hbb_sales_count=$52, hbb_target=$53, hbb_achieved_percentage=$54, hbb_payout=$55, hbb_extra_payout=$56, baqati_sales_count=$57, baqati_target=$58, baqati_achieved_percentage=$59, baqati_payout=$60, baqati_extra_payout=$61, upgrade_sales_count=$62, upgrade_target=$63, uprade_achieved_percentage=$64, upgrade_payout=$65, upgrade_extra_payout=$66, total_payout=$67, hayyak_sales_count_per_plan=$68, hbb_sales_count_per_plan=$69, baqati_sales_count_per_plan=$70, upgrade_sales_count_per_plan=$71, modified=$72",
		&usr.UserId,
		&usr.Position,
		&usr.ManagerId,
		&usr.ManagerPosition,
		&usr.Region,
		&usr.Outlet,
		&usr.Latitude,
		&usr.Longitude,
		&usr.HayyakSalesCount,
		&usr.HayyakTargetCount,
		&usr.HayyakAchievedPercentage,
		&usr.HayyakPayout,
		&usr.HayyakExtraPayout,
		&usr.HBBSalesCount,
		&usr.HBBTargetCount,
		&usr.HBBAchievedPercentage,
		&usr.HBBPayout,
		&usr.HBBExtraPayout,
		&usr.BaqatiSalesCount,
		&usr.BaqatiTargetCount,
		&usr.BaqatiAchievedPercentage,
		&usr.BaqatiPayout,
		&usr.BaqatiExtraPayout,
		&usr.UpgradeSalesCount,
		&usr.UpgradeTargetCount,
		&usr.UpradeAchievedPercentage,
		&usr.UpgradePayout,
		&usr.UpgradeExtraPayout,
		&usr.TotalPayout,
		&usr.HayyakSaleCountPerPlanArray,
		&usr.HBBSaleCountPerPlanArray,
		&usr.BaqatiSaleCountPerPlanArray,
		&usr.UpgradeSaleCountPerPlanArray,
		&usr.CreatedAt,
		&usr.ModifiedAt,
		isQuarterReport,
		month,
		quarter,
		year,
		&usr.Position,
		&usr.ManagerId,
		&usr.ManagerPosition,
		&usr.Region,
		&usr.Outlet,
		&usr.Latitude,
		&usr.Longitude,
		&usr.HayyakSalesCount,
		&usr.HayyakTargetCount,
		&usr.HayyakAchievedPercentage,
		&usr.HayyakPayout,
		&usr.HayyakExtraPayout,
		&usr.HBBSalesCount,
		&usr.HBBTargetCount,
		&usr.HBBAchievedPercentage,
		&usr.HBBPayout,
		&usr.HBBExtraPayout,
		&usr.BaqatiSalesCount,
		&usr.BaqatiTargetCount,
		&usr.BaqatiAchievedPercentage,
		&usr.BaqatiPayout,
		&usr.BaqatiExtraPayout,
		&usr.UpgradeSalesCount,
		&usr.UpgradeTargetCount,
		&usr.UpradeAchievedPercentage,
		&usr.UpgradePayout,
		&usr.UpgradeExtraPayout,
		&usr.TotalPayout,
		&usr.HayyakSaleCountPerPlanArray,
		&usr.HBBSaleCountPerPlanArray,
		&usr.BaqatiSaleCountPerPlanArray,
		&usr.UpgradeSaleCountPerPlanArray,
		&usr.ModifiedAt,
	)

	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertManagerReport: insert failed - " + err.Error(),
		}
	}

	rows, _ := rslt.RowsAffected()

	if rows != 1 {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertManagerReport: rows affected - " + strconv.Itoa(int(rows))}
	}

	return nil
}
