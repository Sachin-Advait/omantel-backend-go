package models

import "time"

type DbClientConfig struct {
	Env      string
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

type HttpClientConfig struct {
	CertificatePath       string
	ConnectionTimeout     time.Duration
	KeepAliveDuration     time.Duration
	MaxIdleConnections    int
	IdleConnectionTimeout time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	Timeout               time.Duration
}
