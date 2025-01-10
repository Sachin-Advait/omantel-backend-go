package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"

	_ "github.com/lib/pq"
)

type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// dbClient is the default client
type dbClient struct {
	*sql.DB
}

var db *dbClient

// InitDatabase is used to initialise the database
func InitDatabase(conf *models.DbClientConfig) {

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require", conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)

	conn, err := sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("PsqlDatabaseService - conn failed")
	}

	conn.SetMaxIdleConns(conf.MaxIdleConnections)
	conn.SetMaxOpenConns(conf.MaxOpenConnections)
	conn.SetConnMaxIdleTime(conf.MaxIdletime) //30
	conn.SetConnMaxLifetime(conf.MaxLifetime) //24 * time.Hour

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err = conn.PingContext(ctx)
	if err != nil {
		logger.Error(context.TODO()).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Err(err).Msg("PsqlDatabaseService - ping failed")
		logger.FileLogger.Fatal().Err(err).Msg("PsqlDatabaseService - ping failed")
	}

	// now create the database
	db = &dbClient{DB: conn}
	logger.Info(context.TODO()).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Msg("PsqlDatabaseService - ping successfull")
}

// Get is used to get the database instance client
func Get() Database {
	return db
}

// Close is used to close the database instance
func Close() error {
	return db.Close()
}
