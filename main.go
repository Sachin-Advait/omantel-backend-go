package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminApi "gitlab.com/rendezvousPoint/incentive_calculator/adapters/admin"
	healthApi "gitlab.com/rendezvousPoint/incentive_calculator/adapters/health"
	schedulerAdp "gitlab.com/rendezvousPoint/incentive_calculator/adapters/scheduler"
	userApi "gitlab.com/rendezvousPoint/incentive_calculator/adapters/user"
	"gitlab.com/rendezvousPoint/incentive_calculator/api"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/scheduler"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/cache"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/fcm"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/http"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
	middlewares "gitlab.com/rendezvousPoint/incentive_calculator/utils/middlewares"
)

func init() {

	go handleSigTerm()

	initEnvConfigs()
	initLogger()
	initFileLogger()

	initHTTPClient()
	initDatabase()
	initDbLogger()
	initCache()
	initFCM()
}

func main() {

	defer closeDatabase()
	defer closeCache()
	defer closeFCM()

	components := LoadApplicationComponents()
	healthApi.InitHealthCheck(components.HealthProvider)
	adminApi.InitAdminAPI(components.AdminProvider)
	userApi.InitUserAPI(components.UserProvider)
	schedulerAdp.InitScheduler(components.Scheduler)

	startCron()
	startRouter()
}

func handleSigTerm() {
	done := make(chan os.Signal, 1)
	signal.Notify(done,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	sig := <-done

	switch sig {
	case syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT:
		closeDatabase()
		closeCache()
		os.Exit(0)
	}
}

func initEnvConfigs() {
	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())
	logger.Info(context.TODO()).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Msg("initialized configs for env")
}

func initLogger() {

	loggerConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Err(err).Msg("error loading logger config file, env - " + flags.Env())
	}

	logger.InitLogger(logger.Level(loggerConfig.GetString(constant.LogLevelKey)))
}

func initFileLogger() {

	applicationConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Err(err).Msg("error loading application config file, env - " + flags.Env())
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

	if err != nil {
		logger.Fatal(context.TODO()).Err(err).Msg("unable to initialize file logger")
	}
}

func initHTTPClient() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("HTTP: error loading application config file, env - " + flags.Env())
	}

	//http config load
	httpConfig := &models.HttpClientConfig{
		CertificatePath:       applicationConfig.GetString(constant.HTTPTLSCertificateFileLocation),
		ConnectionTimeout:     time.Millisecond * applicationConfig.GetDuration(constant.HTTPConnectTimeoutInMillisKey),
		KeepAliveDuration:     time.Millisecond * applicationConfig.GetDuration(constant.HTTPKeepAliveDurationInMillisKey),
		MaxIdleConnections:    applicationConfig.GetInt(constant.HTTPMaxIdleConnectionsKey),
		IdleConnectionTimeout: time.Millisecond * applicationConfig.GetDuration(constant.HTTPIdleConnectionTimeoutInMillisKey),
		TLSHandshakeTimeout:   time.Millisecond * applicationConfig.GetDuration(constant.HTTPTlsHandshakeTimeoutInMillisKey),
		ExpectContinueTimeout: time.Millisecond * applicationConfig.GetDuration(constant.HTTPExpectContinueTimeoutInMillisKey),
		Timeout:               time.Millisecond * applicationConfig.GetDuration(constant.HTTPTimeoutInMillisKey),
	}

	http.InitHttpClient(httpConfig)
}

func initDatabase() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("DB: error loading application config file, env - " + flags.Env())
	}

	// pswd, err := base64.StdEncoding.DecodeString(applicationConfig.GetString(constant.PsqlPassword))
	// if err != nil {
	// 	logger.FileLogger.Fatal().Err(err).Msg("DB: error decoding password, env - " + flags.Env())
	// }

	//db config load
	dbConfig := &models.DbClientConfig{
		Host:               applicationConfig.GetString(constant.PsqlHost),
		Port:               applicationConfig.GetInt(constant.PsqlPort),
		DbName:             applicationConfig.GetString(constant.PsqlDbName),
		UserName:           applicationConfig.GetString(constant.PsqlUser),
		Password:           applicationConfig.GetString(constant.PsqlPassword),
		MaxOpenConnections: applicationConfig.GetInt(constant.PsqlMaxOpenConn),
		MaxIdleConnections: applicationConfig.GetInt(constant.PsqlMaxIdleConn),
		MaxLifetime:        time.Duration(applicationConfig.GetInt(constant.ConnMaxLifetime)) * time.Hour,
		MaxIdletime:        time.Duration(applicationConfig.GetInt(constant.ConnMaxIdleTime)) * time.Minute,
	}

	database.InitDatabase(dbConfig)
}

func initDbLogger() {
	dblogger.InitDBLogger()
}

func closeDatabase() {
	err := database.Close()
	if err != nil {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("error closing database")
	}

	logger.FileLogger.Info().Msg("database graceful shut down")
}

func initCache() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("Cache: error loading application config file, env - " + flags.Env())
	}

	pswd, err := base64.StdEncoding.DecodeString(applicationConfig.GetString(constant.UserIncentiveRedisPassword))
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("Cache: error decoding password, env - " + flags.Env())
	}

	redisConfig := &models.RedisConfigStruct{
		Address:                     applicationConfig.GetString(constant.UserIncentiveRedisAddress),
		Password:                    string(pswd),
		Db:                          applicationConfig.GetInt(constant.UserIncentiveRedisDb),
		PoolSize:                    applicationConfig.GetInt(constant.UserIncentiveRedisPoolSize),
		MinIdleConn:                 applicationConfig.GetInt(constant.UserIncentiveRedisMinIdleConn),
		MaxConnAgeInHour:            applicationConfig.GetDuration(constant.UserIncentiveRedisMaxConnAgeInHour) * time.Hour,
		PoolTimeoutInMilliSecond:    applicationConfig.GetDuration(constant.UserIncentiveRedisPoolTimeoutInMilliSecond) * time.Millisecond,
		ReadTimeoutInMilliSecond:    applicationConfig.GetDuration(constant.UserIncentiveRedisReadTimeoutInMilliSecond) * time.Millisecond,
		WriteTimeoutInMilliSecond:   applicationConfig.GetDuration(constant.UserIncentiveRedisWriteTimeoutInMilliSecond) * time.Millisecond,
		ContextTimeoutInMilliSecond: applicationConfig.GetDuration(constant.UserIncentiveRedisContextTimeoutInMilliSecond) * time.Millisecond,
		Mode:                        applicationConfig.GetString(constant.UserIncentiveRedisMode),
	}

	cache.InitCache(redisConfig)
}

func closeCache() {
	err := cache.Close()
	if err != nil {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("error closing cache")
	}

	logger.FileLogger.Info().Msg("cache graceful shut down")
}

func initFCM() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("FCM: error loading application config file, env - " + flags.Env())
	}

	fcmConfig := flags.BaseConfigPath() + "/" + flags.Env() + "/" + applicationConfig.GetString(constant.FCMConfigFileName)

	fcm.InitFCM(fcmConfig)
}

func closeFCM() {
	err := fcm.Close()
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Stack().Err(err).Msg("error closing fcm")
	}
}

func startCron() {
	err := scheduler.StartCronJob()
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Stack().Msg("error getting cron")
	}

	logger.FileLogger.Info().Str("env", flags.Env()).Msg("Cron started...")
}

func startRouter() {

	router := api.GetRouter(middlewares.Logger(middlewares.LoggerMiddlewareOptions{}), middlewares.ContextCreator(), middlewares.CORSMiddleware(), middlewares.SecurityCheck())

	err := router.Run(fmt.Sprintf(":%d", flags.Port()))
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("error starting router")
	}
}
