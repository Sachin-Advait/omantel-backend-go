package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/user"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func VerifyOTP(ctx *gin.Context) {

	userReq := &models.VerifyOTP{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("verifyOTP : binding request body")
		logger.LogChan <- &logger.Logging{
			UserID:     "",
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "request binding",
			Code:       http.StatusBadRequest,
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: errors.ErrorCodeMap[errors.RequestDecodeErrorCode],
		})
		return
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("verifyOTP : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     userReq.ClientId,
		URL:        ctx.Request.URL.String(),
		Type:       "trace",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	if err := validateVerifyOTP(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("verifyOTP : invalid request")
		logger.LogChan <- &logger.Logging{
			UserID:     "",
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invaild req",
			Code:       http.StatusBadRequest,
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: err.Error(),
		})
		return
	}

	//since userid not found in token
	ctx.Set(logger.UserIDLogParam, userReq.ClientId)

	ctx.JSON(user.VerifyOTP(ctx, userReq))
}

func validateVerifyOTP(req *models.VerifyOTP) error {

	if req.ClientId == "" || req.Otp == 0 {
		return errors.New("BLANK USERID OR OTP")
	}

	return nil
}
