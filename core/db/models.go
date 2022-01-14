package db

import (
	"sync"
	"time"
)

type UserRecord struct {
	EmpId                        int
	Created                      time.Time
	ReportingManager             int
	Isquarterly                  bool
	Month                        int
	Quarter                      int
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
}

type Transaction struct {
	Reqion          string
	Outlet          string
	Id              int
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
	Id        int
	Name      string
	Position  string
	Outlet    string
	Region    string
	MobileNos string
	Remark    string
}

type SaleRecord struct {
	PreviousProduct string
	CurrentProduct  string
	Quantity        int
	Category        string
	SaleDate        time.Time
}

type SaleTransaction struct {
	EmpID       int
	Position    string
	Outlet      string
	Region      string
	SaleRecords map[string][]*SaleRecord
}

type IndividualIncentive struct {
	Payout      float64 `json:"payout"`
	ExtraPayout float64 `json:"extra_payout"`
	Error       string  `json:"error"`
}

type TotalIncentive struct {
	HBB      *IndividualIncentive
	Baqati   *IndividualIncentive
	Hayyak   *IndividualIncentive
	Upgrade  *IndividualIncentive
	TotalPay float64 `json:"total_pay"`
}

type UserAllInfo struct {
	TotalIncentive               *TotalIncentive
	ChartDataResponse            *ChartDataResponse
	HBBSaleCountPerPlanArray     map[string]int
	HayyakSaleCountPerPlanArray  map[string]int
	BaqatiSaleCountPerPlanArray  map[string]int
	UpgradeSaleCountPerPlanArray map[string]int
}

type TargetAchievementData struct {
	SalesCount         int     `json:"sales_count"`
	Target             float64 `json:"target"`
	AchievedPercentage float64 `json:"achieved_percentage"`
	Error              string  `json:"error"`
}

type ChartDataResponse struct {
	HBB     *TargetAchievementData `json:"HBB"`
	Baqati  *TargetAchievementData `json:"Baqati"`
	Hayyak  *TargetAchievementData `json:"Hayyak"`
	Upgrade *TargetAchievementData `json:"Upgrade"`
}

type SaleCountPerPlan struct {
	CurrentPlan string `json:"current_plan"`
	//PreviousPlan string `json:"previous_plan"`
	Quantity int `json:"qty"`
}

type UserHierarchy struct {
	EmployeeID     int
	OutletManager  int
	Manager        int
	SeniorManager  int
	GeneralManager int
	VicePresident  int
}

type ProductGrpTarget struct {
	Baqati  int
	HBB     int
	Hayyak  int
	Upgrade int
}

type ProductWiseSalesCountPerOutlet struct {
	OutletName string
	Region     string
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
	BaqatiSales             int
	BaqatiTarget            float64
	BaqatiAchievedPercent   float64
	HBBSales                int
	HBBTarget               float64
	HBBAchievedPercent      float64
	HayyakSales             int
	HayyakTarget            float64
	HayyakAchievedPercent   float64
	UpgradeSales            int
	UpgradeTarget           float64
	UpgradeAchievedPercent  float64
	TotalSales              int
	TotalTarget             float64
	TotalAchievedPercentage float64
}

type ProductWiseSalesTargetPerOutlet struct {
	OutletName                string
	Region                    string
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
	SaleTargetMap map[string]map[string][12]*ProductGrpTarget
}

type PlanRate struct {
	Mutex       sync.RWMutex
	PlanRateMap map[int]*ProductPlans
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
	VopMap map[string]map[int]*VOP
}

type MaxPayout struct {
	Mutex     sync.RWMutex
	Payoutmap map[string]float64
}

type OtpModel struct {
	StaffID       int
	StaffName     string
	StaffPosition string
	OTP           int
	Retry         int
	CreatedTime   time.Time
	ExpiryTime    time.Time
}

type ManagerHierarchy struct {
	ManagerID             int
	UnderManagerEmployees []int
}

type UserLead struct {
	EmpID    int
	ReportTo int
}

type ReportingStaffIDAndName struct {
	ID     int
	Name   string
	Outlet string
	Region string
}

var SalesTrgt SaleTarget
var PlanRates PlanRate
var ValOfPoint ValueOfPoint
var MaxPay MaxPayout
