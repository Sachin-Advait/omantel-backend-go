package constant

const (
	AppName    = "feed-service"
	AppVersion = "v1.0.0"
)

const (
	Authorization      = "Authorization"
	SuperAuthorization = "Super-Authorization"
	XSourceID          = "X-SourceID"
	XClientLocalIP     = "X-ClientLocalIP"
	XClientPublicIP    = "X-ClientPublicIP"
	XMACAddress        = "X-MACAddress"
	XSystemInfo        = "X-SystemInfo"
	XLocation          = "X-Location"
	XAppID             = "X-AppID"
	XProductVersion    = "X-ProductVersion"
	XOperatingSystem   = "X-OperatingSystem"
	XPlatform          = "X-Platform"
	XAppVersion        = "X-App-version"
	XOSVersion         = "X-OS-version"
	XDeviceID          = "X-DeviceID"
	XDeviceType        = "X-DeviceType"
	XRequestID         = "X-Request-Id"
	XCorrelationID     = "X-Correlation-Id"
	XUserId            = "X-User-Id"
	XUserType          = "X-UserType"
	XAPIVersion        = "X-Api-Version"
)

// Request Params
const (
	Region  = "region"
	Outlet  = "outlet"
	Product = "product"
	Plan    = "plan"
	Target  = "target"
	All     = "all"
)

//db configuration keywords
const (
	PsqlDatabase    = "psqlDB"
	PsqlHost        = PsqlDatabase + "." + "host"
	PsqlPort        = PsqlDatabase + "." + "port"
	PsqlDbName      = PsqlDatabase + "." + "dbName"
	PsqlUser        = PsqlDatabase + "." + "user"
	PsqlPassword    = PsqlDatabase + "." + "password"
	PsqlMaxOpenConn = PsqlDatabase + "." + "maxOpenConnections"
	PsqlMaxIdleConn = PsqlDatabase + "." + "maxIdleConnections"
	ConnMaxLifetime = PsqlDatabase + "." + "connMaxLifetimeInHour"
	ConnMaxIdleTime = PsqlDatabase + "." + "connMaxIdleTimeInMinute"

	LogPsqlDatabase    = "logPsqlDB"
	LogPsqlHost        = LogPsqlDatabase + "." + "host"
	LogPsqlPort        = LogPsqlDatabase + "." + "port"
	LogPsqlDbName      = LogPsqlDatabase + "." + "dbName"
	LogPsqlUser        = LogPsqlDatabase + "." + "user"
	LogPsqlPassword    = LogPsqlDatabase + "." + "password"
	LogPsqlMaxOpenConn = LogPsqlDatabase + "." + "maxOpenConnections"
	LogPsqlMaxIdleConn = LogPsqlDatabase + "." + "maxIdleConnections"
	LogConnMaxLifetime = LogPsqlDatabase + "." + "connMaxLifetimeInHour"
	LogConnMaxIdleTime = LogPsqlDatabase + "." + "connMaxIdleTimeInMinute"
)

//log configuration keywords
const (
	LogLevelKey              = "level"
	ConsoleLoggingEnabledKey = "consoleLoggingEnabled"
	FileLoggingEnabledKey    = "fileLoggingEnabled"
	DebugLoggingEnabledKey   = "debugLoggingEnabled"
	EncodeLogAsJsonKey       = "encodeLogsAsJson"
	DirectoryKey             = "directory"
	FileNameKey              = "fileName"
	MaxSizeInMBKey           = "maxSizeInMB"
	MaxBackUpsKey            = "maxBackups"
	MaxAgeInDaysKey          = "maxAgeInDays"
	Compress                 = "compress"
)

//http configuration keywords
const (
	HTTPTLSCertificateFileLocation       = "http.tlsCertificate"
	HTTPConnectTimeoutInMillisKey        = "http.connectTimeoutInMillis"
	HTTPKeepAliveDurationInMillisKey     = "http.keepAliveDurationInMillis"
	HTTPMaxIdleConnectionsKey            = "http.maxIdleConnections"
	HTTPIdleConnectionTimeoutInMillisKey = "http.idleConnectionTimeoutInMillis"
	HTTPTlsHandshakeTimeoutInMillisKey   = "http.tlsHandshakeTimeoutInMillis"
	HTTPExpectContinueTimeoutInMillisKey = "http.expectContinueTimeoutInMillis"
	HTTPTimeoutInMillisKey               = "http.timeoutInMillis"
)

const (
	HeaderInContext = "Context-Header"
)

//jwt tokens
const (
	TokenExpired = -1
	TokenInvalid = 0
	TokenValid   = 1

	AccessTokenKey  = "ACCESS-TOKEN"
	RefreshTokenKey = "REFRESH-TOKEN"
	IsAdminTokenKey = "IS-ADMIN"
)

// Log Levels
const (
	TraceLevel = "trace"
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

//feed configuration keywords
const (
	JwtKey                   = "secretKey"
	AdminJwtKey              = "adminSecretKey"
	OTPEndpoint              = "otpUrl"
	SessionExpiryTimeMinutes = "sessionExpiryTimeMinutes"
	RefreshTokenExpiryHours  = "refreshTokenExpiryHours"
	CaptchaSitesEnableKey    = "captchaSitesEnableKey"
	CaptchaSecretKey         = "captchaSecretKey"
	VerifyCaptchaGoogleApi   = "verifyCaptchaGoogleApi"
	ShortCodeKey             = "shortCode"
	RequestorChannelIdKey    = "requestorChannelId"
)

//config file names
const (
	LoggerConfigFile      = "logger"
	MoxyConfigFile        = "moxy"
	ApplicationConfigFile = "application"
)

//routes
const (
	SwaggerRoute           = "/swagger/*any"
	ActuatorRoute          = "/actuator/*any"
	SalesIncentiveRoute    = "/sales/incentive"
	SalesIncentiveDocRoute = "/sales/incentive/v1"

	SalesIncentiveV2Route    = "/sales/incentive/v2"
	SalesIncentiveV2DocRoute = "/sales/incentive/v2"

	ClientPath                 = "/client"
	TargetChartPath            = "/target-chart"
	SalePerPlanPath            = "/sale-per-plan"
	GenerateOTPPath            = "/generate-otp"
	VerifyOTPPath              = "/verify-otp"
	UserReportPath             = "/user-report"
	ReportToOutletManagerPath  = "/report_to_outlet_manager"
	ReportToManagerPath        = "/report_to_manager"
	ReportToSeniorManagerPath  = "/report_to_senior_manager"
	ReportToGeneralManagerPath = "/report_to_general_manager"
	ReportToVicePresidentPath  = "/report_to_vice_president"
	UserPayoutPath             = "/usr_payout"
	UserTargetPath             = "/usr_target"
	UsrSalesPerPlanPath        = "/usr_sales_per_plan"
	OutletSalesCountPath       = "/outlet_sales_count"
	OutletSalesTargetPath      = "/outlet_sales_target"
	RegionSalesCountPath       = "/region_sales_count"
	RegionSalesTargetPath      = "/region_sales_target"
	TopTenSalesCountPath       = "/top_ten_sales_count"
	TopTenSalesTargetPath      = "/top_ten_sales_target"
	RenewTokenPath             = "/renew_token"

	CaptchaVerificationPath = "/captcha_verify"

	FetchRegionPath   = "/fetch_region"
	UpsertRegionPath  = "/upsert_region"
	DeleteRegionPath  = "/delete_region"
	FetchOutletPath   = "/fetch_outlet"
	UpsertOutletPath  = "/upsert_outlet"
	DeleteOutletPath  = "/delete_outlet"
	FetchProductPath  = "/fetch_product"
	UpsertProductPath = "/upsert_product"
	DeleteProductPath = "/delete_product"
	FetchPlanPath     = "/fetch_plan"
	UpsertPlanPath    = "/upsert_plan"
	DeletePlanPath    = "/delete_plan"
	FetchTargetPath   = "/fetch_target"
	UpsertTargetPath  = "/upsert_target"
	DeleteTargetPath  = "/delete_target"
)

var (
	Table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

const (
	HBBCategory     = "HBB"
	BaqatiCategory  = "BAQATI"
	HayyakCategory  = "HAYYAK"
	UpgradeCategory = "UPGRADE"
)

const (
	VicePresident  = "VICE PRESIDENT"
	GeneralManager = "GM"
	SeniorManager  = "SM"
	Manager        = "MANAGER"
	TeamLeader     = "TL"
	OutletManager  = "OUTLET MANAGER"
	Promoter       = "PROMOTER"
	CSR            = "CSR"
)

var (
	ManagerMap = map[string]string{
		OutletManager:  "outlet_manager",
		TeamLeader:     "team_lead",
		Manager:        "manager",
		SeniorManager:  "senior_manager",
		GeneralManager: "general_manager",
		VicePresident:  "vice_president",
	}
)

var (
	EmployeePositionMap = map[string]string{
		Promoter:       "A",
		CSR:            "A",
		OutletManager:  "OM",
		TeamLeader:     "TL",
		Manager:        "M",
		SeniorManager:  "SM",
		GeneralManager: "GM",
		VicePresident:  "VP",
	}
)

var ProductGroupMap = map[string]struct{}{
	HBBCategory:     struct{}{},
	BaqatiCategory:  struct{}{},
	HayyakCategory:  struct{}{},
	UpgradeCategory: struct{}{},
}

const (
	MaxAllowedLoginAttempt    = 5
	OtpGenerateBlockTimeInMin = 15
	OtpRegenerateAllowedCount = 5
	OtpTimeoutInMin           = 15
)

var GetTemplate = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
   <soapenv:Header/>
   <soapenv:Body>
      <sen:SendNotificationRq>
         <sen:RequestHeader>
            <FunctionId>17002004</FunctionId>
            <TransactionId>0</TransactionId>
            <RequestId>{{.HeaderRequestID}}</RequestId>
            <RequestorChannelId>{{.HeaderRequestID}}</RequestorChannelId>
            <RequestorUserId>{{.HeaderStaffID}}</RequestorUserId>
            <Timestamp>{{.HeaderTimeStamp}}</Timestamp>
         </sen:RequestHeader>
         <sen:RequestBody>
            <sen:Notification>
               <NotificationType>SMS</NotificationType>
               <NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
               <NotificationLanguage>1</NotificationLanguage>
               <ShortCode>Omantel</ShortCode>
               <IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
               <NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
               <Body>{{.Message}}</Body>
               <!--0 for English text 
				8 for other language text
				4 for binary messages like WAP
				16 for flash sms
				and It will be commmon for all SMSs-->
               <BodyCoding>0</BodyCoding>
               <!--0 for English text
						1 for other language text
						2 for binary messages like WAP
						It will be commmon for all SMSs-->
               <BodyContentType>0</BodyContentType>
               <ReceiverList>
                  <!--One receiver only 8 digit Omantel MSISDN without country code-->
                  <Receiver>{{.ContactNumber}}</Receiver>
               </ReceiverList>
            </sen:Notification>
         </sen:RequestBody>
      </sen:SendNotificationRq>
   </soapenv:Body>
</soapenv:Envelope>
`
