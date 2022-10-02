package business

import (
	"crypto/rand"
	"io"
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
	otpModel.StaffID = clientID

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

	//otpLastUpdateTime updated only while generating OTP more than block time
	if time.Since(otpModel.OtpLastUpdateTime).Minutes() <= constant.OtpGenerateBlockTimeInMin {

		//too many times sent within 15 mins
		if otpModel.OTPResendCount > constant.OtpRegenerateAllowedCount {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.OTPLimitError,
				Details:    "max otp resend reached",
			}
		}

		//too many failure attempts within 15 mins
		if otpModel.OTPVerifyFailedCount >= constant.MaxAllowedLoginAttempt {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.RetryLoginLimitReached,
				Details:    "incorrect otp attempt reached",
			}
		}

		otpModel.OTPResendCount++
	} else {
		//more than 15 mins hence restart counter
		otpModel.OTPResendCount = 0
		otpModel.OTPVerifyFailedCount = 0
		otpModel.OtpLastUpdateTime = time.Now()
	}

	//TODO: generate random num as OTP
	otpModel.OTP = 123456
	otpModel.OTPUsed = false
	otpModel.OtpCreateTime = time.Now()

	otpModel.UpdateReason = "otp generated"
	otpModel.UpdatedTime = time.Now()

	if otpModel.OTPResendCount > 0 {
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

func generateOTP(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
