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
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func GenerateOTP(ctx *gin.Context) {

	userReq := &models.GenerateOTP{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("generateOTP : binding request body")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: "", URL: ctx.Request.URL.String(), Type: "error", Msg: "error binding request body",
			Code: http.StatusBadRequest, Text: err.Error(), CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: errors.ErrorCodeMap[errors.RequestDecodeErrorCode],
		})
		return
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("generateOTP : req body")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: "", URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := validateGenerateOTP(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("generateOTP : invalid request")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: userReq.ClientId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invaild req", Code: http.StatusBadRequest,
			Text: err.Error(), CreateTime: time.Now().Local()})

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: err.Error(),
		})
		return
	}

	//since userid not found in token
	ctx.Set(logger.UserIDLogParam, userReq.ClientId)

	ctx.JSON(user.GenerateOTP(ctx, userReq))
}

func validateGenerateOTP(req *models.GenerateOTP) error {

	if req.ClientId == "" {
		return errors.New("BLANK USERID")
	}

	return nil
}
