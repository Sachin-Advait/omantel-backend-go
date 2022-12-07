package models

import (
	"sync"
	"time"
)

type SingleUserIncentiveRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
}

type SalesCountRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Product   string `json:"product"`
}

type GenerateOTP struct {
	ClientId string `json:"client_id"`
}

type VerifyOTP struct {
	ClientId string `json:"client_id"`
	Otp      int    `json:"otp"`
}

type VerifyOTPResponse struct {
	Error         string    `json:"error"`
	StaffName     string    `json:"staff_name"`
	StaffPosition string    `json:"staff_position"`
	Authenticated bool      `json:"authenticated"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

type UserReportRequest struct {
	ClientId string `json:"client_id"`
	Month    int    `json:"month"`
	Quarter  int    `json:"quarter"`
}

type UserReporting struct {
	ClientId string `json:"client_id"`
}

type NewRegionRequest struct {
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
	IsUpdate   bool   `json:"isUpdate"`
}

type NewRegionFailedResponse struct {
	RegionCode string `json:"region_code"`
	Message    string `json:"message"`
}

type DeleteRegionRequest struct {
	RegionCodes []string `json:"region_codes"`
}

type GetRegionResponse struct {
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
}

type NewOutletRequest struct {
	OutletCode string  `json:"outlet_code"`
	OutletName string  `json:"outlet_name"`
	RegionCode string  `json:"region_code"`
	Lat        float64 `json:"lat"`
	Lang       float64 `json:"lang"`
	IsUpdate   bool    `json:"isUpdate"`
}

type NewOutletFailedResponse struct {
	OutletCode string `json:"outlet_code"`
	Message    string `json:"message"`
}

type DeleteOutletRequest struct {
	OutletCodes []string `json:"outlet_codes"`
}

type GetOutletResponse struct {
	OutletCode string  `json:"outlet_code"`
	OutletName string  `json:"outlet_name"`
	RegionCode string  `json:"region_code"`
	RegionName string  `json:"region_name"`
	Lat        float64 `json:"lat"`
	Lang       float64 `json:"lang"`
}

type NewProductRequest struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Weightage   float64 `json:"weightage"`
	IsUpdate    bool    `json:"isUpdate"`
}

type NewProductFailedResponse struct {
	ProductCode string `json:"product_code"`
	Message     string `json:"message"`
}

type DeleteProductRequest struct {
	ProductCodes []string `json:"product_codes"`
}

type GetProductResponse struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Weightage   float64 `json:"weightage"`
}

type NewPlanRequest struct {
	PlanCode    string  `json:"plan_code"`
	PlanName    string  `json:"plan_name"`
	ProductCode string  `json:"product_code"`
	PlanPrice   float64 `json:"plan_price"`
	SIPPlanRate float64 `json:"sip_plan_rate"`
	IsUpdate    bool    `json:"isUpdate"`
}

type NewPlanFailedResponse struct {
	PlanCode string `json:"plan_code"`
	Message  string `json:"message"`
}

type DeletePlanRequest struct {
	PlanCodes []string `json:"plan_codes"`
}

type GetPlanResponse struct {
	PlanCode    string  `json:"plan_code"`
	PlanName    string  `json:"plan_name"`
	ProductCode string  `json:"product_code"`
	PlanPrice   float64 `json:"plan_price"`
	SIPPlanRate float64 `json:"sip_plan_rate"`
}

type NewTargetRequest struct {
	TargetCode string  `json:"target_code"`
	RegionCode string  `json:"region_code"`
	OutletCode string  `json:"outlet_code"`
	Year       int     `json:"year"`
	Month      int     `json:"month"`
	Baqati     float64 `json:"baqati"`
	Hayyak     float64 `json:"hayyak"`
	HBB        float64 `json:"hbb"`
	Upgrade    float64 `json:"upgrade"`
	IsUpdate   bool    `json:"isUpdate"`
}

type NewTargetFailedResponse struct {
	TargetCode string `json:"target_code"`
	Message    string `json:"message"`
}

type DeleteTargetRequest struct {
	TargetCodes []string `json:"target_codes"`
}

type GetTargetRequest struct {
	TargetCode []string `json:"target_codes"`
	OutletCode []string `json:"outlet_codes"`
	Month      int      `json:"month"`
	Year       int      `json:"year"`
}

type GetTargetResponse struct {
	TargetCode string  `json:"target_code"`
	RegionCode string  `json:"region_code"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Baqati     float64 `json:"baqati"`
	Hayyak     float64 `json:"hayyak"`
	HBB        float64 `json:"hbb"`
	Upgrade    float64 `json:"upgrade"`
}

/*

 */

type DbClientConfig struct {
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
	CreateAt                     time.Time
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
	PreviousProduct string
	CurrentProduct  string
	Quantity        int
	Category        string
	SaleDate        time.Time
}

type SaleTransaction struct {
	EmpID       string
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
	SaleTargetMap map[string]map[string][12]*ProductGrpTarget
}

type PlanRate struct {
	Mutex       sync.RWMutex
	PlanRateMap map[string]*ProductPlans
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

var SalesTrgt SaleTarget
var PlanRates PlanRate
var ValOfPoint ValueOfPoint
var MaxPay MaxPayout

var (
	ProductTypeMap = map[string]string{
		"EBU Fixed":    "HBB",
		"CBU Fixed":    "HBB",
		"EBU Prepaid":  "Hayyak",
		"CBU Prepaid":  "Hayyak",
		"EBU Postpaid": "Baqati",
		"CBU Postpaid": "Baqati",
	}
)

var (
	ShortCode          = "OMSIP"
	RequestorChannelId = "OMSIP"
)

var SMSRequestTemplate = []byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
<soapenv:Header/>
<soapenv:Body>
   <sen:SendNotificationRq>
	  <sen:RequestHeader>
		 <FunctionId>{{.FunctionId}}</FunctionId>
		 <TransactionId>{{.TransactionId}}</TransactionId>
		 <RequestId>{{.RequestId}}</RequestId>
		 <RequestorChannelId>{{.RequestorChannelId}}</RequestorChannelId>
		 <RequestorUserId>{{.RequestorUserId}}</RequestorUserId>
		 <Timestamp>{{.Timestamp}}</Timestamp>
	  </sen:RequestHeader>
	  <sen:RequestBody>
		 <sen:Notification>
			<NotificationType>SMS</NotificationType>
			<NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
			<NotificationLanguage>1</NotificationLanguage>
			<ShortCode>{{.ShortCode}}</ShortCode>
			<IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
			<NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
			<Body>{{.Body}}</Body>
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

// var SMSResponseTemplate = []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
