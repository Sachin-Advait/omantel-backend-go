package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func captchaVerification(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq := &models.CaptchRequest{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("captchaVerification : binding request body")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
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
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("captchaVerification : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "trace",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	if err := validateCaptcha(ctxHeader.XUserId, userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("captchaVerification : invalid request")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
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
	ctx.Set(logger.UserIDLogParam, ctxHeader.XUserId)

	ctx.JSON(adapters.CaptchaVerification(ctx, userReq))
}

func validateCaptcha(userId string, req *models.CaptchRequest) error {

	if userId == "" || req.Token == "" {
		return errors.New("BLANK CAPTCHA TOKEN OR USERID")
	}

	return nil
}
