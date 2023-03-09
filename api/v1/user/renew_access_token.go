package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/user"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func RenewAccessToken(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, "request").Msg("renewAccessToken : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "trace",
		Msg:        "request body",
		Text:       "",
		CreateTime: time.Now().Local(),
	}

	if err := validateRenewToken(ctxHeader.XUserId, ctxHeader.SuperAuthorization); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("renewAccessToken : invalid request")
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

	ctx.JSON(user.RenewAccessToken(ctx))
}

func validateRenewToken(userID, renewToken string) error {

	if userID == "" || renewToken == "" {
		return errors.New("BLANK USERID OR TOKEN")
	}

	return nil
}
