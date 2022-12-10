package logger

import (
	"context"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"io"
	"os"
	"path"
)

// Configuration for logging
type FileConfig struct {
	// Enable console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
	//Compress rolled files in .gz format
	Compress bool
}

type Logger struct {
	*zerolog.Logger
}

var FileLogger *Logger

// InitFileLogger sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func InitFileLogger(config FileConfig) error {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Bool("compress", config.Compress).
		Msg("logging configured")

	logfileObj := &Logger{
		Logger: &logger,
	}
	FileLogger = logfileObj
	return nil
}

// CtxTrace is the file context logging for trace log
func (f *Logger) CtxTrace(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Trace())
}

// CtxDebug is the file context logging for debug log
func (f *Logger) CtxDebug(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Debug())
}

// CtxInfo is the file context logging for info log
func (f *Logger) CtxInfo(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Info())
}

// CtxWarn is the file context logging for warn log
func (f *Logger) CtxWarn(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Warn())
}

// CtxError is the file context logging for error log
func (f *Logger) CtxError(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Error())
}

// CtxPanic is the file context logging for panic log
func (f *Logger) CtxPanic(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Panic().Stack())
}

// CtxFatal is the file context logging for fatal log
func (f *Logger) CtxFatal(ctx context.Context) *zerolog.Event {
	return withContext(ctx, f.Fatal().Stack())
}

func newRollingFile(config FileConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		Fatal(nil).Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		Compress:   config.Compress,   //compress files
	}
}
