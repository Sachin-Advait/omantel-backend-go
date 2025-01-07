package models

import "time"

type UserReportDB struct {
	UserId                       string
	UserName                     string
	Position                     string
	ManagerId                    string
	ManagerPosition              string
	Region                       string
	Outlet                       string
	Latitude                     float64
	Longitude                    float64
	HayyakSalesCount             int
	HayyakTargetCount            float64
	HayyakPayout                 float64
	HayyakExtraPayout            float64
	HayyakAchievedPercentage     float64
	HBBSalesCount                int
	HBBTargetCount               float64
	HBBPayout                    float64
	HBBExtraPayout               float64
	HBBAchievedPercentage        float64
	BaqatiSalesCount             int
	BaqatiTargetCount            float64
	BaqatiPayout                 float64
	BaqatiExtraPayout            float64
	BaqatiAchievedPercentage     float64
	UpgradeSalesCount            int
	UpgradeTargetCount           float64
	UpgradePayout                float64
	UpgradeExtraPayout           float64
	UpradeAchievedPercentage     float64
	VOCPayout                    float64
	VOCExtraPayout               float64
	VOCAchievedPercentage        float64
	TotalPayout                  float64
	HBBSaleCountPerPlanArray     string
	HayyakSaleCountPerPlanArray  string
	BaqatiSaleCountPerPlanArray  string
	UpgradeSaleCountPerPlanArray string
	CreatedAt                    time.Time
	ModifiedAt                   time.Time
}

type OutletReportDB struct {
	Outlet             string
	Region             string
	Latitude           float64
	Longitude          float64
	HBBSalesCount      int
	HBBTargetCount     int
	HayyakSalesCount   int
	HayyakTargetCount  int
	BaqatiSalesCount   int
	BaqatiTargetCount  int
	UpgradeSalesCount  int
	UpgradeTargetCount int
	CreatedAt          time.Time
	ModifiedAt         time.Time
}

type RegionReportDB struct {
	Region             string
	HBBSalesCount      int
	HBBTargetCount     int
	HayyakSalesCount   int
	HayyakTargetCount  int
	BaqatiSalesCount   int
	BaqatiTargetCount  int
	UpgradeSalesCount  int
	UpgradeTargetCount int
	CreatedAt          time.Time
	ModifiedAt         time.Time
}

type UserSalesCountPerPlanDB struct {
	Product      string `json:"product"`
	BusinessUnit struct {
		Consumer   *SaleCalculation `json:"consumer"`
		Corporate  *SaleCalculation `json:"corporate"`
		Government *SaleCalculation `json:"government"`
	} `json:"businessUnit"`
	SubCategory struct {
		CBU map[string]interface{} `json:"CBU"`
		EBU map[string]interface{} `json:"EBU"`
	} `json:"subCategory"`
	Plans struct {
		CBU map[string]int `json:"CBU"`
		EBU map[string]int `json:"EBU"`
	} `json:"plans"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type HbbUserSalesCountPerPlanDB struct {
	Product      string `json:"product"`
	BusinessUnit struct {
		Consumer   *SaleCalculation `json:"consumer"`
		Corporate  *SaleCalculation `json:"corporate"`
		Government *SaleCalculation `json:"government"`
	} `json:"businessUnit"`
	Category struct {
		CBUWFBB *SaleCalculation `json:"CBU WFBB"`
		CBUHBB  *SaleCalculation `json:"CBU HBB"`
		EBUWFBB *SaleCalculation `json:"EBU WFBB"`
		EBUHBB  *SaleCalculation `json:"EBU HBB"`
	} `json:"category"`
	SubCategory struct {
		CBUWFBB map[string]interface{} `json:"CBU WFBB"`
		CBUHBB  map[string]interface{} `json:"CBU HBB"`
		EBUWFBB map[string]interface{} `json:"EBU WFBB"`
		EBUHBB  map[string]interface{} `json:"EBU HBB"`
	} `json:"subCategory"`
	ServiceCategory struct {
		CBUWFBB map[string]interface{} `json:"CBU WFBB"`
		CBUHBB  map[string]interface{} `json:"CBU HBB"`
		EBUWFBB map[string]interface{} `json:"EBU WFBB"`
		EBUHBB  map[string]interface{} `json:"EBU HBB"`
	} `json:"serviceCategory"`
	Plans struct {
		CBUWFBB map[string]int `json:"CBU WFBB"`
		CBUHBB  map[string]int `json:"CBU HBB"`
		EBUWFBB map[string]int `json:"EBU WFBB"`
		EBUHBB  map[string]int `json:"EBU HBB"`
	} `json:"plans"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type SaleCalculation struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percent"`
}

type UserLoginDetailsDB struct {
	UserId    string
	Month     int
	Quarter   int
	Year      int
	LoginTime time.Time
}

type UserBasicDetailsDB struct {
	UserId          string
	Name            string
	Position        string
	MobileNumber    string
	Region          string
	Outlet          string
	ManagerId       string
	ManagerPosition string
	IsAdmin         bool
	IsActive        bool
}

type NotificationDB struct {
	Id         string
	Category   string
	Subject    string
	Body       string
	FileType   string
	FileLink   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  time.Time
	IsUpdate   bool
}

type RedisConfigStruct struct {
	Address                     string
	Password                    string
	Db                          int
	PoolSize                    int
	MinIdleConn                 int
	MaxConnAgeInHour            time.Duration
	PoolTimeoutInMilliSecond    time.Duration
	ReadTimeoutInMilliSecond    time.Duration
	WriteTimeoutInMilliSecond   time.Duration
	ContextTimeoutInMilliSecond time.Duration
	Mode                        string
}
