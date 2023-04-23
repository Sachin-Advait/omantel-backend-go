package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"

	_ "github.com/lib/pq"
)

var dbClient *models.DbClientConfig

func init() {
	initEnvConfigs()
	initLogger()
	initFileLogger()
	initDB()
}

func main() {

	logger.FileLogger.Info().Msg("Application is up")

	if dbClient == nil {
		logger.FileLogger.Error().Err(errors.New("failed to load db config")).Msg("DbConfig load config failed")
		return
	}

	dbConn, err := NewDatabaseSetup(dbClient)
	if err != nil {
		logger.FileLogger.Error().Err(err).Msg("DbConfig failed connection")
		return
	}

	logger.FileLogger.Info().Interface("dbClient", dbConn).Msg("DB connection successfull")
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

func initDB() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("initHTTPClient: error loading application config file, env - " + flags.Env())
	}

	//db config load
	dbClient = &models.DbClientConfig{
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

}

// NewDatabaseSetup start a new db connection
func NewDatabaseSetup(conf *models.DbClientConfig) (*sql.DB, error) {

	logger.FileLogger.Info().Interface("config", conf).Msg("PsqlDatabaseService - config")

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)

	dbConn, err := sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		return nil, fmt.Errorf("db conn error - %v", err)
	}

	dbConn.SetMaxIdleConns(conf.MaxIdleConnections)
	dbConn.SetMaxOpenConns(conf.MaxOpenConnections)
	dbConn.SetConnMaxIdleTime(conf.MaxIdletime) //30
	dbConn.SetConnMaxLifetime(conf.MaxLifetime) //24 * time.Hour

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err = dbConn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping error - %v", err)
	}

	return dbConn, nil
}
