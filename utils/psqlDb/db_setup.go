package psqldb

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	"time"

	_ "github.com/lib/pq"
)

// DB connection
//var dbConn *sql.DB

type PsqlOperation interface {
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	Begin() (*sql.Tx, error)
	PingContext(ctx context.Context) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DBImplementation struct {
	Operation PsqlOperation
}

//NewDatabaseSetup start a new db connection
func NewDatabaseSetup(conf *DbClientConfig) *DBImplementation {

	logger.FileLogger.Info().Interface("config", conf).Msg("PsqlDatabaseService - config")

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)

	var err error
	var dbConn PsqlOperation
	dbConn, err = sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("PsqlDatabaseService - open conn failed")
	}

	dbConn.SetMaxIdleConns(conf.MaxIdleConnections)
	dbConn.SetMaxOpenConns(conf.MaxOpenConnections)
	dbConn.SetConnMaxIdleTime(conf.MaxIdletime) //30
	dbConn.SetConnMaxLifetime(conf.MaxLifetime) //24 * time.Hour

	ctx := context.Background()

	err = dbConn.PingContext(ctx)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("PsqlDatabaseService - ping failed")
	}

	return &DBImplementation{Operation: dbConn}
}
