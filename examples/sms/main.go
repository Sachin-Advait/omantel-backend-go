package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

// req := Request{
// 	UserId:                      userId,
// 	FunctionId:                  "17002004",
// 	TransactionId:               "0",
// 	RequestId:                   strconv.Itoa(int(time.Now().Unix())),
// 	RequestorChannelId:          "OMSIP",
// 	RequestorUserId:             "OMSIP",
// 	Timestamp:                   time.Now().AddDate(2, 0, 0).Format("2006-01-02 15:04:05"),
// 	NotificationType:            "SMS",
// 	NotificationDateTime:        time.Now().Format("2006-01-02 15:04:05"),
// 	NotificationLanguage:        "1",
// 	ShortCode:                   "OMSIP",
// 	IsDeliveryStatusRequired:    "1",
// 	NotificaitionValidityPeriod: time.Now().Format("2006-01-02 15:04:05"),
// 	Body:                        fmt.Sprintf("Dear %v, Your OTP is %v", userId, otp),
// 	BodyCoding:                  "8",
// 	BodyContentType:             "1",
// 	Receiver:                    mobile,
// 	OTP:                         strconv.Itoa(otp),
// }

var SMSRequestTemplate = []byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
<soapenv:Header/>
<soapenv:Body>
   <sen:SendNotificationRq>
	  <sen:RequestHeader>
		 <FunctionId>{{.FunctionId}}</FunctionId>
		 <TransactionId>{{.TransactionId}}</TransactionId>
		 <RequestId>OMSIP_{{.RequestId}}</RequestId>
		 <RequestorChannelId>ZONE</RequestorChannelId>
		 <RequestorUserId>OMSIP</RequestorUserId>
		 <Timestamp>{{.Timestamp}}</Timestamp>
	  </sen:RequestHeader>
	  <sen:RequestBody>
		 <sen:Notification>
			<NotificationType>SMS</NotificationType>
			<NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
			<NotificationLanguage>1</NotificationLanguage>
			<ShortCode>OMSIP</ShortCode>
			<IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
			<NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
			<Body>Dear {{.UserId}}, Your OTP is {{.OTP}}</Body>
			<BodyCoding>8</BodyCoding>
			<BodyContentType>1</BodyContentType>
			<ReceiverList>
			   <Receiver>{{.Receiver}}</Receiver>
			</ReceiverList>
		 </sen:Notification>
	  </sen:RequestBody>
   </sen:SendNotificationRq>
</soapenv:Body>
</soapenv:Envelope>`)

// var SMSRequestTemplate = []byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
// <soapenv:Header/>
// <soapenv:Body>
//    <sen:SendNotificationRq>
// 	  <sen:RequestHeader>
// 		 <FunctionId>{{.FunctionId}}</FunctionId>
// 		 <TransactionId>0</TransactionId>
// 		 <RequestId>OMSIP_{{.RequestId}}</RequestId>
// 		 <RequestorChannelId>OMSIP</RequestorChannelId>
// 		 <RequestorUserId>OMSIP</RequestorUserId>
// 		 <Timestamp>{{.Timestamp}}</Timestamp>
// 	  </sen:RequestHeader>
// 	  <sen:RequestBody>
// 		 <sen:Notification>
// 			<NotificationType>SMS</NotificationType>
// 			<NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
// 			<NotificationLanguage>1</NotificationLanguage>
// 			<ShortCode>ROP Permit</ShortCode>
// 			<IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
// 			<NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
// 			<Body>Dear {{.UserId}}, Your OTP is {{.OTP}}

// 		</Body>
// 			<BodyCoding>8</BodyCoding>
// 			<BodyContentType>1</BodyContentType>
// 			<ReceiverList>
// 			   <Receiver>{{.Receiver}}</Receiver>
// 			</ReceiverList>
// 		 </sen:Notification>
// 	  </sen:RequestBody>
//    </sen:SendNotificationRq>
// </soapenv:Body>
// </soapenv:Envelope>`)

// var response = []byte(`<?xml version="1.0" encoding="UTF-8"?>
// <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
//    <soapenv:Body>
//       <NS1:SendNotificationRs xmlns:NS1="http://www.omantel.om/esb/messages/sendnotification">
//          <NS1:ResponseHeader>
//             <FunctionId>17002004</FunctionId>
//             <TransactionId>36896</TransactionId>
//             <RequestId>OMSIP_1669540763</RequestId>
//             <ReturnCode>0</ReturnCode>
//             <ReturnDesc>Success/ Successfully validated the request</ReturnDesc>
//          </NS1:ResponseHeader>
//       </NS1:SendNotificationRs>
//    </soapenv:Body>
// </soapenv:Envelope>`)

var client *http.Client
var otpUrl string
var fileName string

type Request struct {
	//Values are set in below fields as per the request
	UserId                      string
	FunctionId                  string
	TransactionId               string
	RequestId                   string
	RequestorChannelId          string
	RequestorUserId             string
	Timestamp                   string
	NotificationType            string
	NotificationDateTime        string
	NotificationLanguage        string
	ShortCode                   string
	IsDeliveryStatusRequired    string
	NotificaitionValidityPeriod string
	Body                        string
	BodyCoding                  string
	BodyContentType             string
	Receiver                    string
	OTP                         string
}

type Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soapenv string   `xml:"soapenv,attr"`
	Body    struct {
		Text               string `xml:",chardata"`
		SendNotificationRs struct {
			Text           string `xml:",chardata"`
			NS1            string `xml:"NS1,attr"`
			ResponseHeader struct {
				Text          string `xml:",chardata"`
				FunctionId    string `xml:"FunctionId"`
				TransactionId string `xml:"TransactionId"`
				RequestId     string `xml:"RequestId"`
				ReturnCode    string `xml:"ReturnCode"`
				ReturnDesc    string `xml:"ReturnDesc"`
			} `xml:"ResponseHeader"`
		} `xml:"SendNotificationRs"`
	} `xml:"Body"`
}

func init() {

	initEnvConfigs()
	initLogger()
	initFileLogger()
	initHTTPClient()
}

func main() {

	fmt.Println("application is running: config loaded successfully - " + time.Now().String())
	fmt.Println("otp url: " + otpUrl)
	fmt.Println("certificate fileName: " + fileName)

	var Mobile string

	fmt.Printf("Enter Your Mobile Number: ")
	fmt.Scanln(&Mobile)
	fmt.Println("Your Mobile Number is: " + Mobile)

	t := time.Now()
	fmt.Println("Initiating SOAP call...")
	if err := callSOAPClientSteps("TESTUSER", Mobile, 123456); err != nil {
		fmt.Println("SOAP call failed: "+err.Error(), "timeOut: "+time.Since(t).String())
		return
	}

	fmt.Println("SOAP call successfull")
}

func initEnvConfigs() {
	configs.Init(flags.BaseConfigPath() + "/" + flags.Env())
	logger.Info(context.TODO()).Str("env", flags.BaseConfigPath()+"/"+flags.Env()).Msg("initialized configs for env")
}

func initLogger() {

	loggerConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Stack().Err(err).Msg("error loading logger config file, env - " + flags.Env())
	}

	logger.InitLogger(logger.Level(loggerConfig.GetString(constant.LogLevelKey)))
}

func initFileLogger() {

	applicationConfig, err := configs.Get(constant.LoggerConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Stack().Err(err).Msg("error loading application config file, env - " + flags.Env())
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

	config := logger.FileConfig{
		EncodeLogsAsJson:      applicationConfig.GetBool(constant.EncodeLogAsJsonKey),
		ConsoleLoggingEnabled: applicationConfig.GetBool(constant.ConsoleLoggingEnabledKey),
		FileLoggingEnabled:    applicationConfig.GetBool(constant.FileLoggingEnabledKey),
		Directory:             applicationConfig.GetString(constant.DirectoryKey),
		Filename:              applicationConfig.GetString(constant.FileNameKey),
		MaxSize:               applicationConfig.GetInt(constant.MaxSizeInMBKey),
		MaxBackups:            applicationConfig.GetInt(constant.MaxBackUpsKey),
		MaxAge:                applicationConfig.GetInt(constant.MaxAgeInDaysKey),
		Compress:              applicationConfig.GetBool(constant.Compress),
	}

	if err != nil {
		logger.Fatal(context.TODO()).Err(err).Msg("unable to initialize file logger")
	}

	logger.Info(context.TODO()).Interface("fileConfig", config).Msg("fileConfig: load configuration")
}

func initHTTPClient() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Stack().Err(err).Msg("initHTTPClient: error loading application config file, env - " + flags.Env())
	}

	certificate := applicationConfig.GetString(constant.HTTPTLSCertificateFileLocation)
	connectTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPConnectTimeoutInMillisKey)
	keepAliveDuration := time.Millisecond * applicationConfig.GetDuration(constant.HTTPKeepAliveDurationInMillisKey)
	maxIdleConnections := applicationConfig.GetInt(constant.HTTPMaxIdleConnectionsKey)
	idleConnectionTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPIdleConnectionTimeoutInMillisKey)
	tLSHandshakeTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPTlsHandshakeTimeoutInMillisKey)
	expectContinueTimeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPExpectContinueTimeoutInMillisKey)
	timeout := time.Millisecond * applicationConfig.GetDuration(constant.HTTPTimeoutInMillisKey)

	otpUrl = applicationConfig.GetString(constant.OTPEndpoint)
	fileName = certificate

	dialer := &net.Dialer{
		Timeout:   connectTimeout,
		KeepAlive: keepAliveDuration,
	}

	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client = &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
			TLSClientConfig: &tls.Config{
				//RootCAs:            caCertPool,
				InsecureSkipVerify: true,
			},
			MaxIdleConns:          maxIdleConnections,
			IdleConnTimeout:       idleConnectionTimeout,
			TLSHandshakeTimeout:   tLSHandshakeTimeout,
			ExpectContinueTimeout: expectContinueTimeout,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
		Timeout: timeout,
	}

	logger.Info(context.TODO()).Dur("timeout", timeout).Msg("httpConfig: load configuration")
}

func callSOAPClientSteps(userID, mobile string, otp int) error {

	req := populateRequest(userID, mobile, otp)

	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		return err
	}

	_, err = soapCall(httpReq)
	if err != nil {
		return err
	}

	return nil
}

func populateRequest(userId, mobile string, otp int) *Request {

	req := Request{
		UserId:                      userId,
		FunctionId:                  "17002004",
		TransactionId:               generateOTP(5),
		RequestId:                   strconv.Itoa(int(time.Now().Unix())),
		RequestorChannelId:          "OMSIP",
		RequestorUserId:             "OMSIP",
		Timestamp:                   time.Now().Format("2006-01-02T15:04:05"),
		NotificationType:            "SMS",
		NotificationDateTime:        time.Now().Format("2006-01-02T15:04:05"),
		NotificationLanguage:        "1",
		ShortCode:                   "OMSIP",
		IsDeliveryStatusRequired:    "1",
		NotificaitionValidityPeriod: time.Now().Format("2006-01-02T15:04:05"),
		Body:                        fmt.Sprintf("Dear %v, Your OTP is %v", userId, otp),
		BodyCoding:                  "8",
		BodyContentType:             "1",
		Receiver:                    mobile,
		OTP:                         strconv.Itoa(otp),
	}

	return &req
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("smsRequest").Parse(string(SMSRequestTemplate))
	if err != nil {
		return nil, errors.New("parsing-" + err.Error())
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, req)
	if err != nil {
		return nil, errors.New("execute-" + err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return nil, errors.New("encode-" + err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, otpUrl, bytes.NewBuffer(doc.Bytes()))
	if err != nil {
		return nil, errors.New("request-" + err.Error())
	}

	r.Header.Set("user-agent", "sampleTest")
	r.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("SA-OmantelSIP:Des!B@buEnglishMem@7533")))
	r.Header.Set("SOAPAction", "SendNotification")

	fmt.Printf("requestID: %v, transactionID: %v\n", req.RequestId, req.TransactionId)
	fmt.Println("request type", http.MethodPost)
	fmt.Println("request url", otpUrl)
	fmt.Println("request headers", r.Header)
	fmt.Println("request body", doc.String())

	return r, nil
}

func soapCall(req *http.Request) (*Response, error) {

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("call-" + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("readall-" + err.Error())
	}
	defer resp.Body.Close()

	fmt.Printf("response raw body: %+v\n", string(body))

	var response Response

	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("unmarshall-" + err.Error())
	}

	fmt.Printf("response code: %+v\n", resp.StatusCode)
	fmt.Printf("response body: %+v\n", response)

	if response.Body.SendNotificationRs.ResponseHeader.ReturnCode != "0" {
		return nil, errors.New(response.Body.SendNotificationRs.ResponseHeader.ReturnDesc)
	}

	return &response, nil
}

func generateOTP(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
