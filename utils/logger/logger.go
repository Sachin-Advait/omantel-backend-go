package logger

import (
	"context"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log Params
const (
	RequestIDLogParam   = "requestId"
	UserIDLogParam      = "userId"
	StatusCodeLogParam  = "statusCode"
	EndTimeLogParam     = "endTime"
	StartTimeLogParam   = "startTime"
	PathLogParam        = "path"
	QueryLogParam       = "query"
	HeaderLogParams     = "header"
	LatencyLogParam     = "latency"
	ClientIPLogParam    = "clientIP"
	MethodLogParam      = "method"
	BodyLogParam        = "body"
	ErrorLogParam       = "error"
	ResponseParam       = "response"
	CorrelationLogParam = "correlationId"
)

type Level string

func (l Level) zeroLogLevel() zerolog.Level {
	switch l {
	case constant.TraceLevel:
		return zerolog.TraceLevel
	case constant.DebugLevel:
		return zerolog.DebugLevel
	case constant.InfoLevel:
		return zerolog.InfoLevel
	case constant.WarnLevel:
		return zerolog.WarnLevel
	case constant.ErrorLevel:
		return zerolog.ErrorLevel
	case constant.FatalLevel:
		return zerolog.FatalLevel
	case constant.PanicLevel:
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}

// InitLogger is used to initialize console logger
func InitLogger(level Level) {
	zerolog.ErrorStackMarshaler = getErrorStackMarshaller()
	zerolog.SetGlobalLevel(level.zeroLogLevel())
	log.Logger = log.With().Caller().Logger()
}

// Trace is the console context logging for trace log
func Trace(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Trace())
}

// Debug is the console context logging for debug log
func Debug(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Debug())
}

// Info is the console context logging for info log
func Info(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Info())
}

// Warn is the console context logging for warn log
func Warn(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Warn())
}

// Error is the console context logging for error log
func Error(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Error().Stack())
}

// Panic is the console context logging for panic log
func Panic(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Panic().Stack())
}

// Fatal is the console context logging for fatal log
func Fatal(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Fatal().Stack())
}

func getErrorStackMarshaller() func(err error) interface{} {
	return func(err error) interface{} {
		return string(debug.Stack())
	}
}

func withIDAndPath(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if ctx == nil {
		return event
	}
	id := ctx.Value(RequestIDLogParam)
	if id != nil {
		event.Interface(RequestIDLogParam, id)
	}
	path := ctx.Value(PathLogParam)
	if path != nil {
		event.Interface(PathLogParam, path)
	}
	return event
}
