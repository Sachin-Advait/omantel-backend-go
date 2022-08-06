package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GenerateOTP(ctx *gin.Context, clientID int) (map[string]string, *errors.Error) {

	var isUpdate = true

	otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	if otpModel.StaffID == 0 {
		name, err := impl.PsqlDb.FetchUserName(clientID)
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

		isUpdate = false
	}

	// contactNumber, err := impl.PsqlDb.FetchEmployeeContactNumber(clientID)
	// if err != nil {
	// 	return nil, err
	// }

	// if otpModel.StaffID != 0 && !otpModel.IsActiveUser {
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

	if time.Now().Sub(otpModel.OtpLastUpdateTime).Minutes() < constant.OtpGenerateBlockTimeInMin {

		if otpModel.OTPResendCount >= constant.OtpRegenerateAllowedCount {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.UserBlockedError,
				Details:    "max otp resend reached",
			}
		}

		if otpModel.OTPVerifyFailedCount >= constant.MaxAllowedLoginAttempt {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.UserBlockedError,
				Details:    "incorrect otp attempt reached",
			}
		}
	}

	if otpModel.OTP == 0 || otpModel.OTPUsed {
		otpModel.OTP = 123456
		otpModel.OTPUsed = false
	}

	otpModel.OTPResendCount++
	otpModel.StaffID = clientID
	otpModel.IsActiveUser = false
	otpModel.OtpLastUpdateTime = time.Now()
	otpModel.OtpCreateTime = time.Now()
	otpModel.UpdateReason = "otp generated"
	otpModel.UpdatedTime = time.Now()

	if otpModel.OTPResendCount > 1 {
		otpModel.UpdateReason = "otp re-generated"
	}

	//store OTP in db
	if err := impl.PsqlDb.InsertStaffOTP(isUpdate, otpModel); err != nil {
		return nil, err
	}

	//reqSMS := &restapi.SMSParams{
	//	HeaderRequestID:             requestID,
	//	HeaderStaffID:               strconv.Itoa(clientID),
	//	HeaderTimeStamp:             otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	//	NotificationDateTime:        otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	//	NotificaitionValidityPeriod: otpModel.CreatedTime.Add(3 * time.Minute).Format("2006-01-02 15:04:05"),
	//	Message:                     "Your OTP is: " + strconv.Itoa(otpModel.OTP),
	//	ContactNumber:               strconv.Itoa(contactNumber),
	//}
	//
	//if err = restapi.SendOTP(reqSMS); err != nil {
	//	return err
	//}

	return map[string]string{"error": "nil"}, nil
}
