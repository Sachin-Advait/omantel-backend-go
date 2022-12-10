package psqldb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (dbConn *DBImplementation) FetchUserName(userID string) (string, *errors.Error) {

	logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchUserName start")
	defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchUserName done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var name string

	err := dbConn.Operation.QueryRowContext(ctx, "select name from user_master_table where emp_id=$1", userID).Scan(&name)
	if err == nil {
		return name, nil
	}

	if strings.Contains(err.Error(), "no rows") {
		return "", &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameterCode,
			Message:    errors.UserNotFound,
			Details:    "FetchUserName: no rows - " + err.Error(),
		}
	}

	return "", &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "FetchUserName: err - " + err.Error(),
	}
}

// FetchUserBasicDetails - userID, name, position, region, outlet, mobile, isAdmin
func (dbConn *DBImplementation) FetchUserBasicDetails(userID string) (*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchUserBasicDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchUserBasicDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.UserBasicDetailsDB{}

	err := dbConn.Operation.QueryRowContext(ctx, "select emp_id, name, usr_position, region, outlet, mobile_number, isAdmin from user_master_table where emp_id=$1", userID).
		Scan(&user.UserId, &user.Name, &user.Position, &user.Region, &user.Outlet, &user.MobileNumber, &user.IsAdmin)
	if err == nil {
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

//TODO: FetchUserDetails - for basic user info and outlet-manger of agent
func (dbConn *DBImplementation) FetchUserDetails(userID string) (*models.UserBasicDetailsDB, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchAgentUserDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchAgentUserDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.UserBasicDetailsDB{
		ManagerPosition: constant.OutletManager,
	}

	err := dbConn.Operation.QueryRowContext(ctx,
		"select a.emp_id, a.name, a.usr_position, a.mobile_number, a.region, a.outlet, a.isAdmin, b.outlet_manager from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id=b.emp_id where a.emp_id=$1",
		userID).Scan(&user.UserId, &user.Name, &user.Position, &user.MobileNumber, &user.Region, &user.Outlet, &user.IsAdmin, &user.ManagerId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameterCode,
				Message:    errors.UserNotFound,
				Details:    "FetchAgentUserDetails: no rows - " + err.Error(),
			}
		}

		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchAgentUserDetails: err - " + err.Error(),
		}
	}

	dbConn.Operation.QueryRowContext(ctx, "select usr_position from user_master_table where a.emp_id=$1", user.ManagerId).Scan(&user.ManagerPosition)

	return user, nil
}

func (dbConn *DBImplementation) FetchOTPValidity(userID string) (*models.OtpModel, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchOTPValidity start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchOTPValidity done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "UpdateStaffOTP: update failed - " + err.Error(),
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
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "InsertStaffOTP: insert failed - " + err.Error(),
			}
		}
	}

	return nil
}

// FetchTotalNonManagerialEmployees - all non-managerial employees count of a given region and outlet
func (dbConn *DBImplementation) FetchTotalNonManagerialEmployees(region, outlet string) (int, *errors.Error) {

	// logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("start", time.Now()).Msg("FetchTotalNonManagerialEmployees start")
	// defer logger.FileLogger.Info().Str("region", region).Str("outlet", outlet).Time("end", time.Now()).Msg("FetchTotalNonManagerialEmployees done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	//usr_position=='BACK OFFICE'
	err := dbConn.Operation.QueryRowContext(ctx, "select count(*) from user_master_table where region=$1 and outlet=$2 and (usr_position='PROMOTER' OR usr_position='CSR')",
		region, outlet).Scan(&count)
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
