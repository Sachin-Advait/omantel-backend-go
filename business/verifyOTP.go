package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"time"
)

func (impl *IncentiveImpl) VerifyOTP(ctx *gin.Context, clientID, otp int) (*models.VerifyOTPResponse, *errors.Error) {

	otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	resp := &models.VerifyOTPResponse{}

	if otp != otpModel.OTP {
		resp.Error = "otp invalid"
		return resp, nil
	}

	if !otpModel.ExpiryTime.After(time.Now().Local()) {
		resp.Error = "otp has expired"
		return resp, nil
	}

	resp.Authenticated = true
	resp.StaffName = otpModel.StaffName
	otpModel.StaffPosition = commons.EmployeePositionMap[otpModel.StaffPosition]

	return resp, nil
}
