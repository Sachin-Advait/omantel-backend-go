package middlewares

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"

	"github.com/gin-gonic/gin"
)

// LoggerMiddlewareOptions is the set of configurable allowed for log
type LoggerMiddlewareOptions struct {
	NotLogQueryParams  bool
	NotLogHeaderParams bool
}

// Logger is the middlewares to be used for logging the request and response information
// This should be the first middlewares to be added, in case the recovery middlewares is not being used.
// Otherwise, it should be the second one, just after the recovery middlewares.
func Logger(options LoggerMiddlewareOptions) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		start := time.Now()

		id := ctx.GetHeader(constant.XRequestID)

		if id == "" {
			// get a unique id
			uid, err := objecthelper.GenerateUUID()
			if err == nil {
				id = uid.String()
			}
		}
		// apply the id in the context
		ctx.Set(logger.RequestIDLogParam, id)
		ctx.Request.Header.Add(constant.XRequestID, id)

		path := ctx.Request.URL.Path
		// apply the pah in the context
		ctx.Set(logger.PathLogParam, path)

		userId := ctx.GetHeader(constant.XUserId)
		// apply the userID in the context
		ctx.Set(logger.UserIDLogParam, userId)

		// apply the method type in the context
		ctx.Set(logger.MethodLogParam, ctx.Request.Method)

		// Process request
		ctx.Next()

		logger.FileLogger.CtxInfo(ctx).Dur(logger.LatencyLogParam, time.Duration(time.Since(start).Milliseconds())).Interface(logger.ReqHeaderLogParams, ctx.Request.Header).
			Interface(logger.ResHeaderLogParams, ctx.Writer.Header()).Msg("Header Logging Middleware")
	}
}
