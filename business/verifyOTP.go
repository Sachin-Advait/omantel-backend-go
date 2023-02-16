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

func (impl *IncentiveImpl) VerifyOTP(ctx *gin.Context, req *models.VerifyOTP) (*models.VerifyOTPResponse, *errors.Error) {

	// currTime := time.Now()
	// currYear, currMonth, _ := currTime.Date()
	// currQuarter := impl.getQuarterFromMonth(int(currMonth))

	otpModel, err := impl.PsqlDb.FetchOTPValidity(req.ClientId)
	if err != nil {
		return nil, err
	}

	if otpModel.StaffID == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.OTPInvalidCode,
			Message:    "Please Generate OTP",
			Details:    "otp not found",
		}
	}

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
	if req.Otp != otpModel.OTP {
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

	user, err := impl.PsqlDb.FetchUser(req.ClientId)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    "YOU ARE NOT AUTHORIZED TO VIEW REPORTS",
			Details:    "position:" + user.Position + "contact:" + user.MobileNumber,
		}
	}

	//store OTP in db
	if err := impl.PsqlDb.InsertStaffOTP(true, otpModel); err != nil {
		return nil, err
	}

	accessTokenExpiry := impl.Config.SessionExpiryTime
	refreshTokenExpiry := impl.Config.RefreshTokenExpiryTime

	var (
		jwtErr       error
		accessToken  string
		refreshToken string
	)

	if user.IsAdmin {
		accessToken, _, jwtErr = jwthelper.CreateAccessToken(req.ClientId, accessTokenExpiry, jwthelper.AdminJwtSecretKey, true)
	} else {
		accessToken, _, jwtErr = jwthelper.CreateAccessToken(req.ClientId, accessTokenExpiry, jwthelper.JwtSecretKey, false)
	}

	if jwtErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "create at failed - " + err.Error(),
		}
	}

	if user.IsAdmin {
		refreshToken, _, jwtErr = jwthelper.CreateRefreshToken(req.ClientId, refreshTokenExpiry, jwthelper.AdminJwtSecretKey, true)
	} else {
		refreshToken, _, jwtErr = jwthelper.CreateRefreshToken(req.ClientId, refreshTokenExpiry, jwthelper.JwtSecretKey, false)
	}

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
		IsAdmin:       user.IsAdmin,
		StaffName:     user.Name,
		StaffPosition: constant.EmployeePositionMap[user.Position],
		StaffOutlet:   user.Outlet,
		StaffRegion:   user.Region,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		LastLoginTime: time.Now(),
	}

	return resp, nil
}
