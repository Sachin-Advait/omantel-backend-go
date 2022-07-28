package logger

import (
	"testing"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
)

func TestInitLogger(t *testing.T) {
	InitLogger(constant.TraceLevel)
	InitLogger(constant.DebugLevel)
	InitLogger(constant.InfoLevel)
	InitLogger(constant.WarnLevel)
	InitLogger(constant.ErrorLevel)
	InitLogger(constant.PanicLevel)
	InitLogger(constant.FatalLevel)
	Trace(nil)
	Debug(nil)
	Info(nil)
	Warn(nil)
	Error(nil)
	Panic(nil)
	Fatal(nil)
}
