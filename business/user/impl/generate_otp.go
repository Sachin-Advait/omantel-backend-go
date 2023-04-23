package impl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) GenerateOTP(ctx *gin.Context, userID string) (map[string]string, *errors.Error) {

	// currTime := time.Now()
	// currYear, currMonth, _ := currTime.Date()
	// currQuarter := getQuarterFromMonth(int(currMonth))

	var isUpdate = true

	user, err := impl.DbService.FetchUserBasicDetails(userID, 0, 0, true)
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

	userPosition, ok := constant.EmployeePositionMap[user.Position]
	if !ok || user.MobileNumber == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "USER POSITION OR CONTANCT NUMBER NOT FOUND",
			Details:    "position:" + user.Position + ", contact:" + user.MobileNumber,
		}
	}

	otpModel, err := impl.DbService.FetchOTPValidity(userID)
	if err != nil {
		return nil, err
	}

	if otpModel.StaffID == "" {
		isUpdate = false
		otpModel.StaffID = userID
		otpModel.StaffName = user.Name
		otpModel.StaffPosition = user.Position
	}

	// //otpLastUpdateTime updated only while generating OTP more than block time
	// if time.Since(otpModel.OtpLastUpdateTime).Minutes() <= constant.OtpGenerateBlockTimeInMin {

	// 	//too many times sent within 15 mins
	// 	if otpModel.OTPResendCount > constant.OtpRegenerateAllowedCount {
	// 		return nil, &errors.Error{
	// 			StatusCode: http.StatusForbidden,
	// 			Code:       errors.UserBlockedCode,
	// 			Message:    errors.OTPLimitError,
	// 			Details:    "max otp resend reached",
	// 		}
	// 	}

	// 	//too many failure attempts within 15 mins
	// 	if otpModel.OTPVerifyFailedCount >= constant.MaxAllowedLoginAttempt {
	// 		return nil, &errors.Error{
	// 			StatusCode: http.StatusForbidden,
	// 			Code:       errors.UserBlockedCode,
	// 			Message:    errors.RetryLoginLimitReached,
	// 			Details:    "incorrect otp attempt reached",
	// 		}
	// 	}

	// 	otpModel.OTPResendCount++
	// } else {
	// 	//more than 15 mins hence restart counter
	// 	otpModel.OTPResendCount = 0
	// 	otpModel.OTPVerifyFailedCount = 0
	// 	otpModel.OtpLastUpdateTime = time.Now()
	// }

	// generate random num as OTP
	otpModel.OTP = 123456
	otpModel.OTPUsed = false
	otpModel.OtpCreateTime = time.Now()

	otpModel.UpdateReason = "otp generated"
	otpModel.UpdatedTime = time.Now()

	// if otpModel.OTPResendCount > 0 {
	// 	otpModel.UpdateReason = "otp re-generated"
	// }

	//store OTP in db
	if err := impl.DbService.InsertStaffOTP(isUpdate, otpModel); err != nil {
		return nil, err
	}

	// if err := impl.SOAPService.CallSOAPClientSteps(userID, user.Name, user.MobileNumber, impl.Config.OTPUrl, otpModel.OTP); err != nil {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Code:       errors.InternalServerError,
	// 		Message:    "SOAP CALL FAILED",
	// 		Details:    "sms call failed: " + err.Error(),
	// 	}
	// }

	return map[string]string{"error": "nil", "position": userPosition}, nil
}
