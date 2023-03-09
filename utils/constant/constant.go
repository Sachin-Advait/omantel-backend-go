package constant

const (
	AppName    = "feed-service"
	AppVersion = "v1.0.0"
)

const (
	HeaderInContext = "Context-Header"
)

// Redis Mode and Keys
const (
	Cluster    = "CLUSTER"
	NonCluster = "NON_CLUSTER"

	RedisUserReportKey         = "%s|%v|%v|%v|Report"    // userId|month|quarter|year|Report
	RedisOutletRegionReportKey = "%s|%s|%v|%v|%v|Report" // outlet|region|month|quarter|year|Report
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

var (
	Table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

const (
	HBBCategory     = "HBB"
	BaqatiCategory  = "BAQATI"
	HayyakCategory  = "HAYYAK"
	UpgradeCategory = "UPGRADE"
	VOCCategory     = "VOC"
)

const (
	MaxAllowedLoginAttempt    = 5
	OtpGenerateBlockTimeInMin = 15
	OtpRegenerateAllowedCount = 5
	OtpTimeoutInMin           = 15
)

const (
	VicePresident  = "vice_president"
	GeneralManager = "general_manager"
	SeniorManager  = "senior_manager"
	Manager        = "manager"
	//TeamLeader     = "team_lead"
	OutletManager = "outlet_manager"
	Promoter      = "promoter"
	CSR           = "csr"

	VicePresidentCC  = "VICE PRESIDENT"
	GeneralManagerCC = "GENERAL MANAGER"
	SeniorManagerCC  = "SENIOR MANAGER"
	ManagerCC        = "MANAGER"
	//TeamLeaderCC     = "TEAM LEAD"
	OutletManagerCC = "OUTLET MANAGER"
	PromoterCC      = "PROMOTER"
	CSRCC           = "CSR"

	VicePresidentAbb  = "VP"
	GeneralManagerAbb = "GM"
	SeniorManagerAbb  = "SM"
	ManagerAbb        = "M"
	//TeamLeaderAbb     = "TL"
	OutletManagerAbb = "OM"
	PromoterAbb      = "A"
	CSRAbb           = "A"
)

const (
	VOCUnsatisfactory = "unsatisfactory"
	VOCAlmostAchieved = "almost_achieved"
	VOCAchieved       = "achieved"
	VOCWellAchieved   = "well_achieved"
	VOCOutstanding    = "outstanding"
)

var (
	ManagerMap = map[string]string{
		OutletManager:    "OM",
		OutletManagerAbb: "OM",
		OutletManagerCC:  "OM",

		//TeamLeader:     "TL",
		//TeamLeaderAbb:     "TL",
		//TeamLeaderCC:     "TL",

		Manager:    "M",
		ManagerAbb: "M",
		ManagerCC:  "M",

		SeniorManager:    "SM",
		SeniorManagerAbb: "SM",
		SeniorManagerCC:  "SM",

		GeneralManager:    "GM",
		GeneralManagerAbb: "GM",
		GeneralManagerCC:  "GM",

		VicePresident:    "VP",
		VicePresidentAbb: "VP",
		VicePresidentCC:  "VP",
	}
)

var (
	EmployeePositionMap = map[string]string{
		Promoter:    "A",
		PromoterAbb: "A",
		PromoterCC:  "A",

		CSR:   "A",
		CSRCC: "A",

		OutletManager:    "OM",
		OutletManagerAbb: "OM",
		OutletManagerCC:  "OM",

		//TeamLeader:     "TL",
		//TeamLeaderAbb:     "TL",
		//TeamLeaderCC:     "TL",

		Manager:    "M",
		ManagerAbb: "M",
		ManagerCC:  "M",

		SeniorManager:    "SM",
		SeniorManagerAbb: "SM",
		SeniorManagerCC:  "SM",

		GeneralManager:    "GM",
		GeneralManagerAbb: "GM",
		GeneralManagerCC:  "GM",

		VicePresident:    "VP",
		VicePresidentAbb: "VP",
		VicePresidentCC:  "VP",
	}
)

var ProductGroupMap = map[string]struct{}{
	HBBCategory:     struct{}{},
	BaqatiCategory:  struct{}{},
	HayyakCategory:  struct{}{},
	UpgradeCategory: struct{}{},
}

var NotificationTypeMap = map[string]struct{}{
	"NOTIFICATION": struct{}{},
	"OFFER":        struct{}{},
	"TRAINING":     struct{}{},
}

var SupportedFileTypeMap = map[string]struct{}{
	"PDF":   struct{}{},
	"EXCEL": struct{}{},
	"VIDEO": struct{}{},
}
