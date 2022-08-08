package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func generateOTP(ctx *gin.Context) {

	userReq := &models.GenerateOTP{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("generateOTP : error binding request body")
		logger.LogChan <- &logger.Logging{
			UserID:     "",
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "unauthorized request",
			Text:       "auth token not found",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: errors.ErrorCodeMap[errors.RequestDecodeErrorCode],
		})
		return
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(constant.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("generateOTP : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     userReq.ClientId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	ctx.JSON(adapters.GenerateOTP(ctx, userReq))
}
