package middlewares

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
	"time"

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

		id := ctx.GetHeader(constant.XRequestID)

		if id == "" {
			// get a unique id
			uid, err := objecthelper.GenerateUUID()
			if err == nil {
				id = uid.String()
			}
		}
		// apply the id in the context
		ctx.Set(constant.RequestIDLogParam, id)
		ctx.Request.Header.Add(constant.XRequestID, id)

		path := ctx.Request.URL.Path
		// apply the pah in the context
		ctx.Set(constant.PathLogParam, path)

		userId := ctx.GetHeader(constant.XUserId)
		// apply the userID in the context
		ctx.Set(constant.UserIDLogParam, userId)

		// Start timer
		start := time.Now()

		// log the initial request
		event := logger.Info(ctx).
			Str(constant.UserIDLogParam, userId).
			Time(constant.StartTimeLogParam, start).
			Str(constant.MethodLogParam, ctx.Request.Method)

		if !options.NotLogQueryParams {
			event.Str(constant.QueryLogParam, ctx.Request.URL.RawQuery)
		}

		if !options.NotLogHeaderParams {
			event.Interface(constant.HeaderLogParams, ctx.Request.Header)
		}

		event.Send()

		// Process request
		ctx.Next()

		// stop timer
		end := time.Now()
		latency := end.Sub(start)

		if logger.FileLogger != nil {
			// log the final response details in file
			logger.FileLogger.CtxInfo(ctx).
				Str(constant.UserIDLogParam, userId).
				Int(constant.StatusCodeLogParam, ctx.Writer.Status()).
				Str(constant.MethodLogParam, ctx.Request.Method).
				Str(constant.QueryLogParam, ctx.Request.URL.RawQuery).
				Interface(constant.HeaderLogParams, ctx.Request.Header).
				Time(constant.StartTimeLogParam, start).
				Time(constant.EndTimeLogParam, end).
				Dur(constant.LatencyLogParam, latency).
				Str(constant.ClientIPLogParam, ctx.ClientIP()).
				Str(constant.ErrorLogParam, ctx.Errors.ByType(gin.ErrorTypePrivate).String()).
				Send()
		} else {
			// log the final response details in console
			logger.Info(ctx).
				Str(constant.UserIDLogParam, userId).
				Int(constant.StatusCodeLogParam, ctx.Writer.Status()).
				Str(constant.MethodLogParam, ctx.Request.Method).
				Str(constant.QueryLogParam, ctx.Request.URL.RawQuery).
				Interface(constant.HeaderLogParams, ctx.Request.Header).
				Time(constant.StartTimeLogParam, start).
				Time(constant.EndTimeLogParam, end).
				Dur(constant.LatencyLogParam, latency).
				Str(constant.ClientIPLogParam, ctx.ClientIP()).
				Str(constant.ErrorLogParam, ctx.Errors.ByType(gin.ErrorTypePrivate).String()).
				Send()
		}
	}
}

func HeaderLogging() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()

		headers := make(map[string][]string, 0)

		for key, val := range ctx.Request.Header {

			if key == constant.Authorization && len(val[0]) > 50 {
				out := []rune(val[0][7:])
				out = append(out[:len(out)-10], []rune("XXXXXXXXXX")...)

				headers[key] = []string{string(out)}
				continue
			}

			headers[key] = val
		}

		//Process request
		ctx.Next()

		logger.FileLogger.CtxInfo(ctx).Dur(constant.LatencyKey, time.Now().Sub(start)).Interface(constant.ReqHeader, headers).
			Interface(constant.ResHeader, ctx.Writer.Header()).Msg("Header Logging Middleware")
	}
}