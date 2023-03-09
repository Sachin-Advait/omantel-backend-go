package middlewares

import (
	"net/http"
	"strings"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"

	"github.com/gin-gonic/gin"
)

// ContextCreator:
//	* fetches values from header and sets into contextHeader in upperCase
//	* sets the context Header in the request context with key as 'Context-Header'
func ContextCreator() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		ctxHeader := context.ContextHeader{}

		if err := ctx.ShouldBindHeader(&ctxHeader); err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("ContextCreator Middleware - error binding request header")
			ctx.JSON(http.StatusBadRequest, models.APIResponse{
				ErrorCode: errors.InvalidReqParameterCode,
				Message:   errors.InvalidReqParameters,
				Data:      "INCORRECT HEADER FOUND",
			})
			ctx.Abort()
			return
		}

		ctxHeader.XAppID = strings.ToUpper(ctxHeader.XAppID)
		ctxHeader.XAppVersion = strings.ToUpper(ctxHeader.XAppVersion)
		ctxHeader.XDeviceID = strings.ToUpper(ctxHeader.XDeviceID)
		ctxHeader.XDeviceType = strings.ToUpper(ctxHeader.XDeviceType)
		ctxHeader.XMACAddress = strings.ToUpper(ctxHeader.XMACAddress)
		ctxHeader.XOSVersion = strings.ToUpper(ctxHeader.XOSVersion)
		ctxHeader.XOperatingSystem = strings.ToUpper(ctxHeader.XOperatingSystem)
		ctxHeader.XPlatform = strings.ToUpper(ctxHeader.XPlatform)
		ctxHeader.XProductVersion = strings.ToUpper(ctxHeader.XProductVersion)
		ctxHeader.XSystemInfo = strings.ToUpper(ctxHeader.XSystemInfo)

		ctx.Set(constant.HeaderInContext, ctxHeader)
		//ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
