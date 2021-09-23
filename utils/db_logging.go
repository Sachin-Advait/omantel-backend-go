package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/rs/xid"
	"github.com/spf13/viper"
)

// Logs structure
type Logs struct {
	LogType                 string
	LogMessage              string
	StartTime               time.Time
	RequestBody             string
	ResponseBody            string
	ProcessingTime          int64
	ResponseStatus          int
	Url                     string
	HttpMethod              string
	RemoteAddress           string
	RequestHeaderAndValues  string
	ResponseHeaderAndValues string
	RequestID               string
	RequestReceivedDate     time.Time
	ServerIP                string
}

var (
	pgConn *sql.DB
)

func getRefreshTime() time.Time {
	indiaTime, _ := time.LoadLocation("Asia/Calcutta")
	now := time.Now().In(indiaTime)
	year, month, day := now.Date()
	return time.Date(year, month, day, 07, 00, 0, 0, now.Location())
}

//Get Server IP
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//PostgresConnection Connection Open
func postgresConnect() {

	user := viper.GetString("database.username")
	password := viper.GetString("database.password")
	db := viper.GetString("database.name")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	dbMaxOpenConn := viper.GetInt("database.maxOpenConn")
	dbMaxIdleConn := viper.GetInt("database.maxIdelConn")

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=Asia/Kolkata",
		host, port, user,
		password, db)

	var err error

	// Create connection pool
	pgConn, err = sql.Open("postgres", dbinfo+" connect_timeout=1")
	if err != nil {
		LogObj.Error("failed to connect db for logging query, err: " + err.Error())
		return
	}

	pgConn.SetMaxOpenConns(dbMaxOpenConn)
	pgConn.SetMaxIdleConns(dbMaxIdleConn)
	pgConn.SetConnMaxLifetime(24 * time.Hour)

	ctx := context.Background()

	err = pgConn.PingContext(ctx)
	if err != nil {
		LogObj.Error("failed to ping db for logging query, err: " + err.Error())
		return
	}

	return
}

func LogHelper(r *http.Request) Logs {
	logs := Logs{}
	logs.StartTime = time.Now()
	logs.RequestHeaderAndValues = fmt.Sprint(r.Header)
	logs.HttpMethod = r.Method
	logs.RemoteAddress = r.Header.Get("X-ClientIP")
	logs.LogMessage = xid.New().String()
	logs.RequestID = logs.LogMessage
	logs.RequestReceivedDate = time.Now()
	logs.ServerIP = getLocalIP()
	logs.Url = r.URL.Path
	return logs
}

func structToQueryParams(params interface{}) []interface{} {

	var queryParams []interface{}

	values := reflect.ValueOf(params).Elem()
	for i := 0; i < values.NumField(); i++ {
		queryParams = append(queryParams, values.Field(i).Interface())
	}

	return queryParams
}

func insertRequestResponse(output Logs) {
	ctx := context.Background()

	query := "insert into tbl_amx_web_logging(requestbody,responsebody,processingtimeinms,responsestatus,url,httpmethod,remoteaddress,requestheaderandvalues,responseheaderandvalues,requestid,requestreceiveddate,server_ip) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)"
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	_, err := pgConn.ExecContext(ctx, query, structToQueryParams(&output)[3:]...)
	if err != nil {
		LogObj.Error("failed to insert trx logs, err: " + err.Error())
	}

}

func WriteInDB(output Logs) {
	// if output.ProcessingTime > 1300 {
	// 	go email.SendFailureEmail("Pledge API :"+output.Url, " API response Took more than a second", "Request Params :"+output.RequestBody+" time taken :"+strconv.FormatInt(output.ProcessingTime, 10)+"ms")
	// }

	if pgConn == nil {
		postgresConnect()
	}

	if pgConn != nil {
		insertRequestResponse(output)
	}
}
