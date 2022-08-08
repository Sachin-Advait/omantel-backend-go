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

// GetFeeds godoc
// @Summary Fetches User Feeds
// @Description Fetches Feeds based on users interests
// @ID getFeedsV1
// @Tags Feed V1
// @Produce json
// @Param X-UserType header string true "X-UserType Header"
// @Param X-User-Id header string true "X-User-Id Header"
// @Param X-SourceID header string true "X-SourceID Header"
// @Param X-DeviceID header string true "X-DeviceID Header"
// @Param X-OperatingSystem header string false "X-OperatingSystem Header"
// @Param X-Request-Id header string false "X-Request-Id Header"
// @Success 200 {object} models.APIResponse{data=[]models.GetFeedResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Security AuthorizationBearerToken
// @Router /feed/v1/getFeeds [get]
func getSingleUserIncentive(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Stack().Msg("getSingleUserIncentive : unauthorized request")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "unauthorized request",
			Text:       "auth token not found",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	userReq := &models.SingleUserIncentiveRequest{}

	if err := ctx.ShouldBindJSON(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("getSingleUserIncentive : error binding request body")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "error binding request body",
			Text:       err.Error(),
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error: errors.ErrorCodeMap[errors.RequestDecodeErrorCode],
		})
		return
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(constant.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("getSingleUserIncentive : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	ctx.JSON(adapters.GetSingleUserIncentive(ctx, userReq))
}
