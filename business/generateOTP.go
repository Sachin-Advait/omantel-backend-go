package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"net/http"
	"strconv"
	"time"
)

func (impl *IncentiveImpl) GenerateOTP(ctx *gin.Context, clientID int) (map[string]string, *errors.Error) {

	otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	contactNumber, err := impl.PsqlDb.FetchEmployeeContactNumber(clientID)
	if err != nil {
		return nil, err
	}

	fmt.Println(contactNumber)

	if !otpModel.ExpiryTime.After(time.Now().Local()) {

		otpString, err := impl.generateOTP(6)
		if err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "generate otp failed",
				Details:    "generate otp failed: " + err.Error(),
			}
		}

		otpString = "123456"

		otpModel.OTP, _ = strconv.Atoi(otpString)
		otpModel.Retry = 3
		otpModel.StaffID = clientID
		otpModel.CreatedTime = time.Now().Local()
		otpModel.ExpiryTime = otpModel.CreatedTime.Add(5 * time.Minute)

		//store OTP in db
		if err := impl.PsqlDb.InsertStaffOTP(otpModel); err != nil {
			return nil, err
		}
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
