package logger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
)

type Logging struct {
	UserID     string
	URL        string
	Type       string
	Msg        string
	Text       string
	CreateTime time.Time
}

var LogChan chan *Logging

func InitDBLogger(conf *models.DbClientConfig) {

	LogChan = make(chan *Logging, 10)

	FileLogger.Info().Interface("config", conf).Msg("LogDatabaseService - config")

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)

	dbConn, err := sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		FileLogger.Fatal().Err(err).Msg("LogDatabaseService - open conn failed")
	}

	dbConn.SetMaxIdleConns(conf.MaxIdleConnections)
	dbConn.SetMaxOpenConns(conf.MaxOpenConnections)
	dbConn.SetConnMaxIdleTime(conf.MaxIdletime) //30
	dbConn.SetConnMaxLifetime(conf.MaxLifetime) //24 * time.Hour

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err = dbConn.PingContext(ctx)
	if err != nil {
		FileLogger.Fatal().Err(err).Msg("LogDatabaseService - ping failed")
	}

	FileLogger.Info().Msg("LogDatabaseService - DB ping successful")
	FileLogger.Info().Msg("LogDatabaseService - DB connection successful")

	go insertLogs(dbConn)
}

func insertLogs(dbConn *sql.DB) {

	for val := range LogChan {

		go func(val *Logging) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := dbConn.QueryContext(ctx, "INSERT INTO logs_table(emp_id, url, type, msg, text, created) VALUES ($1,$2,$3,$4,$5,$6);",
				val.UserID,
				val.URL,
				val.Type,
				val.Msg,
				val.Text,
				val.CreateTime,
			)

			if err != nil {
				FileLogger.Error().Interface("msg", val).Err(err).Msg("LogDatabaseService - insert failed")
			}
		}(val)
	}
}
