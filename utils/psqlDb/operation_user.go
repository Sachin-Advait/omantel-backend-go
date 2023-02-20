package psqldb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// FetchUserBasicDetails - userID, name, position, region, outlet, mobile, isAdmin, isActive
func (dbConn *DBImplementation) FetchUserBasicDetails(userID string, quarter, year int, args ...bool) (*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchUserBasicDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchUserBasicDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		err  error
		user = &models.UserBasicDetailsDB{}
	)

	if args != nil {
		err = dbConn.Operation.QueryRowContext(ctx, "select emp_id, name, usr_position, region, outlet, mobile_number, isAdmin, isActive from user_master_table where deleted IS NULL and emp_id=$1",
			userID).Scan(&user.UserId, &user.Name, &user.Position, &user.Region, &user.Outlet, &user.MobileNumber, &user.IsAdmin, &user.IsActive)
	} else {
		err = dbConn.Operation.QueryRowContext(ctx, "select emp_id, name, usr_position, region, outlet, mobile_number, isAdmin, isActive from user_master_table where deleted IS NULL and emp_id=$1 and quarter=$2 and year=$3",
			userID, quarter, year).Scan(&user.UserId, &user.Name, &user.Position, &user.Region, &user.Outlet, &user.MobileNumber, &user.IsAdmin, &user.IsActive)
	}

	if err == nil {
		position, ok := constant.EmployeePositionMap[user.Position]
		if !ok {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
				Details:    "FetchUserBasicDetails: incorrect user position: userId-" + user.UserId + " position-" + user.Position,
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

//TODO: FetchUserDetails - for basic user info and manger of that user
func (dbConn *DBImplementation) FetchUserDetails(userID string, quarter, year int) (*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchAgentUserDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchAgentUserDetails done")

	user, err := dbConn.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("FetchUserDetails: %v", err.Details)
		return nil, err
	}

	var query string

	switch {
	case user.Position == constant.VicePresident || user.Position == constant.VicePresidentAbb || user.Position == constant.VicePresidentCC:
		return user, nil
	case user.Position == constant.GeneralManager || user.Position == constant.GeneralManagerAbb || user.Position == constant.GeneralManagerCC:
		user.ManagerPosition = constant.VicePresidentAbb
		query = "select vice_president from user_hierarchy_table where deleted IS NULL and general_manager= $1 and quarter=$2 and year=$3 limit 1"

	case user.Position == constant.SeniorManager || user.Position == constant.SeniorManagerAbb || user.Position == constant.SeniorManagerCC:
		user.ManagerPosition = constant.GeneralManagerAbb
		query = "select general_manager from user_hierarchy_table where deleted IS NULL and senior_manager= $1 and quarter=$2 and year=$3 limit 1"

	case user.Position == constant.Manager || user.Position == constant.ManagerAbb || user.Position == constant.ManagerCC:
		user.ManagerPosition = constant.SeniorManagerAbb
		query = "select senior_manager from user_hierarchy_table where deleted IS NULL and manager= $1 and quarter=$2 and year=$3 limit 1"

	case user.Position == constant.OutletManager || user.Position == constant.OutletManagerAbb || user.Position == constant.OutletManagerCC:
		user.ManagerPosition = constant.ManagerAbb
		query = "select manager from user_hierarchy_table where deleted IS NULL and outlet_manager= $1 and quarter=$2 and year=$3 limit 1"

	case user.Position == constant.Promoter || user.Position == constant.PromoterAbb || user.Position == constant.PromoterCC || user.Position == constant.CSR ||
		user.Position == constant.CSRAbb || user.Position == constant.CSRCC:
		user.ManagerPosition = constant.OutletManagerAbb
		query = "select outlet_manager from user_hierarchy_table where deleted IS NULL and emp_id= $1 and quarter=$2 and year=$3 limit 1"

	default:
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "INVALID USER POSITION",
			Details:    "FetchUserDetails: id-" + userID + " incorrect position - " + user.Position,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	scanErr := dbConn.Operation.QueryRowContext(ctx, query, user.UserId, quarter, year).Scan(&user.ManagerId)
	if scanErr != nil {
		if strings.Contains(scanErr.Error(), "no rows") {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    errors.UserNotFound,
				Details:    "FetchUserDetails: id-" + userID + " no rows mngr - " + scanErr.Error(),
			}
		}

		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchUserDetails: id-" + userID + " err mngr - " + scanErr.Error(),
		}
	}

	return user, nil
}

func (dbConn *DBImplementation) FetchOTPValidity(userID string) (*models.OtpModel, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchOTPValidity start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchOTPValidity done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	otp := &models.OtpModel{}

	err := dbConn.Operation.QueryRowContext(ctx, "select emp_id, user_active, otp, otp_resend_count, otp_verify_failed_count, otp_used, otp_create_time, otp_update_time, last_login_time from user_authentication_table where emp_id=$1", userID).
		Scan(&otp.StaffID, &otp.IsActiveUser, &otp.OTP, &otp.OTPResendCount, &otp.OTPVerifyFailedCount, &otp.OTPUsed, &otp.OtpCreateTime, &otp.OtpLastUpdateTime, &otp.LastLoginTime)
	if err == nil {
		return otp, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return otp, nil
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchOTPValidity: err - " + err.Error(),
	}
}

func (dbConn *DBImplementation) InsertStaffOTP(update bool, otpModel *models.OtpModel) *errors.Error {

	// logger.FileLogger.Info().Interface("clientID", otpModel.StaffID).Time("start", time.Now()).Msg("InsertStaffOTP start")
	// defer logger.FileLogger.Info().Interface("clientID", otpModel.StaffID).Time("end", time.Now()).Msg("InsertStaffOTP done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var err error

	if update {
		_, err = dbConn.Operation.QueryContext(ctx, "UPDATE user_authentication_table SET user_active=$1, otp=$2, otp_resend_count=$3, otp_verify_failed_count=$4, otp_used=$5, otp_create_time=$6, otp_update_time=$7, last_login_time=$8, modified=$9, comment=$10",
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
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "UpdateStaffOTP: update failed - " + err.Error(),
			}
		}

	} else {
		_, err = dbConn.Operation.QueryContext(ctx, "INSERT INTO user_authentication_table(emp_id, user_active, otp, otp_resend_count, otp_verify_failed_count, otp_used, otp_create_time, otp_update_time, last_login_time, created, modified, comment) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);",
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
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "InsertStaffOTP: insert failed - " + err.Error(),
			}
		}
	}

	return nil
}

// FetchTotalNonManagerialEmployees - all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchTotalNonManagerialEmployees(region, outlet string, quarter, year int) (int, *errors.Error) {

	// logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("start", time.Now()).Msg("FetchTotalNonManagerialEmployees start")
	// defer logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("end", time.Now()).Msg("FetchTotalNonManagerialEmployees done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var count int

	//usr_position=='BACK OFFICE'
	err := dbConn.Operation.QueryRowContext(ctx, "select count(*) from user_master_table where deleted IS NULL and isactive=true and region=$1 and outlet=$2 and (usr_position='PROMOTER' OR usr_position='CSR') and quarter=$3 and year=$4",
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
