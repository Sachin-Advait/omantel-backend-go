package dblogger

import (
	"context"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

type Logging struct {
	UserID     string
	URL        string
	Type       string
	Msg        string
	Code       int
	Text       string
	CreateTime time.Time
}

var LogChan chan *Logging

func PrintDBLogs(log *Logging) {
	LogChan <- log
}

func InitDBLogger() {

	LogChan = make(chan *Logging, 100)

	logger.FileLogger.Info().Msg("LogDatabaseService - config")

	go insertLogs()
}

func insertLogs() {

	db := database.Get()

	for val := range LogChan {

		go func(val *Logging) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			_, err := db.ExecContext(ctx, "INSERT INTO logs_table(emp_id, url, type, msg, code, text, created) VALUES ($1,$2,$3,$4,$5,$6,$7);",
				val.UserID,
				val.URL,
				val.Type,
				val.Msg,
				val.Code,
				val.Text,
				val.CreateTime,
			)

			if err != nil {
				logger.FileLogger.Error().Str(logger.RequestIDLogParam, val.UserID).Interface("msg", val).Err(err).Msg("LogDatabaseService - insert failed")
			}
		}(val)
	}
}
