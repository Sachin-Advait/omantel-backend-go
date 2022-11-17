package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
)

func (impl *IncentiveImpl) VerifyOTP(ctx *gin.Context, clientID string, otp int) (*models.VerifyOTPResponse, *errors.Error) {

	_, otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	if otpModel.StaffID == "" {
		name, _, err := impl.PsqlDb.FetchUserName(clientID)
		if err != nil {
			return nil, err
		}

		if name == "" {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.UserNotFoundCode,
				Message:    errors.UserIDInvalid,
				Details:    "userId is invalid",
			}
		}

		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.OTPInvalidCode,
			Message:    "Please Generate OTP",
			Details:    "otp not found",
		}
	}

	// if !otpModel.IsActiveUser {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusForbidden,
	// 		Code:       errors.UserBlockedCode,
	// 		Message:    errors.UserBlockedError,
	// 		Details:    "user is blocked",
	// 	}
	// }

	//TODO: reset loginAttemptCount
	// if otpModel.LoginAttemptCount > constant.MaxAllowedLoginAttempt {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusForbidden,
	// 		Code:       errors.RetryLoginLimitReachedCode,
	// 		Message:    errors.RetryLoginLimitReached,
	// 		Details:    "retry login limit reached",
	// 	}
	// }

	//otp generated more than 30 mins ago
	if time.Since(otpModel.OtpCreateTime).Minutes() >= constant.OtpTimeoutInMin {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.OTPExpiredCode,
			Message:    errors.OTPExpiredError,
			Details:    "otp is expired",
		}
	}

	//too many login attempts within 15 mins
	if otpModel.OTPVerifyFailedCount > constant.MaxAllowedLoginAttempt {
		return nil, &errors.Error{
			StatusCode: http.StatusForbidden,
			Code:       errors.UserBlockedCode,
			Message:    errors.RetryLoginLimitReached,
			Details:    "user is blocked, max failed count",
		}
	}

	//otp has been used
	if otpModel.OTPUsed {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.OTPUsedCode,
			Message:    errors.OTPUsedError,
			Details:    "otp is used",
		}
	}

	//otp unmatched
	if otp != otpModel.OTP {
		otpModel.OTPVerifyFailedCount++
		otpModel.IsActiveUser = false
		otpModel.UpdateReason = "incorrect otp"
		otpModel.UpdatedTime = time.Now()

		//store OTP in db
		if err := impl.PsqlDb.InsertStaffOTP(true, otpModel); err != nil {
			return nil, err
		}

		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.OTPInvalidCode,
			Message:    errors.OTPInvalidError,
			Details:    "otp is incorrect",
		}
	}

	//otp matched
	otpModel.OTPUsed = true
	otpModel.OTPVerifyFailedCount = 0
	otpModel.IsActiveUser = true
	otpModel.LastLoginTime = time.Now()
	otpModel.UpdateReason = "otp login success"
	otpModel.UpdatedTime = time.Now()

	//store OTP in db
	if err := impl.PsqlDb.InsertStaffOTP(true, otpModel); err != nil {
		return nil, err
	}

	accessTokenExpiry := impl.Config.SessionExpiryTime
	refreshTokenExpiry := impl.Config.RefreshTokenExpiryTime

	accessToken, _, jwtErr := jwthelper.CreateAccessToken(clientID, accessTokenExpiry, jwthelper.JwtSecretKey)
	if jwtErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "create at failed - " + err.Error(),
		}
	}

	refreshToken, _, jwtErr := jwthelper.CreateRefreshToken(clientID, refreshTokenExpiry, jwthelper.JwtSecretKey)
	if jwtErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "create rt failed - " + err.Error(),
		}
	}

	resp := &models.VerifyOTPResponse{
		Authenticated: true,
		StaffName:     otpModel.StaffName,
		StaffPosition: constant.EmployeePositionMap[otpModel.StaffPosition],
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		LastLoginTime: time.Now(),
	}

	return resp, nil
}
