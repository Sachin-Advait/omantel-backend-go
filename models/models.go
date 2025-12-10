package models

import (
	"sync"
	"time"
)

type GenerateOTP struct {
	ClientId string `json:"client_id"`
}

type VerifyOTP struct {
	IsAdmin  bool   `json:"is_admin"`
	ClientId string `json:"client_id"`
	Otp      int    `json:"otp"`
	Token    string `json:"token"`
}

type UserMapRequest struct {
	OutletName string `json:"outlet"`
	RegionName string `json:"region"`
	Quarterly  string `json:"quarterly"`
	Monthly    string `json:"monthly"`
}

type UserMapResponse struct {
	ClientId   string `json:"client_id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	RegionName string `json:"region"`
	OutletName string `json:"outlet"`
	IsAdmin    bool   `json:"isAdmin"`
}

type UserIncentiveRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Channel   string `json:"channel"`
}

type UserSalesRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Product   string `json:"product"`
	Channel   string `json:"channel"`
}

type MNPSalesRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Channel   string `json:"channel"`
}

type UserReportRequest struct {
	ClientId string `json:"client_id"`
	Month    int    `json:"month"`
	Quarter  int    `json:"quarter"`
	Channel  string `json:"channel"`
}

type VerifyOTPResponse struct {
	Error         string    `json:"error,omitempty"`
	IsAdmin       bool      `json:"is_admin"`
	StaffName     string    `json:"staff_name"`
	StaffPosition string    `json:"staff_position"`
	StaffRole     string    `json:"staff_role"`
	StaffOutlet   string    `json:"staff_outlet"`
	StaffRegion   string    `json:"staff_region"`
	Authenticated bool      `json:"authenticated"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

type UserRecord struct {
	EmpId                        string
	ReportingManager             string
	Isquarterly                  bool
	Month                        int
	Quarter                      int
	Year                         int
	Hayyak_sales_count           int
	Hayyak_target                float64
	Hayyak_payout                float64
	Hayyak_extra_payout          float64
	Hayyak_achieved_percentage   float64
	Hbb_sales_count              int
	Hbb_target                   float64
	Hbb_payout                   float64
	Hbb_extra_payout             float64
	Hbb_achieved_percentage      float64
	Baqati_sales_count           int
	Baqati_target                float64
	Baqati_payout                float64
	Baqati_extra_payout          float64
	Baqati_achieved_percentage   float64
	Upgrade_sales_count          int
	Upgrade_target               float64
	Upgrade_payout               float64
	Upgrade_extra_payout         float64
	Uprade_achieved_percentage   float64
	Total_payout                 float64
	HBBSaleCountPerPlanArray     string
	HayyakSaleCountPerPlanArray  string
	BaqatiSaleCountPerPlanArray  string
	UpgradeSaleCountPerPlanArray string
	CreatedAt                    time.Time
}

type Transaction struct {
	Reqion          string
	Outlet          string
	Id              string
	CategoryType    string
	PreviousProduct string
	CurrentProduct  string
	Quantity        int
	Category        string
	ProductGrp      string
	SaleDate        time.Time
	Quarter         int
}

type Employee struct {
	Id        string
	Name      string
	Position  string
	Outlet    string
	Region    string
	MobileNos string
	Remark    string
}

type SaleRecord struct {
	TrxId           int
	PreviousProduct string
	CurrentProduct  string
	Quantity        int
	Category        string
	MNPDonor        string
	SaleDate        time.Time
}

type SaleTransaction struct {
	EmpID       string
	Position    string
	Outlet      string
	Region      string
	SaleRecords map[string][]*SaleRecord
}

type SaleCountRecord struct {
	ProductGroup         string
	CurrentPlan          string
	Category             string
	BusinessUnit         string
	SubCategory          string
	LineType             string
	FinalServiceCategory string
	MNPDonor             string
	Quantity             int
}

type IndividualIncentive struct {
	Payout      float64 `json:"payout"`
	ExtraPayout float64 `json:"extra_payout"`
	// Terminated   int     `json:"terminated"`
	// DeductAmount float64 `json:"deduct_Amount"`
	Error string `json:"error,omitempty"`
}

type TotalIncentive struct {
	HBB                     *IndividualIncentive
	Baqati                  *IndividualIncentive
	Hayyak                  *IndividualIncentive
	Upgrade                 *IndividualIncentive
	VOC                     *IndividualIncentive
	TotalPayBeforeDeduction float64 `json:"total_pay_before_deduction"` // after adjustment
	TotalPay                float64 `json:"total_pay"`                  // before adjustment
}

type MNPSalesCount struct {
	Total      int              `json:"total"`
	Baqati     *MNPCategorySale `json:"Baqati"`
	Hayyak     *MNPCategorySale `json:"Hayyak"`
	CreatedAt  time.Time        `json:"createdAt"`
	ModifiedAt time.Time        `json:"modifiedAt"`
}

type MNPCategorySale struct {
	Count      int                     `json:"count"`
	Percentage float64                 `json:"percentage"`
	Error      string                  `json:"error,omitempty"`
	Types      map[string]*EachMNPSale `json:"types"`
}

type EachMNPSale struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type MNPPlansCount struct {
	Total      map[string]map[string]int `json:"total"`
	Baqati     map[string]map[string]int `json:"Baqati"`
	Hayyak     map[string]map[string]int `json:"Hayyak"`
	Error      string                    `json:"error,omitempty"`
	CreatedAt  time.Time                 `json:"createdAt"`
	ModifiedAt time.Time                 `json:"modifiedAt"`
}

type TargetAchievementData struct {
	SalesCount         int     `json:"sales_count"`
	Target             float64 `json:"target"`
	AchievedPercentage float64 `json:"achieved_percentage"`
	Error              string  `json:"error,omitempty"`
}

type ChartDataResponse struct {
	HBB     *TargetAchievementData `json:"HBB"`
	Baqati  *TargetAchievementData `json:"Baqati"`
	Hayyak  *TargetAchievementData `json:"Hayyak"`
	Upgrade *TargetAchievementData `json:"Upgrade"`
}

type UserSaleCount struct {
	Product []struct {
		BusinessUnit []struct {
			SubCategory []struct {
				PlanName []struct {
					Quantity int `json:"qty"`
				} `json:"current_plan"`
			} `json:"sub_category"`
		} `json:"business_unit"`
	} `json:"product"`
}

type SaleCountPerPlan struct {
	CurrentPlan string `json:"current_plan"`
	//PreviousPlan string `json:"previous_plan"`
	Quantity int `json:"qty"`
}

type UserHierarchy struct {
	EmployeeID     string
	OutletManager  string
	Manager        string
	SeniorManager  string
	GeneralManager string
	VicePresident  string
}

type ProductGrpTarget struct {
	Baqati  int
	HBB     int
	Hayyak  int
	Upgrade int
}

type ProductWiseSalesCountPerOutlet struct {
	OutletName string `json:"OutletName,omitempty"`
	Region     string `json:"Region,omitempty"`
	Latitude   float64
	Longitude  float64
	Baqati     int
	HBB        int
	Hayyak     int
	Upgrade    int
	Total      int
}

type ProductWiseSalesCountPerStaff struct {
	Name    string
	Baqati  int
	HBB     int
	Hayyak  int
	Upgrade int
	Total   int
}

type ProductWiseSalesTargetPerStaff struct {
	Name                    string
	HBBSales                int
	HBBTarget               float64
	HBBAchievedPercent      float64
	HayyakSales             int
	HayyakTarget            float64
	HayyakAchievedPercent   float64
	BaqatiSales             int
	BaqatiTarget            float64
	BaqatiAchievedPercent   float64
	UpgradeSales            int
	UpgradeTarget           float64
	UpgradeAchievedPercent  float64
	TotalSales              int
	TotalTarget             float64
	TotalAchievedPercentage float64
}

type ProductWiseSalesTargetPerOutlet struct {
	OutletName                string `json:"OutletName,omitempty"`
	Region                    string `json:"Region,omitempty"`
	BaqatiSales               int
	BaqatiTarget              int
	BaqatiAchievedPercentage  float64
	HBBSales                  int
	HBBTarget                 int
	HBBAchievedPercentage     float64
	HayyakSales               int
	HayyakTarget              int
	HayyakAchievedPercentage  float64
	UpgradeSales              int
	UpgradeTarget             int
	UpgradeAchievedPercentage float64
	TotalSales                int
	TotalTarget               int
	TotalAchievedPercentage   float64
}

type ProductPlans struct {
	PlanName    string
	PlanPrice   float64
	SIPRatePlan float64
}

type SaleTarget struct {
	Mutex         sync.RWMutex
	SaleTargetMap map[int]map[string]map[string][12]*ProductGrpTarget
}

type PlanRate struct {
	Mutex       sync.RWMutex
	PlanRateMap map[int]map[int]map[string]*ProductPlans
}

type Calculator struct {
	PrevPlanID int
	Qty        int
}

type VOP struct {
	Point    float64
	ExtraPay float64
}

type ValueOfPoint struct {
	Mutex  sync.RWMutex
	VopMap map[int]map[int]map[string]map[int]*VOP
}

type MaxPayout struct {
	Mutex     sync.RWMutex
	Payoutmap map[int]map[int]map[string]map[string]float64
}

type MaxExtraPayout struct {
	Mutex          sync.RWMutex
	ExtraPayoutmap map[int]map[int]map[string]map[string]float64
}

type ThresholdCheckout struct {
	Mutex sync.RWMutex
	Limit map[int]map[int]map[string]map[string]map[int]float64
}

type VOCDetails struct {
	Mutex    sync.RWMutex
	Criteria map[int]map[string]float64
	MaxPay   map[int]map[string]float64
}

type OtpModel struct {
	StaffID              string
	StaffName            string
	StaffPosition        string
	OTP                  int
	IsActiveUser         bool
	OtpLastUpdateTime    time.Time
	OTPResendCount       int
	OTPVerifyFailedCount int
	OTPUsed              bool
	OtpCreateTime        time.Time
	LastLoginTime        time.Time
	CreatedTime          time.Time
	UpdatedTime          time.Time
	UpdateReason         string
}

type ManagerHierarchy struct {
	ManagerID             string
	UnderManagerEmployees []string
}

type UserLead struct {
	EmpID    string
	ReportTo string
}

type ReportingStaffIDAndName struct {
	ID       string
	Position string
	Name     string
	Outlet   string
	Region   string
}

// CaptchRequest is struct for requests of /verifyCaptcha API
type CaptchRequest struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserLoginCountRequest struct {
	Region string    `json:"region"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type UserLoginCountResponse struct {
	TotalUsers         int     `json:"totalUser"`
	LoginCount         int     `json:"loginCount"`
	NonLoginCount      int     `json:"nonLoginCount"`
	LoginPercentage    float64 `json:"loginPercentage"`
	NonLoginPercentage float64 `json:"nonLoginPercentage"`
}

type UserLoginDetailsRequest struct {
	Page   int       `json:"page"`
	Size   int       `json:"size"`
	Region string    `json:"region"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Login  bool      `json:"login"`
}

type UserLoginDetailsResponse struct {
	NextPage int           `json:"nextPage"`
	Size     int           `json:"size"`
	Users    []*LoginUsers `json:"users"`
}

type LoginUsers struct {
	StaffId            string    `json:"userId"`
	StaffName          string    `json:"staff_name"`
	StaffMobileNo      string    `json:"staff_mobile_no"`
	StaffPosition      string    `json:"staff_position"`
	StaffOutlet        string    `json:"staff_outlet"`
	StaffRegion        string    `json:"staff_region"`
	SuperVisorId       string    `json:"supervisor_id"`
	SuperVisorName     string    `json:"supervisor_name"`
	SuperVisorPosition string    `json:"supervisor_position"`
	LoginTime          time.Time `json:"loginTime"`
}

var SalesTrgt SaleTarget
var PlanRates PlanRate
var ValOfPoint ValueOfPoint
var MaxPay MaxPayout
var MaxExtraPay MaxExtraPayout
var ThresholdCheck ThresholdCheckout
var VOCDetail VOCDetails
