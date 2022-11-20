package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func init() {

	initEnvConfigs()
	initLogger()
	initFileLogger()
	initHTTPClient()
}

func main() {

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("application is running: config loaded successfully - " + time.Now().String())
	}
}

func initEnvConfigs() {
	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())
	logger.Info(context.TODO()).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Msg("initialized configs for env")
}

func initLogger() {

	loggerConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Stack().Err(err).Msg("error loading logger config file, env - " + flags.Env())
	}

	logger.InitLogger(logger.Level(loggerConfig.GetString(constant.LogLevelKey)))
}

func initFileLogger() {

	applicationConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Stack().Err(err).Msg("error loading application config file, env - " + flags.Env())
	}

	err = logger.InitFileLogger(logger.FileConfig{
		EncodeLogsAsJson:      applicationConfig.GetBool(constant.EncodeLogAsJsonKey),
		ConsoleLoggingEnabled: applicationConfig.GetBool(constant.ConsoleLoggingEnabledKey),
		FileLoggingEnabled:    applicationConfig.GetBool(constant.FileLoggingEnabledKey),
		Directory:             applicationConfig.GetString(constant.DirectoryKey),
		Filename:              applicationConfig.GetString(constant.FileNameKey),
		MaxSize:               applicationConfig.GetInt(constant.MaxSizeInMBKey),
		MaxBackups:            applicationConfig.GetInt(constant.MaxBackUpsKey),
		MaxAge:                applicationConfig.GetInt(constant.MaxAgeInDaysKey),
		Compress:              applicationConfig.GetBool(constant.Compress),
	})

	config := logger.FileConfig{
		EncodeLogsAsJson:      applicationConfig.GetBool(constant.EncodeLogAsJsonKey),
		ConsoleLoggingEnabled: applicationConfig.GetBool(constant.ConsoleLoggingEnabledKey),
		FileLoggingEnabled:    applicationConfig.GetBool(constant.FileLoggingEnabledKey),
		Directory:             applicationConfig.GetString(constant.DirectoryKey),
		Filename:              applicationConfig.GetString(constant.FileNameKey),
		MaxSize:               applicationConfig.GetInt(constant.MaxSizeInMBKey),
		MaxBackups:            applicationConfig.GetInt(constant.MaxBackUpsKey),
		MaxAge:                applicationConfig.GetInt(constant.MaxAgeInDaysKey),
		Compress:              applicationConfig.GetBool(constant.Compress),
	}

	if err != nil {
		logger.Fatal(context.TODO()).Err(err).Msg("unable to initialize file logger")
	}

	logger.Info(context.TODO()).Interface("fileConfig", config).Msg("fileConfig: load configuration")
}

func initHTTPClient() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("initHTTPClient: error loading application config file, env - " + flags.Env())
	}

	connectTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPConnectTimeoutInMillisKey)
	keepAliveDuration := time.Millisecond * applicationConfig.GetDuration(constant.HTTPKeepAliveDurationInMillisKey)
	maxIdleConnections := applicationConfig.GetInt(constant.HTTPMaxIdleConnectionsKey)
	idleConnectionTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPIdleConnectionTimeoutInMillisKey)
	tLSHandshakeTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPTlsHandshakeTimeoutInMillisKey)
	expectContinueTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPExpectContinueTimeoutInMillisKey)
	timeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPTimeoutInMillisKey)

	dialer := &net.Dialer{
		Timeout:   connectTimeout,
		KeepAlive: keepAliveDuration,
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext:           dialer.DialContext,
			MaxIdleConns:          maxIdleConnections,
			IdleConnTimeout:       idleConnectionTimeout,
			TLSHandshakeTimeout:   tLSHandshakeTimeout,
			ExpectContinueTimeout: expectContinueTimeout,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
		Timeout: timeout,
	}

	logger.Info(context.TODO()).Dur("timeout", timeout).Interface("client", client).Msg("httpConfig: load configuration")
}
