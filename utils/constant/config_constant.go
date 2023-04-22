package constant

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

// application configuration keywords
const (
	JwtKey                   = "secretKey"
	AdminJwtKey              = "adminSecretKey"
	OTPEndpoint              = "otpUrl"
	CorsOriginEnabled        = "corsOrigin"
	SessionExpiryTimeMinutes = "sessionExpiryTimeMinutes"
	RefreshTokenExpiryHours  = "refreshTokenExpiryHours"
	CaptchaSitesEnableKey    = "captchaSitesEnableKey"
	CaptchaSecretKey         = "captchaSecretKey"
	VerifyCaptchaGoogleApi   = "verifyCaptchaGoogleApi"
	ShortCodeKey             = "shortCode"
	RequestorChannelIdKey    = "requestorChannelId"

	FCMConfigFileName       = "fcmConfigFileName"
	FCMOfferTopicKey        = "fcmOfferTopic"
	FCMTrainingTopicKey     = "fcmTrainingTopic"
	FCMNotificationTopicKey = "fcmNotificationTopic"
)

// config file names
const (
	LoggerConfigFile      = "logger"
	MoxyConfigFile        = "moxy"
	ApplicationConfigFile = "application"
)

// db configuration keywords
const (
	MasterCronTime = "masterCronTime"

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

	UserIncentiveRedisAddress                     = "userIncentiveRedis.address"
	UserIncentiveRedisPassword                    = "userIncentiveRedis.password"
	UserIncentiveRedisDb                          = "userIncentiveRedis.db"
	UserIncentiveRedisPoolSize                    = "userIncentiveRedis.poolSize"
	UserIncentiveRedisMinIdleConn                 = "userIncentiveRedis.minIdleConn"
	UserIncentiveRedisMaxConnAgeInHour            = "userIncentiveRedis.maxConnAgeInHour"
	UserIncentiveRedisPoolTimeoutInMilliSecond    = "userIncentiveRedis.poolTimeoutInMilliSecond"
	UserIncentiveRedisReadTimeoutInMilliSecond    = "userIncentiveRedis.readTimeoutInMilliSecond"
	UserIncentiveRedisWriteTimeoutInMilliSecond   = "userIncentiveRedis.writeTimeoutInMilliSecond"
	UserIncentiveRedisContextTimeoutInMilliSecond = "userIncentiveRedis.contextTimeoutInMilliSecond"
	UserIncentiveRedisKeyTTLInHour                = "userIncentiveRedis.keyTTLInHour"
	UserIncentiveRedisMode                        = "userIncentiveRedis.mode"
)

// log configuration keywords
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

// http configuration keywords
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
