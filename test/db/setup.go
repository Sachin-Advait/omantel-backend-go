package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DbClientConfig struct {
	Host     string
	Port     int
	DbName   string
	UserName string
	Password string
	//MaxOpenConnections <= 0 means unlimited
	MaxOpenConnections int
	//MaxIdleConnections zero means defaultMaxIdleConns; negative means 0
	MaxIdleConnections int
	//MaxLifetime maximum amount of time a connection may be reused
	MaxLifetime time.Duration
	//MaxIdletime maximum amount of time a connection may be idle before being closed
	MaxIdletime time.Duration
}

var db *sql.DB

func initDatabase() {
	conf := &DbClientConfig{
		Host:               "localhost",
		Port:               5432,
		DbName:             "postgres",
		UserName:           "postgres",
		Password:           "",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
		MaxLifetime:        2 * time.Hour,
		MaxIdletime:        30 * time.Minute,
	}

	psqlConf := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)

	conn, err := sql.Open("postgres", psqlConf+" connect_timeout=1")
	if err != nil {
		log.Fatal("PsqlDatabaseService - conn failed")
	}

	conn.SetMaxIdleConns(conf.MaxIdleConnections)
	conn.SetMaxOpenConns(conf.MaxOpenConnections)
	conn.SetConnMaxIdleTime(conf.MaxIdletime) //30
	conn.SetConnMaxLifetime(conf.MaxLifetime) //24 * time.Hour

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err = conn.PingContext(ctx)
	if err != nil {
		log.Fatal("PsqlDatabaseService - ping failed")
		return
	}

	db = conn
}
