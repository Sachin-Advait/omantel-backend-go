package commons

import (
	"sync"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var once sync.Once
var log *zap.SugaredLogger

func SetUpLogging() {
	l, _ := zap.NewProduction()
	log = l.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	once.Do(func() {
		SetUpLogging()
		defer logrus.Info("log setup completed")
	})

	return log
}
