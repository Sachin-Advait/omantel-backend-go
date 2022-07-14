package main

/*
---RUN------
go build
./main --base-config-path ~/go/src/github.com/WeTutelage/login-service/resources
*/

import (
	"database/sql"
	"fmt"

	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/migration/migrate/helper"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

var migrateDB *sql.DB

func init() {
	utils.SetLoggerObj(commons.LOG_FILE_NAME)
}

func main() {

	//check if database is present or not.
	//If not then create a new db with the given dbName in config.
	//initDB()

	setupDB()

	defer migrateDB.Close()

	migrationpath := "../schemas/*"

	helper.StartMigration(migrateDB, migrationpath)
}

func initDB() {

	configs := &config.GetInstance().Database

	// logger.FileLogger.Info().
	// 	Str("host", configs.Host).
	// 	Int("port", configs.Port).
	// 	Str("dbName", configs.Name).
	// 	Str("user", configs.UserName).
	// 	Str("password", configs.Password).
	// 	Int("maxOpenConnections", configs.MaxOpenConnections).
	// 	Int("maxIdleConnections", configs.MaxIdleConnections).
	// 	Dur("connMaxLifetime", configs.MaxLifetime).
	// 	Dur("connMaxIdleTime", configs.MaxIdletime).
	// 	Msg("InitDB Psql Database Configuration")

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", configs.Host, configs.Port, configs.UserName, configs.Password)

	// Create connection pool
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		utils.LogObj.Fatal(err.Error(), "initDB - error psql open connection")
	}

	defer db.Close()

	rows, error := db.Query("SELECT distinct datname FROM pg_catalog.pg_database")
	if error != nil {
		utils.LogObj.Fatal(error.Error(), "initDB - error psql query")
	}

	defer rows.Close()

	for rows.Next() {
		var databaseName string

		err = rows.Scan(&databaseName)
		if err != nil {
			utils.LogObj.Fatal(err.Error(), "initDB - error psql scan")
		}

		if databaseName == configs.Name {
			return
		}
	}

	_, errCreate := db.Exec(fmt.Sprintf("CREATE  DATABASE \"%v\"", configs.Name))
	if errCreate != nil {
		utils.LogObj.Fatal(errCreate.Error(), "initDB - error psql execution")
	}
}

func setupDB() {

	//db config load
	dbConfig := &config.GetInstance().Database

	// logger.FileLogger.Info().
	// 	Str("host", dbConfig.Host).
	// 	Int("port", dbConfig.Port).
	// 	Str("dbName", dbConfig.Name).
	// 	Str("user", dbConfig.UserName).
	// 	Str("password", dbConfig.Password).
	// 	Int("maxOpenConnections", dbConfig.MaxOpenConnections).
	// 	Int("maxIdleConnections", dbConfig.MaxIdleConnections).
	// 	Dur("connMaxLifetime", dbConfig.MaxLifetime).
	// 	Dur("connMaxIdleTime", dbConfig.MaxIdletime).
	// 	Msg("SetupDB Psql Database Configuration")

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.Name)

	// Create connection pool
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		utils.LogObj.Fatal(err.Error(), "setupDB - error psql connection")
	}

	migrateDB = db
}
