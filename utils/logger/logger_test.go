package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	InitLogger(TraceLevel)
	InitLogger(DebugLevel)
	InitLogger(InfoLevel)
	InitLogger(WarnLevel)
	InitLogger(ErrorLevel)
	InitLogger(PanicLevel)
	InitLogger(FatalLevel)
	Trace(nil)
	Debug(nil)
	Info(nil)
	Warn(nil)
	Error(nil)
	Panic(nil)
	Fatal(nil)
}
