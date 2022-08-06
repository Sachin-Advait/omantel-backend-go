package main

import (
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/api"
	"gitlab.com/rendezvousPoint/incentive_calculator/business"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	middlewares "gitlab.com/rendezvousPoint/incentive_calculator/utils/middlewares"
)

func init() {

	initEnvConfigs()
	initLogger()
	initFileLogger()
}

func main() {

	//cronService := crons.NewCronServices(&config.GetInstance().Database)
	//cronContext := crons.NewCronContext(cronService)
	//
	//go cronContext.StartCrons()

	incentiveProvider := buildIncentiveProvider()
	adapters.InitIncentive(incentiveProvider)

	adapters.InitIncentiveCron()

	startRouter()

	//cronService := crons.NewCronServices(&config.GetInstance().Database)
	//cronContext := crons.NewCronContext(cronService)
	//
	//go cronContext.StartCrons()
	//
	//serverContext := infra.NewServerContext(&config.GetInstance().Server)
	//
	//trxService := handler.NewSalesIncentiveService(&config.GetInstance().Database)
	//trxContext := handler.NewSalesIncentiveContext(trxService)
	//
	////err := trxService.UpdateEachPromoterData()
	////if err != nil {
	////	fmt.Println("promoter error:", err)
	////	return
	////}
	////
	////fmt.Println("UpdateEachPromoterData success")
	////
	////err = trxService.UpdateOtherManagersData()
	////if err != nil {
	////	fmt.Println("manager error:", err)
	////	return
	////}
	////
	////fmt.Println("UpdateOtherManagersData success")
	//
	//defer trxContext.SafeClose()
	//
	//serverContext.Mount(baseURL, trxContext.NewSalesIncentiveRouter())
	//serverContext.Mount(health, handler.NewHealthCheckRouter())
	//
	//// Start server
	//address := config.GetInstance().Server.GetAddressString()
	//serverContext.StartServer(address)
}

func initEnvConfigs() {
	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())
	logger.Info(nil).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Msg("initialized configs for env")
}

func initLogger() {

	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())

	loggerConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(nil).Stack().Err(err).Msg("error loading logger config file, env - " + flags.Env())
	}

	logger.InitLogger(logger.Level(loggerConfig.GetString(constant.LogLevelKey)))
}

func initFileLogger() {

	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())

	applicationConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(nil).Stack().Err(err).Msg("error loading application config file, env - " + flags.Env())
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
		logger.Fatal(nil).Err(err).Msg("unable to initialize file logger")
	}
}

func buildIncentiveProvider() models.IncentiveProvider {

	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("error loading application config file, env - " + flags.Env())
	}

	jwthelper.JwtSecretKey = applicationConfig.GetString(constant.JwtKey)
	if jwthelper.JwtSecretKey == "" {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("jwt secret key not found, env - " + flags.Env())
	}

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
		MaxIdletime:        time.Duration(applicationConfig.GetInt(constant.ConnMaxIdleTime)) * time.Hour,
	}

	incentiveConfig := &configs.IncentiveConfig{
		SessionExpiryTime:      time.Duration(applicationConfig.GetInt(constant.SessionExpiryTimeMinutes)) * time.Minute,
		RefreshTokenExpiryTime: time.Duration(applicationConfig.GetInt(constant.RefreshTokenExpiryHours)) * time.Hour,
	}

	return business.NewIncentiveService(incentiveConfig, dbConfig)
}

func startRouter() {

	router := api.GetRouter(middlewares.Logger(middlewares.LoggerMiddlewareOptions{}), middlewares.ContextCreator(), middlewares.CORSMiddleware())

	err := router.Run(fmt.Sprintf(":%d", flags.Port()))
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("error starting router")
	}
}
