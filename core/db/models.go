package db

import "sync"

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
}

type Employee struct {
	Id       int
	Name     string
	Position string
	Outlet   string
	Region   string
	Remark   string
}

type SaleRecord struct {
	PreviousProduct string
	CurrentProduct  string
	Quantity        int
	Category        string
}

type SaleTransaction struct {
	EmpID       int
	Position    string
	Outlet      string
	Region      string
	SaleRecords map[string][]*SaleRecord
}

type IndividualIncentive struct {
	Payout      float64
	ExtraPayout float64
	Error       string
}

type TotalIncentive struct {
	HBB     *IndividualIncentive
	Baqati  *IndividualIncentive
	Hayyak  *IndividualIncentive
	Upgrade *IndividualIncentive
}

type ProductGrpTarget struct {
	Baqati  int
	HBB     int
	Hayyak  int
	Upgrade int
}

type ProductPlans struct {
	PlanName    string
	PlanPrice   float64
	SIPRatePlan float64
}

type SaleTarget struct {
	Mutex         sync.RWMutex
	SaleTargetMap map[string]map[string]*ProductGrpTarget
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

var SalesTrgt SaleTarget
var PlanRates PlanRate
var ValOfPoint ValueOfPoint
var MaxPay MaxPayout
