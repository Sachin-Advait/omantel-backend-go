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

// Log Params
const (
	RequestIDLogParam   = "requestId"
	UserIDLogParam      = "userId"
	StatusCodeLogParam  = "statusCode"
	EndTimeLogParam     = "endTime"
	StartTimeLogParam   = "startTime"
	PathLogParam        = "path"
	QueryLogParam       = "query"
	HeaderLogParams     = "header"
	LatencyLogParam     = "latency"
	ClientIPLogParam    = "clientIP"
	MethodLogParam      = "method"
	BodyLogParam        = "body"
	ErrorLogParam       = "error"
	ResponseParam       = "response"
	CorrelationLogParam = "correlationId"
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
	JwtKey = "secretKey"
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
)

const (
	FeedsKey              = "Feed"
	CachePipeSplitterChar = "|"
	CacheSplitterChar     = "~"
)
