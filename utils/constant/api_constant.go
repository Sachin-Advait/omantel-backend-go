package constant

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

// routes
const (
	ApiRoute                    = "/api"
	SwaggerRoute                = "/swagger/*any"
	ActuatorRoute               = "/actuator/*any"
	SalesIncentiveRoute         = "/sales/incentive"
	SalesIncentiveDocRoute      = "/sales/incentive/v1"
	SalesIncentiveAdminRoute    = "/sales/incentive/admin"
	SalesIncentiveDocAdminRoute = "/sales/incentive/admin/v1"

	GenerateTokenPath = "/token"

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
	UserSalesCountPath         = "/usr_sales_count"
	UserSalesPerPlanPath       = "/usr_sales_per_plan"
	UserPlanPerCategoryPath    = "/user_plan_category"
	UserMNPSalesCountPath      = "/usr_mnp_sale_count"
	UserMNPPlanCountPath       = "/usr_mnp_plan_count"
	OutletSalesCountPath       = "/outlet_sales_count"
	OutletSalesTargetPath      = "/outlet_sales_target"
	RegionSalesCountPath       = "/region_sales_count"
	RegionSalesTargetPath      = "/region_sales_target"
	TopTenSalesCountPath       = "/top_ten_sales_count"
	TopTenSalesTargetPath      = "/top_ten_sales_target"
	RenewTokenPath             = "/renew_token"
	GetUserPath                = "/get_user"
	UserLogoutPath             = "/logout"

	CaptchaVerificationPath = "/captcha_verify"

	FetchRegionPath  = "/fetch_region"
	UpsertRegionPath = "/upsert_region"
	DeleteRegionPath = "/delete_region"

	FetchOutletPath  = "/fetch_outlet"
	UpsertOutletPath = "/upsert_outlet"
	DeleteOutletPath = "/delete_outlet"

	FetchProductPath  = "/fetch_product"
	UpsertProductPath = "/upsert_product"
	DeleteProductPath = "/delete_product"

	FetchPlanPath  = "/fetch_plan"
	UpsertPlanPath = "/upsert_plan"
	DeletePlanPath = "/delete_plan"

	FetchTargetPath  = "/fetch_target"
	UpsertTargetPath = "/upsert_target"
	DeleteTargetPath = "/delete_target"

	FetchVOCCriteriaPath  = "/fetch_voc_criteria"
	UpsertVOCCriteriaPath = "/upsert_voc_criteria"
	DeleteVOCCriteriaPath = "/delete_voc_criteria"

	FetchVOCMaxPayPath  = "/fetch_voc_pay"
	UpsertVOCMaxPayPath = "/upsert_voc_pay"
	DeleteVOCMaxPayPath = "/delete_voc_pay"

	FetchVOCUserPath  = "/fetch_voc_user"
	UpsertVOCUserPath = "/upsert_voc_user"
	DeleteVOCUserPath = "/delete_voc_user"

	UpsertNotificationPath = "/upsert_notification"
	FetchNotificationPath  = "/fetch_notification"
	DeleteNotificationPath = "/delete_notification"

	LoginChartPath   = "/login_chart"
	LoginDetailsPath = "/login_details"
)
